package vault

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
)

// authentication engine with pseudo-enum
type authEngine string

const (
	awsIam     authEngine = "aws"
	vaultToken authEngine = "token"
)

// VaultConfig defines vault api client interface
type VaultConfig struct {
	address      string
	insecure     bool
	engine       authEngine
	token        string
	awsMountPath string
	awsRole      string
	snapshotPath string
}

// vault config constructor
func NewVaultConfig() *VaultConfig {
	// vault address default
	address := os.Getenv("VAULT_ADDR")
	if len(address) == 0 {
		address = "http://127.0.0.1:8200"
	} else {
		// vault address validation
		if _, err := url.ParseRequestURI(address); err != nil {
			log.Fatalf("%s is not a valid Vault server address", address)
		}
	}
	// validate insecure
	insecure, err := strconv.ParseBool(os.Getenv("VAULT_SKIP_VERIFY"))
	if err != nil {
		log.Fatal("invalid boolean value for VAULT_SKIP_VERIFY")
	}
	// determine vault auth engine if unspecified
	engine := authEngine(os.Getenv("VAULT_AUTH_ENGINE"))
	token := os.Getenv("VAULT_TOKEN")
	awsMountPath := os.Getenv("VAULT_AWS_MOUNT")

	if len(engine) == 0 {
		log.Print("authentication engine for Vault not specified; using logic from other parameters to assist with determination")

		// validate inputs specified for only one engine
		if len(token) > 0 && len(awsMountPath) > 0 {
			log.Fatal("token and AWS mount path were simultaneously specified; these are mutually exclusive options")
		}
		if len(token) == 0 {
			log.Print("AWS IAM authentication will be utilized with the Vault client")
			engine = awsIam
		} else {
			log.Print("token authentication will be utilized with the Vault client")
			engine = vaultToken
		}
	}
	// validate vault token
	awsRole := os.Getenv("VAULT_AWS_ROLE")

	if engine == vaultToken && len(token) != 28 {
		log.Fatal("the specified Vault Token is invalid")
	} else {
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
	}

	// initialize vault config
	vaultConfig := &VaultConfig{
		address:      address,
		insecure:     insecure,
		engine:       engine,
		token:        token,
		awsMountPath: awsMountPath,
		awsRole:      awsRole,
		snapshotPath: os.Getenv("VAULT_SNAPSHOT_PATH"),
	}

	return vaultConfig
}

// snapshot path reader
func (config *VaultConfig) SnapshotPath() string {
	return config.snapshotPath
}

// vault client configuration
func NewVaultClient(config *VaultConfig) (*vault.Client, error) {
	// initialize vault api config
	vaultConfig := &vault.Config{Address: config.address}
	if err := vaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: config.insecure}); err != nil {
		log.Print("Vault TLS configuration failed to initialize")
		return nil, err
	}

	// initialize vault client
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Print("Vault client failed to initialize")
		return nil, err
	}

	// determine authentication method
	switch config.engine {
	case vaultToken:
		client.SetToken(config.token)
	case awsIam:
		// determine iam role login option
		var loginOption auth.LoginOption

		if len(config.awsRole) > 0 {
			// use explicitly specified iam role
			loginOption = auth.WithRole(config.awsRole)
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
			return nil, errors.New("unable to login to AWS IAM auth method")
		}
		if authInfo == nil {
			return nil, errors.New("no auth info was returned after login")
		}
	}

	// return authenticated vault client
	return client, nil
}
