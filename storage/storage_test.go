package storage

import (
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

const object = "my_bucket"

func TestStorageTransfer(test *testing.T) {
	if err := StorageTransfer(&util.CloudConfig{}, "/foo", true); err == nil {
		test.Error("did not return error as expected for nonexistent snapshot file")
	}
}
