package storage

import (
	"os"
	"strings"
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestSnapshotS3Upload(test *testing.T) {
	// test this fails at s3upload
	fooFile, err := os.Create("./foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	defer fooFile.Close()
	defer os.Remove("./foo")

	if err := snapshotS3Upload(util.Container, fooFile, "prefix-foo"); err == nil || !strings.Contains(err.Error(), "get credentials: failed to refresh cached credentials") {
		test.Errorf("expected error (contains): get credentials: failed to refresh cached credentials, actual: %v", err)
	}
}
