package util

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// while these are public to decode, the individual structs initialized from this are safely private
// storage configs
type AWSConfig struct {
	S3Bucket string `hcl:"s3_bucket"`
	S3Prefix string `hcl:"s3_prefix,optional"`
}

type GCPConfig struct {
	CSBucket string `hcl:"cs_bucket"`
	CSPrefix string `hcl:"cs_prefix"`
}

// vault config
type VaultConfig struct {
	Address      string `hcl:"address,optional"`
	Insecure     bool   `hcl:"insecure,optional"`
	Engine       string `hcl:"auth_engine,optional"`
	Token        string `hcl:"token,optional"`
	AWSMountPath string `hcl:"aws_mount_path,optional"`
	AWSRole      string `hcl:"aws_role,optional"`
	SnapshotPath string `hcl:"snapshot_path,optional"`
}

// overall vault raft backup config
type BackupConfig struct {
	AWSConfig       *AWSConfig   `hcl:"aws_config,block"`
	GCPConfig       *GCPConfig   `hcl:"gcp_config,block"`
	VaultConfig     *VaultConfig `hcl:"vault_config,block"`
	SnapshotCleanup bool         `hcl:"snapshot_cleanup,optional"`
}

// config constructor
func NewBackupConfig(filePath string) (*BackupConfig, error) {
	// determine input structure and return accordingly
	if len(filePath) == 0 {
		return envImportConfig()
	} else {
		return hclDecodeConfig(filePath)
	}
}

// decode hcl config file into vault raft backup config
func hclDecodeConfig(filePath string) (*BackupConfig, error) {
	// initialize config
	var backupConfig BackupConfig
	// decode hcl config file into vault raft backup config struct
	err := hclsimple.DecodeFile(filePath, nil, &backupConfig)
	if err != nil {
		log.Printf("the provided hcl config file at %s could not be parsed into a valid config for vault raft backup", filePath)
	}

	return &backupConfig, err
}

// import environment variables into vault raft backup config
func envImportConfig() (*BackupConfig, error) {
	// import environment variables into vault raft backup config struct
	// validate vault insecure
	insecureEnv := os.Getenv("VAULT_SKIP_VERIFY")
	insecure, err := strconv.ParseBool(insecureEnv)
	if err != nil && len(insecureEnv) > 0 { // assigned value could not be converted to boolean
		log.Print(err)
		log.Printf("invalid boolean value '%s' for VAULT_SKIP_VERIFY", insecureEnv)
		return nil, errors.New("invalid VAULT_SKIP_VERIFY value")
	}
	// validate snapshot cleanup
	cleanupEnv := os.Getenv("SNAPSHOT_CLEANUP")
	cleanup, err := strconv.ParseBool(cleanupEnv)
	if err != nil && len(cleanupEnv) > 0 { // assigned value could not be converted to boolean
		log.Print(err)
		log.Printf("invalid boolean value '%s' for SNAPSHOT_CLEANUP", cleanupEnv)
		return nil, errors.New("invalid SNAPSHOT_CLEANUP value")
	}

	return &BackupConfig{
		AWSConfig: &AWSConfig{
			S3Bucket: os.Getenv("S3_BUCKET"),
			S3Prefix: os.Getenv("S3_PREFIX"),
		},
		GCPConfig: &GCPConfig{
			CSBucket: os.Getenv("CS_BUCKET"),
			CSPrefix: os.Getenv("CS_PREFIX"),
		},
		VaultConfig: &VaultConfig{
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
