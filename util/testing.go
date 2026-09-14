package util

import (
	"log"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

// global test helpers
const (
	VaultAddress = "http://127.0.0.1:8200"
	Container    = "my_bucket"
	Prefix       = "prefix"
	AppRole      = "myAppRole"
	tokenFile    = "/tmp/vault-test-root-token"
)

var (
	VaultToken          = rootToken()
	VaultClient         = basicVaultClient()
	RoleID, SecretID, _ = approleAttrs()
)

// helper for retrieving root token from bootstrap
func rootToken() string {
	// retrieve root token
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		// return unauthenticated client as next best option
		return ""
	}
	return strings.TrimSpace(string(data))
}

// helper for basic vault client
func basicVaultClient() *vault.Client {
	// initialize config and client
	vaultConfig := &vault.Config{Address: VaultAddress}
	vaultConfig.ConfigureTLS(&vault.TLSConfig{Insecure: true})
	client, _ := vault.NewClient(vaultConfig)
	client.SetToken(VaultToken)

	return client
}

// helper for approle auth
func approleAttrs() (string, string, error) {
	// retrieve role id and secret id for testing approle auth in "push" mode
	roleID, err := VaultClient.Logical().Read("auth/approle/role/" + AppRole + "/role-id")
	if err != nil {
		log.Print("failed to retrieve role ID for approle auth")
		return "", "", err
	}
	secretID, err := VaultClient.Logical().Write("auth/approle/role/"+AppRole+"/secret-id", nil)
	if err != nil {
		log.Print("failed to retrieve secret ID for approle auth")
		return "", "", err
	}

	return roleID.Data["role_id"].(string), secretID.Data["secret_id"].(string), nil
}
