package util

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type AWSConfig struct {
	S3Bucket string `hcl:"s3_bucket"`
	S3Prefix string `hcl:"s3_prefix,optional"`
}

type VaultConfig struct {
	Address      string `hcl:"address,optional"`
	Insecure     bool   `hcl:"insecure,optional"`
	Engine       string `hcl:"auth_engine,optional"`
	Token        string `hcl:"token,optional"`
	AWSMountPath string `hcl:"aws_mount_path,optional"`
	AWSRole      string `hcl:"aws_role,optional"`
	SnapshotPath string `hcl:"snapshot_path,optional"`
}

type BackupConfig struct {
	AWSConfig       *AWSConfig   `hcl:"aws_config,block"`
	VaultConfig     *VaultConfig `hcl:"vault_config,block"`
	SnapshotCleanup bool         `hcl:"snapshot_cleanup,optional"`
}

// decode hcl config file into vault raft backup config
func HclDecodeConfig(filePath string) (*BackupConfig, error) {
	// initialize config
	var backupConfig *BackupConfig
	// decode hcl config file into vault raft backup config struct
	err := hclsimple.DecodeFile(filePath, nil, backupConfig)
	if err != nil {
		log.Printf("the provided hcl config file at %s could not be parsed into a valid config for vault raft backup", filePath)
	}

	return backupConfig, err
}

// import environment variables into vault raft backup config
func OSImportConfig() (*BackupConfig, error) {
	// import environment variables into vault raft backup config struct
	// validate vault insecure
	insecure, err := strconv.ParseBool(os.Getenv("VAULT_SKIP_VERIFY"))
	if err != nil { // assigned value could not be converted to boolean
		log.Printf("invalid boolean value '%s' for VAULT_SKIP_VERIFY", os.Getenv("VAULT_SKIP_VERIFY"))
		return nil, errors.New("invalid VAULT_SKIP_VERIFY value")
	}
	// validate snapshot cleanup
	cleanup, err := strconv.ParseBool(os.Getenv("SNAPSHOT_CLEANUP"))
	if err != nil {
		log.Printf("invalid boolean value '%s' for SNAPSHOT_CLEANUP", os.Getenv("SNAPSHOT_CLEANUP"))
		return nil, err
	}

	return &BackupConfig{
		AWSConfig: &AWSConfig{
			S3Bucket: os.Getenv("S3_BUCKET"),
			S3Prefix: os.Getenv("S3_PREFIX"),
		},
		VaultConfig:     &VaultConfig{
			Address:      os.Getenv("VAULT_ADDR"),
			Insecure:     insecure,
			Engine:       os.Getenv("VAULT_AUTH_ENGINE"),
			Token:        os.Getenv("VAULT_TOKEN"),
			AWSMountPath: os.Getenv("VAULT_AWS_MOUNT"),
			AWSRole:      os.Getenv("VAULT_AWS_ROLE"),
			SnapshotPath: os.Getenv("VAULT_SNAPSHOT_PATH"),
		},
		SnapshotCleanup: cleanup,
	}, nil
}
