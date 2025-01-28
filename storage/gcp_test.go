package storage

import (
	"os"
	"testing"

	"github.com/mschuchard/vault-raft-backup/util"
)

func TestSnapshotCSUpload(test *testing.T) {
	// test this fails at upload transfer
	fooFile, err := os.Open("../.gitignore")
	if err != nil {
		test.Error("test short-circuited because file could not be opened")
		return
	}
	defer fooFile.Close()

	if err := snapshotCSUpload(util.Container, fooFile, ""); err == nil || err.Error() != "dialing: google: could not find default credentials. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information" {
		test.Errorf("expected error: dialing: google: could not find default credentials. See https://cloud.google.com/docs/authentication/external/set-up-adc for more information, actual: %v", err)
	}
}
