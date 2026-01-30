package vault

import (
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"

	"github.com/mschuchard/vault-raft-backup/util"
)

// vault raft snapshot creation
func VaultRaftSnapshotCreate(client *vault.Client, snapshotPath string) error {
	// prepare snapshot file for content writing
	snapshotFile, err := os.OpenFile(snapshotPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		log.Printf("snapshot file at '%s' could not be created or opened", snapshotPath)
		return err
	}

	// defer snapshot close
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
	}()

	// execute raft snapshot creation
	err = client.Sys().RaftSnapshot(snapshotFile)
	if err != nil {
		log.Print("Vault Raft snapshot creation failed")
		return err
	}

	log.Printf("snapshot file created on local filesystem at '%s'", snapshotPath)

	return nil
}

// vault raft snapshot restoration
func VaultRaftSnapshotRestore(client *vault.Client, snapshotPath string) error {
	// prepare snapshot file for content reading
	snapshotFile, err := os.OpenFile(snapshotPath, os.O_RDONLY, 0o600)
	if err != nil {
		log.Printf("snapshot file at '%s' could not be opened for reading", snapshotPath)
		return err
	}

	// defer snapshot close
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
	}()

	// execute raft snapshot restore
	err = client.Sys().RaftSnapshotRestore(snapshotFile, false)
	if err != nil {
		log.Print("Vault Raft snapshot restore failed")
		return err
	}

	log.Printf("snapshot file restored to vault from local filesystem at '%s'", snapshotPath)

	return nil
}
