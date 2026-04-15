package util

import (
	"os"
	"testing"
	"time"

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
	os.Remove("foo")
}

func TestSnapshotFileRemove(test *testing.T) {
	genericFile, err := os.Create("foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}

	if err := SnapshotFileRemove(genericFile); err != nil {
		test.Error("failed to remove generic file")
		test.Error(err)
		os.Remove("foo")
	}

	genericFile.Close()
	if _, err = genericFile.Stat(); err == nil {
		test.Error("validation that generic file was removed returned no path error")
	}

	if err = SnapshotFileRemove(genericFile); err == nil || err.Error() != "snapshot not removed" {
		test.Error("unexpected or no error returned")
		test.Errorf("expected: snapshot not removed, actual: %s", err)
	}
}

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
		// leader elected?
		health, err = VaultClient.Sys().Health()
		// ...then continue
		if err == nil && !health.Standby {
			break
		}
		// otherwise wait and try again
		time.Sleep(1 * time.Second)
	}

	// enable auth: aws
	if err := VaultClient.Sys().EnableAuthWithOptions("aws", &vault.EnableAuthOptions{Type: "aws"}); err != nil {
		test.Fatalf("failed to enable aws auth: %s", err)
	}
}
