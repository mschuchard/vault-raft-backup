package vault

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"

	"github.com/mitodl/vault-raft-backup/util"
)

// vault raft snapshot creation
func VaultRaftSnapshot(client *vault.Client, snapshotPath string) (*os.File, error) {
	// prepare snapshot file
	snapshotFile, err := os.OpenFile(snapshotPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fmt.Println("snapshot file at " + snapshotPath + " could not be created")
		fmt.Println(err)
		return nil, err
	}

	// defer snapshot close
	defer util.SnapshotFileClose(snapshotFile)

	// execute raft snapshot
	err = client.Sys().RaftSnapshot(snapshotFile)
	if err != nil {
		snapshotFile.Close()
		fmt.Println("Vault Raft snapshot invocation failed")
		fmt.Println(err)
		return nil, err
	}

	return snapshotFile, nil
}
