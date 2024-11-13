package storage

import (
	"os"
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestSnapshotCSUpload(test *testing.T) {
	// test this fails at upload transfer
	fooFile, err := os.Create("./foo")
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	defer fooFile.Close()
	defer os.Remove("./foo")

	if err := snapshotCSUpload(util.Container, fooFile, "prefix-foo"); err == nil || err.Error() != "dialing: google: could not find default credentials. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information" {
		test.Errorf("expected error: dialing: google: could not find default credentials. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information, actual: %v", err)
	}
}
