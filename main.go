package main

import (
	"log"

	"github.com/mitodl/vault-raft-backup/aws"
	"github.com/mitodl/vault-raft-backup/vault"
)

func main() {
	// construct vaultconfig and awsConfig
	vaultConfig := vault.NewVaultConfig()
	awsConfig := aws.NewAWSConfig()

	// initialize and configure client
	vaultClient, err := vault.VaultClient(vaultConfig)
	if err != nil {
		log.Fatalln("Vault client initialization and configuration failed")
	}

	// vault raft snapshot
	snapshotFile, err := vault.VaultRaftSnapshot(vaultClient, vaultConfig.SnapshotPath())
	if err != nil {
		log.Fatalln("Vault Raft Snapshot failed")
	}

	// execute snapshot upload to s3
	_, err = aws.SnapshotS3Upload(awsConfig, snapshotFile.Name())
	if err != nil {
		log.Fatalln("S3 upload failed")
	}
}
