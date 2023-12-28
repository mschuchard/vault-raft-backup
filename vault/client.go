package vault

import (
	"context"
	"errors"
	"fmt"
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
	// initialize config
	vaultConfig := &vault.Config{Address: config.address}
	err := vaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: config.insecure})
	if err != nil {
		fmt.Println("Vault TLS configuration failed to initialize")
		fmt.Println(err)
		return nil, err
	}

	// initialize client
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		fmt.Println("Vault client failed to initialize")
		fmt.Println(err)
		return nil, err
	}

	// determine authentication method
	if token == "aws-iam" {
		// authenticate with aws iam
		awsAuth, err := auth.NewAWSAuth(auth.WithIAMAuth())
		if err != nil {
			return nil, errors.New("Unable to initialize AWS IAM authentication")
		}

		authInfo, err := client.Auth().Login(context.TODO(), awsAuth)
		if err != nil {
			return nil, errors.New("Unable to login to AWS IAM auth method")
		}
		if authInfo == nil {
			return nil, errors.New("No auth info was returned after login")
		}
	} else {
		// authenticate with token
		if len(token) != 26 {
			return nil, errors.New("The Vault token is invalid")
		}
		client.SetToken(token)
	}

	// return vault client interface
	return client, nil
}
