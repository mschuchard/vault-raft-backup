package storage

import (
	"os"
	"strings"
	"testing"

	"github.com/mschuchard/vault-raft-backup/testutil"
)

func TestSnapshotBlobUpload(test *testing.T) {
	// test this fails at azurecredential
	fooFile, err := os.Open("../.gitignore")
	if err != nil {
		test.Error("test short-circuited because file could not be opened")
		return
	}
	defer fooFile.Close()

	if err := snapshotBlobUpload(testutil.Container, fooFile, "empty", "https://foo.com"); err == nil || !strings.Contains(err.Error(), "DefaultAzureCredential: failed to acquire a token.") {
		test.Errorf("expected error (contains): DefaultAzureCredential: failed to acquire a token., actual: %v", err)
	}
}
