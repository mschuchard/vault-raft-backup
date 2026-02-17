package storage

import (
	"os"
	"testing"

	"github.com/mschuchard/vault-raft-backup/enum"
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

	// create empty file for testing
	if err := os.WriteFile("test.bak", []byte(""), 0644); err != nil {
		test.Error("failed to create test file for local end-to-end compression test")
		test.Error(err)
	}
	// end-to-end test local storage compressed snapshot
	if err := StorageTransfer(&util.CloudConfig{Container: "/tmp", Platform: enum.LOCAL, Prefix: "test-"}, &util.SnapshotConfig{Cleanup: true, Path: "test.bak", CompressionLevel: 1}); err != nil {
		test.Error(err)
	}
	// verify compressed file exists at destination and has expected size for gzip level 1 empty file
	if fileInfo, err := os.Stat("/tmp/test-test.bak"); err != nil || fileInfo.Size() != 23 {
		test.Error("compressed snapshot file was not found at expected destination")
		test.Errorf("expected file size: 23, actual file size: %d", fileInfo.Size())
		test.Error(err)
	}
}

func TestCompressReader(test *testing.T) {
	if _, err := CompressReader(nil, 10); err != nil {
		test.Error("invalid input compression level did not reset to 1")
		test.Error(err)
	}
}
