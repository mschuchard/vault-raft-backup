package aws

import (
	"os"
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
	_, err := os.Create("foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	awsConfig := &AWSConfig{s3Bucket: "my_bucket"}

	if _, err := SnapshotS3Upload(awsConfig, "foo"); err == nil || err.Error() != "MissingRegion: could not find region configuration" {
		test.Errorf("expected error: MissingRegion: could not find region configuration, actual: %v", err)
	}
}
