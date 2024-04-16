package aws

import (
	"os"
	"strings"
	"testing"
)

func TestNewAWSConfig(test *testing.T) {
	os.Setenv("S3_PREFIX", "my_prefix")
	_, err := NewAWSConfig()
	if err == nil || err.Error() != "empty s3 bucket input setting" {
		test.Errorf("expected error: empty s3 bucket input setting, actual %v", err)
	}

	os.Setenv("S3_BUCKET", "my_bucket")
	awsConfig, err := NewAWSConfig()
	if err != nil {
		test.Errorf("constructor unexpectedly errored with %v", err)
	}
	if awsConfig.s3Bucket != "my_bucket" || awsConfig.s3Prefix != "my_prefix" {
		test.Errorf("expected bucket value: my_bucket, actual: %s", awsConfig.s3Bucket)
		test.Errorf("expected prefix value: my_prefix, actual: %s", awsConfig.s3Prefix)
	}
}

func TestSnapshotS3Upload(test *testing.T) {
	// test this fails at s3upload
	os.Setenv("AWS_REGION", "us-west-1")
	_, err := os.Create("foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	awsConfig := &awsConfig{s3Bucket: "my_bucket"}

	if _, err := SnapshotS3Upload(awsConfig, "foo"); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}
}
