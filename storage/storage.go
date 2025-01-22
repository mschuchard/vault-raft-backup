package storage

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/mschuchard/vault-raft-backup/util"
)

// unified function for interfacing with all snapshot storage transfers
func StorageTransfer(config *util.CloudConfig, snapshotPath string, cleanup bool) error {
	// use supplied prefix and snapshot base filename for full name
	snapshotName := config.Prefix + filepath.Base(snapshotPath)

	// open snapshot file
	snapshotFile, err := os.Open(snapshotPath)
	if err != nil {
		log.Printf("failed to open snapshot file %q: %v", snapshotPath, err)
		return err
	}

	// defer snapshot close and remove
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
		if cleanup {
			err = util.SnapshotFileRemove(snapshotFile)
		}
	}()

	// upload snapshot to various storage backends
	switch config.Platform {
	case util.AWS:
		err = snapshotS3Upload(config.Container, snapshotFile, snapshotName)
	case util.GCP:
		err = snapshotCSUpload(config.Container, snapshotFile, snapshotName)
	case util.LOCAL:
		err = snapshotFSCopy(config.Container, snapshotFile, snapshotName)
	default:
		log.Printf("an invalid cloud platform was specified: %s", config.Platform)
		err = errors.New("invalid cloud platform")
	}

	return err
}
