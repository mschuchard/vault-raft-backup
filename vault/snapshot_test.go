package vault

import (
	"os"
	"strings"
	"testing"

	"github.com/mitodl/vault-raft-backup/util"
)

func TestVaultRaftSnapshot(test *testing.T) {
	snapshotFile, err := VaultRaftSnapshot(util.VaultClient, "vault.bak")
	if err == nil || !strings.Contains(err.Error(), "GET http://127.0.0.1:8200/v1/sys/storage/raft/snapshot") {
		test.Errorf("expected error (contains): GET http://127.0.0.1:8200/v1/sys/storage/raft/snapshot, actual: %v", err)
	}
	if _, err = snapshotFile.Stat(); err == nil {
		test.Error("vault raft snapshot file was not actually created")
	}

	os.Remove("./vault.bak")
}
