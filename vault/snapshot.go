package vault

import (
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"

	"github.com/mitodl/vault-raft-backup/util"
)

// vault raft snapshot creation
func VaultRaftSnapshot(client *vault.Client, snapshotPath string) (*os.File, error) {
	// prepare snapshot file
	snapshotFile, err := os.OpenFile(snapshotPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Print("snapshot file at " + snapshotPath + " could not be created")
		return nil, err
	}
	defer util.SnapshotFileClose(snapshotFile)

	// execute raft snapshot
	err = client.Sys().RaftSnapshot(snapshotFile)
	if err != nil {
		log.Print("Vault Raft snapshot creation failed")
		return nil, err
	}

	return snapshotFile, nil
}
