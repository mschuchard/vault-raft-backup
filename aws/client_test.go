package aws

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

const (
	bucket = "my_bucket"
	prefix = "my_prefix"
)

var (
	backupConfig = &util.BackupConfig{
		AWSConfig: &util.AWSConfig{S3Prefix: prefix},
	}
	expectedConfig = awsConfig{s3Bucket: bucket, s3Prefix: prefix}
)

func TestNewAWSConfig(test *testing.T) {
	_, err := NewAWSConfig(backupConfig)
	if err == nil || err.Error() != "empty s3 bucket input setting" {
		test.Errorf("expected error: empty s3 bucket input setting, actual %v", err)
	}

	backupConfig.AWSConfig.S3Bucket = bucket
	config, err := NewAWSConfig(backupConfig)
	if err != nil {
		test.Errorf("constructor unexpectedly errored with %v", err)
	}
	if *config != expectedConfig {
		test.Errorf("expected aws config values: %v", expectedConfig)
		test.Errorf("actual aws config value: %v", *config)
	}
}
