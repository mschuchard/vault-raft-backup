package util

import (
	"log"
	"os"
)

// close snapshot file
func SnapshotFileClose(snapshotFile *os.File) {
	// close file
	err := snapshotFile.Close()
	if err != nil {
		log.Print("Vault raft snapshot file failed to close")
		log.Fatal(err)
	}
}
