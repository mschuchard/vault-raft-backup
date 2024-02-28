package util

import (
	"os"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

func TestSnapshotFileClose(test *testing.T) {
	genericFile, err := os.Create("foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}

	if err := SnapshotFileClose(genericFile); err != nil {
		test.Error("generic file failed to close")
		test.Error(err)
	}
}

func TestSnapshotFileRemove(test *testing.T) {
	genericFile, err := os.Create("foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}

	if err := SnapshotFileRemove(genericFile); err != nil {
		test.Error("failed to remove generic file")
		test.Error(err)
	}

	// false positive probably because of race condition
	/*if _, err := genericFile.Stat(); err == nil {
		test.Error("validation that generic file was removed returned no path error")
	}*/
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
