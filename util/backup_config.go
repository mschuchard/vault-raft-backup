package util

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// platform with pseudo-enum
type platform string

const (
	AWS platform = "aws"
	GCP platform = "gcp"
)

// while these are public to decode, the individual structs initialized from this are safely private
// storage configs
type CloudConfig struct {
	Container string   `hcl:"container"`
	Platform  platform `hcl:"platform"`
	Prefix    string   `hcl:"prefix,optional"`
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
	CloudConfig     *CloudConfig `hcl:"cloud_config,block"`
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
	if err := hclsimple.DecodeFile(filePath, nil, &backupConfig); err != nil {
		log.Printf("the provided hcl config file at %s could not be parsed into a valid config for vault raft backup", filePath)
		return nil, err
	}

	// validate a cloud config block was specified
	if backupConfig.CloudConfig == nil {
		log.Print("the cloud_config block is required in the input configuration file")
		return nil, errors.New("cloud_config block absent")
	}

	return &backupConfig, nil
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

	// validate container and platform were specified
	container := os.Getenv("CONTAINER")
	platform := platform(os.Getenv("PLATFORM"))
	if len(container) == 0 || len(platform) == 0 {
		log.Print("CONTAINER and PLATFORM are both required input values, and one or both was unspecified as an environment variable")
		return nil, errors.New("environment variable absent")
	}

	return &BackupConfig{
		CloudConfig: &CloudConfig{
			Container: container,
			Platform:  platform,
			Prefix:    os.Getenv("PREFIX"),
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
