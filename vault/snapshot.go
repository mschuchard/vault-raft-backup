package vault

import (
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"

	"github.com/mitodl/vault-raft-backup/util"
)

// vault raft snapshot creation
func VaultRaftSnapshot(client *vault.Client, snapshotPath string) (*os.File, error) {
	// prepare snapshot file for content writing
	snapshotFile, err := os.OpenFile(snapshotPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		log.Printf("snapshot file at %s could not be created or opened", snapshotPath)
		return nil, err
	}

	// defer snapshot close
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
	}()

	// execute raft snapshot to file
	err = client.Sys().RaftSnapshot(snapshotFile)
	if err != nil {
		log.Print("Vault Raft snapshot creation failed")
		return nil, err
	}

	return snapshotFile, err
}
