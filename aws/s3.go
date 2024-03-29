package aws

import (
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/mitodl/vault-raft-backup/util"
)

// AWSConfig defines aws client api interaction
type AWSConfig struct {
	s3Bucket string
	s3Prefix string
	s3Region string
}

// aws config constructor
func NewAWSConfig() (*AWSConfig, error) {
	return &AWSConfig{
		s3Bucket: os.Getenv("S3_BUCKET"),
		s3Prefix: os.Getenv("S3_PREFIX"),
	}, nil
}

// snapshot upload to aws s3
func SnapshotS3Upload(config *AWSConfig, snapshotPath string) (*s3manager.UploadOutput, error) {
	// open snapshot file
	snapshotFile, err := os.Open(snapshotPath)
	if err != nil {
		log.Printf("failed to open snapshot file %q: %v", snapshotPath, err)
		return nil, err
	}
	// defer snapshot close
	defer func() {
		err = util.SnapshotFileClose(snapshotFile)
		err = util.SnapshotFileRemove(snapshotFile)
	}()

	// aws session with configuration populated automatically
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// initialize an s3 uploader with the session and default options
	uploader := s3manager.NewUploader(awsSession)
	// determine vault backup base filename for s3 key
	snapshotPathBase := filepath.Base(snapshotPath)

	// upload the snapshot file to the s3 bucket at specified key
	uploadResult, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.s3Bucket),
		Key:    aws.String(config.s3Prefix + "-" + snapshotPathBase),
		Body:   snapshotFile,
	})
	if err != nil {
		log.Printf("Vault backup failed to upload to S3 bucket %s", config.s3Bucket)
		return nil, err
	}

	// output s3 uploader location info
	log.Printf("Vault Raft snapshot uploaded to %s", aws.StringValue(&uploadResult.Location))

	return uploadResult, err
}
