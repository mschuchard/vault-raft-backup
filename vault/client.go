package vault

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
)

// VaultConfig is for vault interface
type VaultConfig struct {
	vaultAddr    string
	token        string
	snapshotPath string
	insecure     bool
}

// vault config constructor
func NewVaultConfig() *VaultConfig {
	// initialize vaultConfig
	insecure, err := strconv.ParseBool(os.Getenv("VAULT_SKIP_VERIFY"))
	if err != nil {
		log.Fatalln("Invalid boolean value for VAULT_SKIP_VERIFY")
	}
	vaultConfig := &VaultConfig{
		vaultAddr:    os.Getenv("VAULT_ADDR"),
		token:        os.Getenv("VAULT_TOKEN"),
		snapshotPath: os.Getenv("VAULT_SNAPSHOT_PATH"),
		insecure:     insecure,
	}

	return vaultConfig
}

// snapshot path reader
func (config *VaultConfig) SnapshotPath() string {
	return config.snapshotPath
}

// vault client configuration
func VaultClient(config *VaultConfig) (*vault.Client, error) {
	// initialize config
	vaultConfig := &vault.Config{Address: config.vaultAddr}
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
	if config.token == "aws-iam" {
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
		if len(config.token) != 26 {
			return nil, errors.New("The Vault token is invalid")
		}
		client.SetToken(config.token)
	}

	// return vault client interface
	return client, nil
}
