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

// vaultConfig defines vault api client interface
type vaultConfig struct {
	address      string
	insecure     bool
	engine       authEngine
	token        string
	awsMountPath string
	awsRole      string
	snapshotPath string
}

// vault config constructor
func NewVaultConfig() (*vaultConfig, error) {
	// vault address default
	address := os.Getenv("VAULT_ADDR")
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
	insecure, err := strconv.ParseBool(os.Getenv("VAULT_SKIP_VERIFY"))
	if err != nil {
		// no value specified so assign based on address
		if len(os.Getenv("VAULT_SKIP_VERIFY")) == 0 {
			// https --> false
			if address[0:5] == "https" {
				insecure = false
			} else { // http --> true
				insecure = true
			}
		} else { // assigned value could not be converted to boolean
			log.Printf("invalid boolean value '%s' for VAULT_SKIP_VERIFY", os.Getenv("VAULT_SKIP_VERIFY"))
			return nil, errors.New("invalid VAULT_SKIP_VERIFY value")
		}
	}

	// determine vault auth engine if unspecified
	engine := authEngine(os.Getenv("VAULT_AUTH_ENGINE"))
	token := os.Getenv("VAULT_TOKEN")
	awsMountPath := os.Getenv("VAULT_AWS_MOUNT")
	awsRole := os.Getenv("VAULT_AWS_ROLE")

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
	} else { // validate engine if unspecified
		if engine != awsIam && engine != vaultToken {
			log.Printf("%v was input as an authentication engine, but only token and aws are supported", engine)
			return nil, errors.New("invalid Vault authentication engine")
		}
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

	// provide snapshot path default if unspecified
	snapshotPath := os.Getenv("VAULT_SNAPSHOT_PATH")
	if len(snapshotPath) == 0 {
		// assign default path in tmp-dir
		snapshotPath = os.TempDir() + "/vault.bak"
		log.Printf("vault raft snapshot path defaulting to '%s'", snapshotPath)
	}

	// return initialized vault config
	return &vaultConfig{
		address:      address,
		insecure:     insecure,
		engine:       engine,
		token:        token,
		awsMountPath: awsMountPath,
		awsRole:      awsRole,
		snapshotPath: snapshotPath,
	}, nil
}

// snapshot path reader
func (config *vaultConfig) SnapshotPath() string {
	return config.snapshotPath
}

// vault client configuration
func NewVaultClient(config *vaultConfig) (*vault.Client, error) {
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
