package aws

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/mitodl/vault-raft-backup/util"
)

// AWSConfig is for aws interaction
type AWSConfig struct {
	s3Bucket string
	s3Prefix string
	s3Region string
}

// aws config constructor
func NewAWSConfig() *AWSConfig {
	return &AWSConfig{
		s3Bucket: os.Getenv("S3_BUCKET"),
		s3Prefix: os.Getenv("S3_PREFIX"),
		s3Region: os.Getenv("AWS_REGION"),
	}
}

// snapshot upload to s3
func SnapshotS3Upload(config *AWSConfig, snapshotPath string) (*s3manager.UploadOutput, error) {
	// open snapshot and defer closing
	snapshotFile, err := os.Open(snapshotPath)
	if err != nil {
		fmt.Printf("Failed to open snapshot file %q: %v", snapshotPath, err)
		return nil, err
	}
	defer util.SnapshotFileClose(snapshotFile)

	// aws session
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.s3Region),
	}))

	// initialize an uploader with the session and default options
	uploader := s3manager.NewUploader(awsSession)

	// determine vault backup base for s3 key
	snapshotPathBase := filepath.Base(snapshotPath)

	// upload the snapshot to the s3bucket at specified key
	uploadResult, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.s3Bucket),
		Key:    aws.String(config.s3Prefix + "-" + snapshotPathBase),
		Body:   snapshotFile,
	})
	if err != nil {
		fmt.Println("Vault backup failed to upload to S3 bucket " + config.s3Bucket)
		fmt.Println(err)
		return nil, err
	}

	// output info
	log.Printf("Vault Raft snapshot uploaded to, %s\n", aws.StringValue(&uploadResult.Location))

	return uploadResult, nil
}
