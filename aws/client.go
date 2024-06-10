package aws

import (
	"errors"
	"log"

	"github.com/mschuchard/vault-raft-backup/util"
)

// AWSConfig defines aws client api interaction
type awsConfig struct {
	s3Bucket string
	s3Prefix string
}

// aws config constructor
func NewAWSConfig(backupConfig *util.BackupConfig) (*awsConfig, error) {
	// validate s3 bucket name input
	if len(backupConfig.AWSConfig.S3Bucket) == 0 {
		log.Print("the name of an AWS S3 bucket is required as an input parameter value")
		return nil, errors.New("empty s3 bucket input setting")
	}

	// return constructor
	return &awsConfig{
		s3Bucket: backupConfig.AWSConfig.S3Bucket,
		s3Prefix: backupConfig.AWSConfig.S3Prefix,
	}, nil
}
