package aws

import (
	"os"
	"strings"
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

func TestSnapshotS3Upload(test *testing.T) {
	// test this fails at s3upload
	os.Setenv("AWS_REGION", "us-west-1")
	_, err := os.Create("foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}

	if _, err := SnapshotS3Upload(&expectedConfig, "foo"); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}
}
