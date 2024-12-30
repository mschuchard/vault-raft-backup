package storage

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestStorageTransfer(test *testing.T) {
	if err := StorageTransfer(&util.CloudConfig{}, "/foo", true); err == nil {
		test.Error("did not return error as expected for nonexistent snapshot file")
	}

	if err := StorageTransfer(&util.CloudConfig{Platform: "doesnotexist"}, "../.gitignore", true); err == nil || err.Error() != "invalid cloud platform" {
		test.Error("unexpected or no error returned")
		test.Errorf("expected: invalid cloud platform, actual: %s", err)
	}
}
