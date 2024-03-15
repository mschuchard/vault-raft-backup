package util

import (
	"errors"
	"log"
	"os"
)

// close snapshot file
func SnapshotFileClose(snapshotFile *os.File) error {
	// close file
	err := snapshotFile.Close()
	if err != nil {
		log.Printf("Vault Raft snapshot file at '%s' failed to close after interactions", snapshotFile.Name())
	}

	return err
}

// remove snapshot file
func SnapshotFileRemove(snapshotFile *os.File) error {
	// assign filename
	filename := snapshotFile.Name()

	// verify file existence
	_, err := snapshotFile.Stat()
	if err == nil {
		// remove file
		err = os.Remove(filename)
		if err == nil {
			log.Printf("removed Vault Raft snapshot at '%s'", filename)
		} else {
			log.Printf("failed to remove Vault Raft snapshot at '%s'", filename)
			log.Print("local snapshot file will need to be removed manually if desired")
			err = errors.New("snapshot not removed")
		}
	} else {
		// filenotfound
		log.Printf("Vault Raft snapshot file does not exist at expected path '%s', and therefore will not be removed automatically", filename)
		log.Print("local snapshot file will need to be removed manually if desired")
		err = errors.New("snapshot not found")
	}

	// need custom error to avoid collision with *os.PathError type from previously executed code since this func is normally deferred
	return err
}
