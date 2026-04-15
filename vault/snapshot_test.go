package vault

import (
	"os"
	"testing"

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

	if err := VaultRaftSnapshotRestore(util.VaultClient, "/foo/vault.bak"); err == nil || err.Error() != "open /foo/vault.bak: no such file or directory" {
		test.Errorf("expected error (contains): open /foo/vault.bak: no such file or directory, actual: %v", err)
	}
}
