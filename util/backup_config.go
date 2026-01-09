package util

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/mschuchard/vault-raft-backup/enum"
)

// while these are public to decode, the individual structs initialized from this are safely private
// storage configs
type CloudConfig struct {
	AZAccountURL string        `hcl:"az_account_url,optional"`
	Container    string        `hcl:"container"`
	Platform     enum.Platform `hcl:"platform"`
	Prefix       string        `hcl:"prefix,optional"`
}

// vault config
type VaultConfig struct {
	Address      string          `hcl:"address,optional"`
	Insecure     bool            `hcl:"insecure,optional"`
	Engine       enum.AuthEngine `hcl:"auth_engine,optional"`
	Token        string          `hcl:"token,optional"`
	AWSMountPath string          `hcl:"aws_mount_path,optional"`
	AWSRole      string          `hcl:"aws_role,optional"`
	SnapshotPath string          `hcl:"snapshot_path,optional"`
}

// snapshot config
type SnapshotConfig struct {
	Cleanup bool `hcl:"cleanup,optional"`
	Restore bool `hcl:"restore,optional"`
}

// overall vault raft backup config
type BackupConfig struct {
	CloudConfig    *CloudConfig    `hcl:"cloud_config,block"`
	VaultConfig    *VaultConfig    `hcl:"vault_config,block"`
	SnapshotConfig *SnapshotConfig `hcl:"snapshot_config,block"`
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

	// validate params
	if err = validateParams(backupConfig.CloudConfig.Platform, backupConfig.VaultConfig.Engine, backupConfig.CloudConfig.AZAccountURL); err != nil {
		return nil, err
	}

	// finalize snapshot path
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
	if err != nil { // assigned value could not be converted to boolean
		log.Printf("invalid boolean value '%s' for VAULT_SKIP_VERIFY", insecureEnv)
		log.Print(err)
		return nil, errors.New("invalid VAULT_SKIP_VERIFY value")
	}

	// validate snapshot cleanup and restore
	cleanupEnv := os.Getenv("SNAPSHOT_CLEANUP")
	cleanup, err := strconv.ParseBool(cleanupEnv)
	if err != nil { // assigned value could not be converted to boolean
		log.Printf("invalid boolean value '%s' for SNAPSHOT_CLEANUP", cleanupEnv)
		log.Print(err)
		return nil, errors.New("invalid SNAPSHOT_CLEANUP value")
	}

	restoreEnv := os.Getenv("SNAPSHOT_RESTORE")
	restore, err := strconv.ParseBool(restoreEnv)
	if err != nil { // assigned value could not be converted to boolean
		log.Printf("invalid boolean value '%s' for SNAPSHOT_RESTORE", restoreEnv)
		log.Print(err)
		return nil, errors.New("invalid SNAPSHOT_RESTORE value")
	}

	// validate container and platform were specified, and platform value
	container := os.Getenv("CONTAINER")
	if len(container) == 0 {
		log.Print("CONTAINER is a required input value, and it was unspecified as an environment variable")
		return nil, errors.New("container environment variable absent")
	}

	// validate params
	platform := enum.Platform(os.Getenv("PLATFORM"))
	authEngine := enum.AuthEngine(os.Getenv("VAULT_AUTH_ENGINE"))
	azAccountURL := os.Getenv("AZ_ACCOUNT_URL")
	if err = validateParams(platform, authEngine, azAccountURL); err != nil {
		return nil, err
	}

	// finalize snapshot path
	snapshotPath, err := defaultSnapshotPath(os.Getenv("VAULT_SNAPSHOT_PATH"))
	if err != nil {
		return nil, err
	}

	return &BackupConfig{
		CloudConfig: &CloudConfig{
			AZAccountURL: azAccountURL,
			Container:    container,
			Platform:     platform,
			Prefix:       os.Getenv("PREFIX"),
		},
		VaultConfig: &VaultConfig{
			Address:      os.Getenv("VAULT_ADDR"),
			Insecure:     insecure,
			Engine:       authEngine,
			Token:        os.Getenv("VAULT_TOKEN"),
			AWSMountPath: os.Getenv("VAULT_AWS_MOUNT"),
			AWSRole:      os.Getenv("VAULT_AWS_ROLE"),
			SnapshotPath: snapshotPath,
		},
		SnapshotConfig: &SnapshotConfig{
			Cleanup: cleanup,
			Restore: restore,
		},
	}, nil
}

// validates various input parameters
func validateParams(platform enum.Platform, authEngine enum.AuthEngine, azAccountURL string) error {
	// validate platform
	if _, err := platform.New(); err != nil {
		return err
	}

	// validate auth engine
	if _, err := authEngine.New(); err != nil {
		return err
	}

	// validate azure account url
	if platform == enum.AZ {
		if len(azAccountURL) == 0 {
			log.Print("azure specified as cloud platform, but co-requisite account url parameter was not specified")
			return errors.New("az_account_url value absent")
		} else if match, _ := regexp.MatchString(`https://.*\.blob\.core\.windows\.net`, azAccountURL); !match {
			log.Print("the azure account url must be of the form: https://<storage-account-name>.blob.core.windows.net")
			return errors.New("invalid az_account_url value")
		}
	}

	return nil
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
