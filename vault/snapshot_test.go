package vault

import (
	"os"
	"strings"
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestVaultRaftSnapshotCreate(test *testing.T) {
	snapshotFile, err := VaultRaftSnapshotCreate(util.VaultClient, "vault.bak")
	if err == nil || !strings.Contains(err.Error(), "GET http://127.0.0.1:8200/v1/sys/storage/raft/snapshot") {
		test.Errorf("expected error (contains): GET http://127.0.0.1:8200/v1/sys/storage/raft/snapshot, actual: %v", err)
	}
	if _, err = snapshotFile.Stat(); err == nil {
		test.Error("vault raft snapshot file was not actually created")
	}

	os.Remove("./vault.bak")

	if _, err = VaultRaftSnapshotCreate(util.VaultClient, "/foo/vault.bak"); err == nil || err.Error() != "open /foo/vault.bak: no such file or directory" {
		test.Errorf("expected error (contains): open /foo/vault.bak: no such file or directory, actual: %s", err)
	}
}

func TestVaultRaftSnapshotRestore(test *testing.T) {
	err := VaultRaftSnapshotRestore(util.VaultClient, "snapshot_test.go")
	if err == nil || !strings.Contains(err.Error(), "POST http://127.0.0.1:8200/v1/sys/storage/raft/snapshot") {
		test.Errorf("expected error (contains): POST http://127.0.0.1:8200/v1/sys/storage/raft/snapshot, actual: %v", err)
	}

	if err = VaultRaftSnapshotRestore(util.VaultClient, "/foo/vault.bak"); err == nil || err.Error() != "open /foo/vault.bak: no such file or directory" {
		test.Errorf("expected error (contains): open /foo/vault.bak: no such file or directory, actual: %v", err)
	}
}
