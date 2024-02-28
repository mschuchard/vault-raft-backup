package util

import (
	"log"
	"os"
)

// close snapshot file
func SnapshotFileClose(snapshotFile *os.File) error {
	// close file
	err := snapshotFile.Close()
	if err != nil {
		log.Printf("Vault Raft snapshot file at '%s' failed to close", snapshotFile.Name())
	}

	return err
}

// remove snapshot file
func SnapshotFileRemove(snapshotFile *os.File) error {
	// assign filename
	filename := snapshotFile.Name()

	// verify file existence
	if _, err := snapshotFile.Stat(); err != nil {
		log.Printf("Vault Raft snapshot file does not exist at '%s'", filename)
		return err
	}
	// remove file
	err := os.Remove(filename)
	if err != nil {
		log.Printf("failed to remove Vault Raft snapshot at '%s'", filename)
	}

	return err
}
