package storage

import (
	"errors"
	"log"

	"github.com/mschuchard/vault-raft-backup/util"
)

// storageConfig defines parameters for storage backend
type storageConfig struct {
	object string
	prefix string
}

func NewConfig(backupConfig *util.AWSConfig) (*storageConfig, error) {
	// validate s3 bucket name input
	if len(backupConfig.S3Bucket) == 0 {
		log.Print("the name of an AWS S3 bucket is required as an input parameter value")
		return nil, errors.New("empty s3 bucket input setting")
	}

	// return constructor
	return &storageConfig{
		object: backupConfig.S3Bucket,
		prefix: backupConfig.S3Prefix,
	}, nil
}
