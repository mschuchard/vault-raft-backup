package aws

import (
	"os"
	"strings"
	"testing"
)

func TestSnapshotS3Upload(test *testing.T) {
	// test this fails at s3upload
	os.Setenv("AWS_REGION", "us-west-1")
	_, err := os.Create("./foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}

	if _, err := SnapshotS3Upload(&expectedConfig, "./foo", true); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}
}
