package util

import (
	"os"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

func TestSnapshotFileClose(test *testing.T) {
	genericFile, err := os.OpenFile("foo", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	SnapshotFileClose(genericFile)
}

// bootstrap vault server for testing
func TestBootstrap(test *testing.T) {
	// check if we should skip bootstrap
	auths, _ := VaultClient.Sys().ListAuth()
	if _, ok := auths["auth/aws/"]; ok {
		test.Skip("Vault server already bootstrapped; skipping")
	}

	// enable auth: aws
	VaultClient.Sys().EnableAuthWithOptions("auth/aws", &vault.EnableAuthOptions{Type: "aws"})
}
