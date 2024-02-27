package util

import (
	"log"
	"os"
)

// close snapshot file (intended for deferral)
func SnapshotFileClose(snapshotFile *os.File) error {
	// close file
	err := snapshotFile.Close()
	if err != nil {
		log.Print("Vault raft snapshot file failed to close")
	}
	return err
}
