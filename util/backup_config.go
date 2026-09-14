package util

import (
	"errors"
	"log"
	"os"
	"regexp"
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
	Address    string          `hcl:"address,optional"`
	Insecure   bool            `hcl:"insecure,optional"`
	Engine     enum.AuthEngine `hcl:"auth_engine,optional"`
	Token      string          `hcl:"token,optional"`
	SecretID   string          `hcl:"secret_id"`
	WrapToken  string          `hcl:"wrap_token"`
	AzResource string          `hcl:"az_resource"`
	AuthMount  string          `hcl:"auth_mount,optional"`
	VaultRole  string          `hcl:"vault_role,optional"`
	Namespace  string          `hcl:"namespace,optional"`
}

// snapshot config
type SnapshotConfig struct {
	Cleanup          bool   `hcl:"cleanup,optional"`
	CompressionLevel int    `hcl:"compression_level,optional"`
	Path             string `hcl:"path,optional"`
	Restore          bool   `hcl:"restore,optional"`
}

// overall vault raft backup config
type BackupConfig struct {
	CloudConfig    *CloudConfig    `hcl:"cloud_config,block"`
	VaultConfig    *VaultConfig    `hcl:"vault_config,block"`
	SnapshotConfig *SnapshotConfig `hcl:"snapshot_config,block"`
}

// config constructor
func NewBackupConfig(filePath string) (*BackupConfig, error) {
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
	if err = validateParams(backupConfig.CloudConfig.Platform, backupConfig.VaultConfig.Engine, backupConfig.CloudConfig.AZAccountURL, backupConfig.SnapshotConfig); err != nil {
		return nil, err
	}

	// finalize snapshot path
	if backupConfig.SnapshotConfig.Path, err = defaultSnapshotPath(backupConfig.SnapshotConfig.Path); err != nil {
		return nil, err
	}

	return &backupConfig, nil
}

// validates various input parameters
func validateParams(platform enum.Platform, authEngine enum.AuthEngine, azAccountURL string, snapshotConfig *SnapshotConfig) error {
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

	// validate snapshot config params if defined
	if snapshotConfig != nil {
		// validate params for restoration scenario
		if snapshotConfig.Restore {
			// validate restore and cleanup are not both true
			if snapshotConfig.Cleanup {
				log.Print("snapshot cleanup is specified as 'true', but this has no effect since restore is also specified as 'true'")
				log.Print("ignoring cleanup parameter value in restoration scenario")
			}
		} else { // validate params for backup scenario
			// validate compression level is between 0 and 3 inclusive
			if snapshotConfig.CompressionLevel < 0 || snapshotConfig.CompressionLevel > 3 {
				log.Printf("snapshot compression level must be an integer between 0 and 3 inclusive, but instead %d was specified", snapshotConfig.CompressionLevel)
				return errors.New("invalid snapshot compression level")
			}
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
