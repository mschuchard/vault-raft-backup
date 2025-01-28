package storage

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestStorageTransfer(test *testing.T) {
	if err := StorageTransfer(&util.CloudConfig{}, "/foo", true); err == nil || err.Error() != "open /foo: no such file or directory" {
		test.Error("did not return error as expected for nonexistent snapshot file")
		test.Error(err)
	}

	if err := StorageTransfer(&util.CloudConfig{Platform: "doesnotexist"}, "../.gitignore", false); err == nil || err.Error() != "invalid cloud platform" {
		test.Error("unexpected or no error returned")
		test.Errorf("expected: invalid cloud platform, actual: %s", err)
	}
}
