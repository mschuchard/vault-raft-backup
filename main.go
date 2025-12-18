package main

import (
	"log"

	"github.com/mschuchard/vault-raft-backup/storage"
	"github.com/mschuchard/vault-raft-backup/util"
	"github.com/mschuchard/vault-raft-backup/vault"
)

func main() {
	// invoke cli parsing
	hclConfigPath := util.Cli()

	// construct vault raft backup config
	backupConfig, err := util.NewBackupConfig(hclConfigPath)
	if err != nil {
		log.Print("vault raft backup configuration failed validation")
		log.Fatal(err)
	}

	// construct configured vault validated client
	vaultClient, err := vault.NewVaultClient(backupConfig.VaultConfig)
	if err != nil {
		log.Print("configured vault client validated construction failed")
		log.Fatal(err)
	}

	// vault raft snapshot
	if backupConfig.SnapshotRestore {
		// restore from snapshot
		if err = vault.VaultRaftSnapshotRestore(vaultClient, backupConfig.VaultConfig.SnapshotPath); err != nil {
			log.Print("vault raft snapshot restore failed")
			log.Fatal(err)
		}
	} else {
		// create snapshot
		snapshotFile, err := vault.VaultRaftSnapshotCreate(vaultClient, backupConfig.VaultConfig.SnapshotPath)
		if err != nil {
			log.Print("vault raft snapshot creation failed")
			log.Fatal(err)
		}

		// transfer snapshot to cloud storage
		if err = storage.StorageTransfer(backupConfig.CloudConfig, snapshotFile.Name(), backupConfig.SnapshotCleanup); err != nil {
			if err.Error() == "snapshot not found" || err.Error() == "snapshot not removed" {
				// log the non-fatal error
				log.Print("cloud storage upload succeeded, but snapshot cleanup failed")
				log.Print(err)
			} else {
				// not an error from failed removal so error is actually fatal
				log.Print("cloud storage upload failed")
				log.Fatal(err)
			}
		}
	}
}
