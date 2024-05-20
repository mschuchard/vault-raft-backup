package aws

import (
	"os"
	"testing"
)

const (
	prefix = "my_prefix"
	bucket = "my_bucket"
)

var expectedConfig awsConfig = awsConfig{s3Bucket: "my_bucket", s3Prefix: "my_prefix"}

func TestNewAWSConfig(test *testing.T) {
	os.Setenv("S3_PREFIX", prefix)
	_, err := NewAWSConfig()
	if err == nil || err.Error() != "empty s3 bucket input setting" {
		test.Errorf("expected error: empty s3 bucket input setting, actual %v", err)
	}

	os.Setenv("S3_BUCKET", bucket)
	config, err := NewAWSConfig()
	if err != nil {
		test.Errorf("constructor unexpectedly errored with %v", err)
	}
	if *config != expectedConfig {
		test.Errorf("expected aws config values: %v", expectedConfig)
		test.Errorf("actual aws config value: %v", *config)
	}
}
