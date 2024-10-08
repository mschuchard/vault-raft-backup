package storage

import (
	"errors"
	"log"

	"github.com/mschuchard/vault-raft-backup/util"
)

// config defines parameters configured for storage backend
type config struct {
	object string
	prefix string
}

func NewConfig(backupConfig *util.AWSConfig) (*config, error) {
	// validate s3 bucket name input
	if len(backupConfig.S3Bucket) == 0 {
		log.Print("the name of a destination storage object is required as an input parameter value")
		return nil, errors.New("empty storage object input setting")
	}

	// return constructor
	return &config{
		object: backupConfig.S3Bucket,
		prefix: backupConfig.S3Prefix,
	}, nil
}
