package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/mschuchard/vault-raft-backup/util"
)

func StorageTransfer(config *config, snapshotPath string, cleanup bool) (*s3manager.UploadOutput, error) {
	// use supplied prefix and snapshot base filename for full name
	snapshotName := config.prefix + "-" + filepath.Base(snapshotPath)

	// open snapshot file
	snapshotFile, err := os.Open(snapshotPath)
	if err != nil {
		log.Printf("failed to open snapshot file %q: %v", snapshotPath, err)
		return nil, err
	}

	// defer snapshot close and remove
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
		if cleanup {
			err = util.SnapshotFileRemove(snapshotFile)
		}
	}()

	// TODO: clobbers deferred err from snapshot close and remove
	return snapshotS3Upload(config, snapshotFile, snapshotName)
}
