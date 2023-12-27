package util

import (
	"fmt"
	"log"
	"os"
)

// close snapshot file
func SnapshotFileClose(snapshotFile *os.File) {
	// close file
	err := snapshotFile.Close()
	if err != nil {
		fmt.Println("Vault raft snapshot file failed to close")
		log.Fatalln(err)
	}
}
