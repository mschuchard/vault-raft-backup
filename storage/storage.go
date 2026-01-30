package storage

import (
	"compress/gzip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mschuchard/vault-raft-backup/enum"
	"github.com/mschuchard/vault-raft-backup/util"
)

// unified function for interfacing with all snapshot storage transfers
func StorageTransfer(config *util.CloudConfig, snapshot *util.SnapshotConfig) error {
	// use supplied prefix and snapshot base filename for full name
	snapshotName := config.Prefix + filepath.Base(snapshot.Path)

	// open snapshot file
	snapshotFile, err := os.Open(snapshot.Path)
	if err != nil {
		log.Printf("failed to open snapshot file %s", snapshot.Path)
		return err
	}

	// defer snapshot close and remove
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
		if snapshot.Cleanup && err == nil {
			err = util.SnapshotFileRemove(snapshotFile)
		}
	}()

	// wrap in streaming compression if enabled
	var reader io.ReadCloser = snapshotFile
	if snapshot.CompressionLevel > 0 {
		// create compressed reader at specified level
		reader, err = CompressReader(snapshotFile, snapshot.CompressionLevel)
		if err != nil {
			log.Printf("failed to compress snapshot file %s", snapshot.Path)
			return err
		}

		// defer reader close
		defer reader.Close()
	} else {
		log.Print("snapshot will be transferred without compression")
	}

	// upload snapshot to various storage backends
	switch config.Platform {
	case enum.AWS:
		err = snapshotS3Upload(config.Container, reader, snapshotName)
	case enum.AZ:
		err = snapshotBlobUpload(config.Container, reader, snapshotName, config.AZAccountURL)
	case enum.GCP:
		err = snapshotCSUpload(config.Container, reader, snapshotName)
	case enum.LOCAL:
		err = snapshotFSCopy(config.Container, reader, snapshotName)
	default:
		log.Printf("an invalid storage platform was specified: %s", config.Platform)
		err = errors.New("invalid storage platform")
	}
	if err != nil {
		log.Print("snapshot storage transfer failed")
		return err
	}

	// potentially return deferred error
	return err
}

// compression helper function that wraps a reader with gzip compression, and returns a reader for the compressed data
func CompressReader(reader io.Reader, level int) (io.ReadCloser, error) {
	log.Print("snapshot will be transferred with compression")

	// convert user input level to gzip level
	switch level {
	case 1:
		level = gzip.BestSpeed
	case 2:
		level = gzip.DefaultCompression
	case 3:
		level = gzip.BestCompression
	default:
		log.Printf("invalid user input compression level %d", level)
		return nil, errors.New("invalid compression level")
	}
	log.Printf("using converted gzip compression level %d", level)

	// initialize pipe
	pipeRead, pipeWrite := io.Pipe()
	// create gzip writer
	gzWriter, err := gzip.NewWriterLevel(pipeWrite, level)
	if err != nil {
		log.Printf("invalid compression level %d", level)
		return nil, err
	}

	// start compression in background
	go func() {
		defer pipeWrite.Close()
		defer gzWriter.Close()

		// compress data from reader into pipe
		if _, err := io.Copy(gzWriter, reader); err != nil {
			pipeWrite.CloseWithError(err)
			return
		}

		// flush gzip writer
		if err := gzWriter.Close(); err != nil {
			pipeWrite.CloseWithError(err)
		}
	}()

	// return the read end of the pipe
	return pipeRead, nil
}
