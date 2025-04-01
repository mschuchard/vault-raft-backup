package util

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/mschuchard/vault-raft-backup/enum"
)

// while these are public to decode, the individual structs initialized from this are safely private
// storage configs
type CloudConfig struct {
	Container string        `hcl:"container"`
	Platform  enum.Platform `hcl:"platform"`
	Prefix    string        `hcl:"prefix,optional"`
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
	err := hclsimple.DecodeFile(filePath, nil, &backupConfig)
	if err != nil {
		log.Printf("the provided hcl config file at %s could not be parsed into a valid config for vault raft backup", filePath)
		return nil, err
	}

	// validate a cloud config block was specified
	if backupConfig.CloudConfig == nil {
		log.Print("the cloud_config block is required in the input configuration file")
		return nil, errors.New("cloud_config block absent")
	}

	// validate parameters and finalize snapshot path
	if _, err = enum.Platform(backupConfig.CloudConfig.Platform).New(); err != nil {
		return nil, err
	}

	backupConfig.VaultConfig.SnapshotPath, err = defaultSnapshotPath(backupConfig.VaultConfig.SnapshotPath)
	if err != nil {
		return nil, err
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
		log.Printf("invalid boolean value '%s' for VAULT_SKIP_VERIFY", insecureEnv)
		log.Print(err)
		return nil, errors.New("invalid VAULT_SKIP_VERIFY value")
	}

	// validate snapshot cleanup
	cleanupEnv := os.Getenv("SNAPSHOT_CLEANUP")
	cleanup, err := strconv.ParseBool(cleanupEnv)
	if err != nil && len(cleanupEnv) > 0 { // assigned value could not be converted to boolean
		log.Printf("invalid boolean value '%s' for SNAPSHOT_CLEANUP", cleanupEnv)
		log.Print(err)
		return nil, errors.New("invalid SNAPSHOT_CLEANUP value")
	}

	// validate container and platform were specified
	container := os.Getenv("CONTAINER")
	if len(container) == 0 {
		log.Print("CONTAINER is a required input value, and it was unspecified as an environment variable")
		return nil, errors.New("environment variable absent")
	}
	platform, err := enum.Platform(os.Getenv("PLATFORM")).New()
	if err != nil {
		return nil, err
	}

	// validate parameters and finalize snapshot path
	snapshotPath, err := defaultSnapshotPath(os.Getenv("VAULT_SNAPSHOT_PATH"))
	if err != nil {
		return nil, err
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
			SnapshotPath: snapshotPath,
		},
		SnapshotCleanup: cleanup,
	}, nil
}

// determines default snapshot path
func defaultSnapshotPath(snapshotPath string) (string, error) {
	// provide snapshot path default if unspecified
	if len(snapshotPath) == 0 {
		// create timestamp for default filename suffix
		timestamp := time.Now().Local().Format("2006-01-02-150405")
		defaultFilename := "vault-" + timestamp + "-*.bak"

		// create random tmp file in tmp dir and then close it for later backup
		snapshotTmpFile, err := os.CreateTemp(os.TempDir(), defaultFilename)
		if err != nil {
			log.Printf("could not create a temporary file for the local snapshot file in the temporary directory '%s'", os.TempDir())
			return "", err
		}
		snapshotTmpFile.Close()

		// assign to snapshot path config field member
		snapshotPath = snapshotTmpFile.Name()
		log.Printf("vault raft snapshot path defaulting to '%s'", snapshotPath)
	}

	return snapshotPath, nil
}
