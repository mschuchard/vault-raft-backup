package vault

import (
	"context"
	"errors"
	"log"
	"net/url"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"

	"github.com/mschuchard/vault-raft-backup/util"
)

// authentication engine with pseudo-enum
type authEngine string

const (
	awsIam     authEngine = "aws"
	vaultToken authEngine = "token"
)

// configured vault client validated constructor
func NewVaultClient(backupVaultConfig *util.VaultConfig) (*vault.Client, error) {
	// vault address default
	address := backupVaultConfig.Address
	if len(address) == 0 {
		address = "http://127.0.0.1:8200"
	} else {
		// vault address validation
		if url, err := url.ParseRequestURI(address); err != nil || len(url.Scheme) == 0 || len(url.Host) == 0 {
			log.Printf("%s is not a valid Vault server address", address)

			// assign err if it is nil
			if err == nil {
				err = errors.New("invalid Vault server address")
			}

			return nil, err
		}
	}

	// validate insecure
	insecure := backupVaultConfig.Insecure
	if !insecure && address[0:5] == "http:" {
		log.Print("insecure input parameter was omitted or specified as false, and address protocol is http")
		log.Print("insecure will be reset to value of true")
		insecure = true
	}

	// initialize locals
	engine := authEngine(backupVaultConfig.Engine)
	token := backupVaultConfig.Token
	awsMountPath := backupVaultConfig.AWSMountPath
	awsRole := backupVaultConfig.AWSRole

	// determine vault auth engine if unspecified
	if len(engine) == 0 {
		log.Print("authentication engine for Vault not specified; using logic from other parameters to assist with determination")

		// validate inputs specified for only one engine
		if len(token) > 0 && (len(awsMountPath) > 0 || len(awsRole) > 0) {
			log.Print("token and AWS mount path or AWS role were simultaneously specified; these are mutually exclusive options")
			log.Print("intended authentication engine could not be determined from other parameters")
			return nil, errors.New("unable to deduce authentication engine")
		}
		if len(token) == 0 {
			log.Print("AWS IAM authentication will be utilized with the Vault client")
			engine = awsIam
		} else {
			log.Print("token authentication will be utilized with the Vault client")
			engine = vaultToken
		}
	} else if engine != awsIam && engine != vaultToken { // validate engine if unspecified
		log.Printf("%v was input as an authentication engine, but only token and aws are supported", engine)
		return nil, errors.New("invalid Vault authentication engine")
	}

	// validate vault token
	if engine == vaultToken && len(token) != 28 {
		log.Print("the specified Vault Token is invalid")
		return nil, errors.New("invalid vault token")
	}

	// default aws mount path and role
	if engine == awsIam {
		if len(awsMountPath) == 0 {
			log.Print("using default AWS authentication mount path at 'aws'")
			awsMountPath = "aws"
		}
		if len(awsRole) == 0 {
			log.Print("using Vault role in utilized AWS authentication engine with the same name as the current utilized AWS IAM Role")
		}
	}

	// initialize vault api config
	vaultConfig := &vault.Config{Address: address}
	if err := vaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: insecure}); err != nil {
		log.Print("Vault TLS configuration failed to initialize")
		return nil, err
	}

	// initialize vault client
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Print("Vault client failed to initialize")
		return nil, err
	}

	// verify vault is unsealed
	sealStatus, err := client.Sys().SealStatus()
	if err != nil {
		log.Print("unable to verify that the Vault cluster is unsealed")
		return nil, err
	}
	if sealStatus.Sealed {
		log.Print("the Vault server cluster is sealed and no operations can be executed")
		return nil, errors.New("vault sealed")
	}

	// determine authentication method
	switch engine {
	case vaultToken:
		client.SetToken(token)
	case awsIam:
		// determine iam role login option
		var loginOption auth.LoginOption

		if len(awsRole) > 0 {
			// use explicitly specified iam role
			loginOption = auth.WithRole(awsRole)
		} else {
			// use default iam role
			loginOption = auth.WithIAMAuth()
		}

		// authenticate with aws iam
		awsAuth, err := auth.NewAWSAuth(loginOption)
		if err != nil {
			return nil, errors.New("unable to initialize AWS IAM authentication")
		}

		// utilize aws authentication with vault client
		authInfo, err := client.Auth().Login(context.Background(), awsAuth)
		if err != nil {
			log.Print("unable to authenticate to Vault via AWS IAM auth method")
			return nil, err
		}
		if authInfo == nil {
			return nil, errors.New("no auth info was returned after login")
		}
	}

	// return authenticated vault client
	return client, nil
}
