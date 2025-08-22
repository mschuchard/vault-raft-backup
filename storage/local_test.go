package storage

import (
	"os"
	"testing"
)

func TestSnapshotFSCopy(test *testing.T) {
	fooFile, err := os.Open("../.gitignore")
	if err != nil {
		test.Error("test short-circuited because file could not be opened")
		return
	}
	defer fooFile.Close()

	if err := snapshotFSCopy("../.gitignore", fooFile, "foo"); err == nil || err.Error() != "open ../.gitignore: not a directory" {
		test.Error("snapshot copy did not fail on nonexistent destination directory")
		test.Error(err)
	}

	if err := snapshotFSCopy("/tmp", fooFile, "forbidden/char"); err == nil || err.Error() != "open /tmp/forbidden/char: no such file or directory" {
		test.Error("snapshot copy did not fail on unsuitable destination target file")
		test.Error(err)
	}

	if err := snapshotFSCopy("/tmp", fooFile, "foo"); err != nil {
		test.Error("snapshot copy failed to write to destination at /tmp/foo")
		test.Error(err)
	}
	copiedFile, err := os.ReadFile("/tmp/foo")
	if err != nil {
		test.Error("copied file could not be opened for additional validation")
		test.Error(err)
		return
	}
	fooFileContent, err := os.ReadFile("../.gitignore")
	if err != nil {
		test.Error("original file could not be read")
		test.Error(err)
	}
	if string(copiedFile) != string(fooFileContent) {
		test.Error("copied file did not contain same contents as original file")
		test.Errorf("original: %s, copied: %s", string(fooFileContent), string(copiedFile))
	}
}
