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
	backupAWSConfig = &util.AWSConfig{S3Prefix: prefix}
	expectedConfig  = awsConfig{s3Bucket: bucket, s3Prefix: prefix}
)

func TestNewAWSConfig(test *testing.T) {
	_, err := NewAWSConfig(backupAWSConfig)
	if err == nil || err.Error() != "empty s3 bucket input setting" {
		test.Errorf("expected error: empty s3 bucket input setting, actual %v", err)
	}

	backupAWSConfig.S3Bucket = bucket
	config, err := NewAWSConfig(backupAWSConfig)
	if err != nil {
		test.Errorf("constructor unexpectedly errored with %v", err)
	}
	if *config != expectedConfig {
		test.Errorf("expected aws config values: %v", expectedConfig)
		test.Errorf("actual aws config value: %v", *config)
	}
}
