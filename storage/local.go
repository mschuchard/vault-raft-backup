package storage

import (
	"io"
	"log"
	"os"

	"github.com/mschuchard/vault-raft-backup/util"
)

// snapshot copy to local filesystem
func snapshotFSCopy(directory string, snapshotFile io.Reader, snapshotName string) error {
	// validate destination directory
	if _, err := os.ReadDir(directory); err != nil {
		log.Printf("the destination directory at %s is unsuitable for copying the snapshot file", directory)
		return err
	}

	// open output file
	destination := directory + snapshotName
	destinationWriter, err := os.Open(destination)
	if err != nil {
		log.Printf("a destination file at %s could not be opened for streaming", destination)
		return err
	}

	// defer snapshot destination close
	defer func() {
		err = util.SnapshotFileClose(destinationWriter)
	}()

	// copy snapshot to destination
	if _, err = io.Copy(destinationWriter, snapshotFile); err != nil {
		log.Printf("the snapshot file at %s could not be copied to the destination at %s", snapshotFile, destination)
		return err
	}
	if err := destinationWriter.Sync(); err != nil {
		log.Printf("the snapshot file at %s could not be copied to the destination at %s", snapshotFile, destination)
		return err
	}

	log.Printf("snapshotfile successfully copied to %s", destination)
	return nil
}
