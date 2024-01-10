package util

import (
	"os"
	"testing"
)

func TestSnapshotFileClose(test *testing.T) {
	genericFile, err := os.OpenFile("foo", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		test.Error("test short-circuited because file could not be created and opened")
	}
	SnapshotFileClose(genericFile)
}
