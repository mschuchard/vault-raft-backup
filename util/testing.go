package util

import (
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

// global test helpers
const (
	VaultAddress = "http://127.0.0.1:8200"
	Container    = "my_bucket"
	Prefix       = "prefix"
	tokenFile    = "/tmp/vault-test-root-token"
)

var (
	VaultToken  = rootToken()
	VaultClient = basicVaultClient()
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
