package util

import (
	"log"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type VaultConfig struct {
	Address      string `hcl:"address"`
	Insecure     bool   `hcl:"insecure"`
	Engine       string `hcl:"auth_engine"`
	Token        string `hcl:"token"`
	AWSMountPath string `hcl:"aws_mount_path"`
	AWSRole      string `hcl:"aws_role"`
	SnapshotPath string `hcl:"snapshot_path"`
}

type AWSConfig struct {
	S3Bucket string `hcl:"s3_bucket"`
	S3Prefix string `hcl:"s3_prefix"`
}

type Config struct {
	AWSConfig       *AWSConfig   `hcl:"aws_config,block"`
	VaultConfig     *VaultConfig `hcl:"vault_config"`
	SnapshotCleanup bool         `hcl:"snapshot_cleanup"`
}

func HclDecodeConfig(filePath string) (*Config, error) {
	// initialize config
	config := &Config{}
	// decode hcl config file into vault raft backup config struct
	err := hclsimple.DecodeFile(filePath, nil, config)
	if err != nil {
		log.Printf("the provided hcl config file at %s could not be parsed into a valid config for vault raft backup", filePath)
	}

	return config, err
}
