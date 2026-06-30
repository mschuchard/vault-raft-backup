package vault

import (
	"log"
	"os"
	"testing"
	"time"

	vault "github.com/hashicorp/vault/api"
	"github.com/mschuchard/vault-raft-backup/util"
)

func TestVaultRaftSnapshotCreate(test *testing.T) {
	if err := VaultRaftSnapshotCreate(util.VaultClient, "vault.bak"); err != nil {
		test.Error("vault raft snapshot creation failed")
		test.Error(err)
	}
	if _, err := os.Stat("./vault.bak"); err != nil {
		test.Error("vault raft snapshot file was not actually created")
		test.Error(err)
	}

	os.Remove("./vault.bak")

	if err := VaultRaftSnapshotCreate(util.VaultClient, "/foo/vault.bak"); err == nil || err.Error() != "open /foo/vault.bak: no such file or directory" {
		test.Errorf("expected error (contains): open /foo/vault.bak: no such file or directory, actual: %s", err)
	}
}

func TestVaultRaftSnapshotRestore(test *testing.T) {
	if err := VaultRaftSnapshotCreate(util.VaultClient, "vault.bak"); err != nil {
		test.Error("vault raft snapshot creation for snapshot restoration test failed")
	}
	if err := VaultRaftSnapshotRestore(util.VaultClient, "vault.bak"); err != nil {
		test.Error("vault raft snapshot restoration failed")
		test.Error(err)
	}

	os.Remove("./vault.bak")

	// some weird bug in gha prevents vault server from restarting properly
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		if err := VaultRaftSnapshotRestore(util.VaultClient, "/foo/vault.bak"); err == nil || err.Error() != "open /foo/vault.bak: no such file or directory" {
			test.Errorf("expected error (contains): open /foo/vault.bak: no such file or directory, actual: %v", err)
		}
	}

	// ensure vault server is available (it is possible it has not finished restarting after restoration)
	client, err := vault.NewClient(&vault.Config{Address: util.VaultAddress})
	if err != nil {
		log.Fatalf("failed to create vault client for validating server: %s", err)
	}
	// ensure vault server is healthy after snapshot restoration during unit tests
	for i := range 16 {
		// cluster healthy and unsealed?
		if health, err := client.Sys().Health(); err == nil && !health.Sealed {
			break
		} else if i == 15 {
			// check if error
			if err != nil {
				log.Print(err)
			}
			// for some reason the server never recovered
			log.Fatalf("vault server was not available after fifteen seconds")
		}
		// otherwise wait and try again
		time.Sleep(1 * time.Second)
	}
}
