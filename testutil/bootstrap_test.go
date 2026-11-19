package testutil

import (
	"os"
	"testing"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// bootstrap vault server for testing
func TestBootstrap(test *testing.T) {
	// check if we should skip bootstrap
	health, err := VaultClient.Sys().Health()
	if err == nil && health.Initialized && !health.Sealed {
		test.Skip("Vault server already initialized; skipping bootstrap")
	}

	// initialize single key unseal
	initResponse, err := VaultClient.Sys().Init(&vault.InitRequest{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		test.Fatalf("vault initialization failed: %s", err)
	}

	// unseal with the single key
	if sealResponse, err := VaultClient.Sys().Unseal(initResponse.Keys[0]); err != nil || !sealResponse.Initialized || sealResponse.Sealed {
		test.Fatalf("vault unseal failed: %s", err)
	}

	// persist root token for subsequent tests
	if err := os.WriteFile(tokenFile, []byte(initResponse.RootToken), 0o600); err != nil {
		test.Fatalf("failed to write root token to %s: %s", tokenFile, err)
	}

	// authenticate the client with root token for further vault configuration
	VaultClient.SetToken(initResponse.RootToken)

	// wait for raft leader election before configuring vault
	for range 10 {
		// cluster healthy and leader elected?
		if health, err := VaultClient.Sys().Health(); err == nil && !health.Standby {
			break
		}
		// otherwise wait and try again
		time.Sleep(1 * time.Second)
	}

	// enable auth: approle, aws, kubernetes (token enabled by default with dev server)
	if err := VaultClient.Sys().EnableAuthWithOptions("approle", &vault.EnableAuthOptions{Type: "approle"}); err != nil {
		test.Fatalf("failed to enable approle auth: %s", err)
	}
	if _, err := VaultClient.Logical().Write("auth/approle/role/"+appRole, map[string]any{
		"token_policies": "default",
		"token_ttl":      "1h",
		"token_max_ttl":  "4h",
	}); err != nil {
		test.Fatalf("failed to configure approle auth: %s", err)
	}
	if err := VaultClient.Sys().EnableAuthWithOptions("aws", &vault.EnableAuthOptions{Type: "aws"}); err != nil {
		test.Fatalf("failed to enable aws auth: %s", err)
	}
	if err := VaultClient.Sys().EnableAuthWithOptions("kubernetes", &vault.EnableAuthOptions{Type: "kubernetes"}); err != nil {
		test.Fatalf("failed to enable kubernetes auth: %s", err)
	}
}
