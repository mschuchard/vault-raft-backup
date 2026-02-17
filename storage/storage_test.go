package storage

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestStorageTransfer(test *testing.T) {
	if err := StorageTransfer(&util.CloudConfig{}, &util.SnapshotConfig{Cleanup: true, Path: "/foo"}); err == nil || err.Error() != "open /foo: no such file or directory" {
		test.Error("did not return error as expected for nonexistent snapshot file")
		test.Error(err)
	}

	if err := StorageTransfer(&util.CloudConfig{Platform: "doesnotexist"}, &util.SnapshotConfig{Cleanup: false, Path: "../.gitignore"}); err == nil || err.Error() != "invalid storage platform" {
		test.Error("unexpected or no error returned")
		test.Errorf("expected: invalid storage platform, actual: %s", err)
	}
}

func TestCompressReader(test *testing.T) {
	if _, err := CompressReader(nil, 10); err != nil {
		test.Error("invalid input compression level did not reset to 1")
		test.Error(err)
	}
}
