package util

import (
	"log"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type VaultConfig struct {
	Address      string `hcl:"address,optional"`
	Insecure     bool   `hcl:"insecure,optional"`
	Engine       string `hcl:"auth_engine,optional"`
	Token        string `hcl:"token,optional"`
	AWSMountPath string `hcl:"aws_mount_path,optional"`
	AWSRole      string `hcl:"aws_role,optional"`
	SnapshotPath string `hcl:"snapshot_path,optional"`
}

type AWSConfig struct {
	S3Bucket string `hcl:"s3_bucket"`
	S3Prefix string `hcl:"s3_prefix,optional"`
}

type Config struct {
	AWSConfig       *AWSConfig   `hcl:"aws_config,block"`
	VaultConfig     *VaultConfig `hcl:"vault_config,block"`
	SnapshotCleanup bool         `hcl:"snapshot_cleanup,optional"`
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
