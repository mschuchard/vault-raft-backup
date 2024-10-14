package storage

import (
	"os"
	"strings"
	"testing"
)

func TestSnapshotS3Upload(test *testing.T) {
	// test this fails at s3upload
	os.Setenv("AWS_REGION", "us-west-1")
	fooFile, err := os.Create("./foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	defer fooFile.Close()
	defer os.Remove("./foo")

	if _, err := snapshotS3Upload(object, fooFile, "prefix-foo"); err == nil || !strings.Contains(err.Error(), "NoCredentialProviders: no valid providers in chain") {
		test.Errorf("expected error (contains): NoCredentialProviders: no valid providers in chain, actual: %v", err)
	}
}
