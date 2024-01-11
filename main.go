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
		log.Print("vault configuration failed validation")
		log.Fatal(err)
	}
	awsConfig := aws.NewAWSConfig()

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
	_, err = aws.SnapshotS3Upload(awsConfig, snapshotFile.Name())
	if err != nil {
		log.Print("S3 upload failed")
		log.Fatal(err)
	}
}
