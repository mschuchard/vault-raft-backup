package aws

import (
	"errors"
	"log"
	"os"
)

// AWSConfig defines aws client api interaction
type awsConfig struct {
	s3Bucket string
	s3Prefix string
}

// aws config constructor
func NewAWSConfig() (*awsConfig, error) {
	// validate s3 bucket name input
	s3Bucket := os.Getenv("S3_BUCKET")
	if len(s3Bucket) == 0 {
		log.Print("the name of an AWS S3 bucket is required as an input parameter value")
		return nil, errors.New("empty s3 bucket input setting")
	}

	// return constructor
	return &awsConfig{
		s3Bucket: s3Bucket,
		s3Prefix: os.Getenv("S3_PREFIX"),
	}, nil
}
