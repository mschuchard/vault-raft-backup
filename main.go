package main

import (
	"log"

	"github.com/mitodl/vault-raft-backup/aws"
	"github.com/mitodl/vault-raft-backup/vault"
)

func main() {
	// construct vault client config and aws client config
	vaultConfig, err := vault.NewVaultConfig()
	if err != nil {
		log.Print("Vault configuration failed validation")
		log.Fatal(err)
	}
	awsConfig, err := aws.NewAWSConfig()
	if err != nil {
		log.Print("AWS configuration failed validation")
		log.Fatal(err)
	}

	// initialize and configure vault client
	vaultClient, err := vault.NewVaultClient(vaultConfig)
	if err != nil {
		log.Print("Vault client initialization and configuration failed")
		log.Fatal(err)
	}

	// vault raft snapshot
	snapshotFile, err := vault.VaultRaftSnapshot(vaultClient, vaultConfig.SnapshotPath())
	if err != nil {
		log.Print("Vault Raft snapshot failed")
		log.Fatal(err)
	}

	// upload snapshot to aws s3
	_, err = aws.SnapshotS3Upload(awsConfig, snapshotFile.Name(), true)
	if err != nil && err.Error() != "snapshot not found" && err.Error() != "snapshot not removed" {
		// not an error from failed removal so error is actually fatal
		log.Print("S3 upload failed")
		log.Fatal(err)
	}
}
