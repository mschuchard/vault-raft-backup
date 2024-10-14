package storage

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

const (
	object = "my_bucket"
	prefix = "my_prefix"
)

var (
	backupAWSConfig = &util.AWSConfig{S3Prefix: prefix}
	expectedConfig  = config{object: object, prefix: prefix}
)

func TestNewConfig(test *testing.T) {
	_, err := NewConfig(backupAWSConfig)
	if err == nil || err.Error() != "empty storage object input setting" {
		test.Errorf("expected error: empty storage object input setting, actual %s", err)
	}

	backupAWSConfig.S3Bucket = object
	config, err := NewConfig(backupAWSConfig)
	if err != nil {
		test.Errorf("constructor unexpectedly errored with %v", err)
	}
	if *config != expectedConfig {
		test.Errorf("expected aws config values: %v", expectedConfig)
		test.Errorf("actual aws config value: %v", *config)
	}
}
