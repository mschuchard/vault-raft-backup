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
	backupConfig, err := util.NewBackupConfig(*hclConfigPath)
	if err != nil {
		log.Print("vault raft backup configuration failed validation")
		log.Fatal(err)
	}

	// construct vault client config
	vaultConfig, err := vault.NewVaultConfig(backupConfig.VaultConfig)
	if err != nil {
		log.Print("vault configuration failed validation")
		log.Fatal(err)
	}

	// initialize and configure vault client
	vaultClient, err := vault.NewVaultClient(vaultConfig)
	if err != nil {
		log.Print("vault client initialization and configuration failed")
		log.Fatal(err)
	}

	// vault raft snapshot
	snapshotFile, err := vault.VaultRaftSnapshot(vaultClient, backupConfig.VaultConfig.SnapshotPath)
	if err != nil {
		log.Print("vault raft snapshot failed")
		log.Fatal(err)
	}

	// transfer snapshot to cloud storage
	err = storage.StorageTransfer(backupConfig.CloudConfig, snapshotFile.Name(), backupConfig.SnapshotCleanup)
	if err != nil && err.Error() != "snapshot not found" && err.Error() != "snapshot not removed" {
		// not an error from failed removal so error is actually fatal
		log.Print("cloud storage upload failed")
		log.Fatal(err)
	}
}
