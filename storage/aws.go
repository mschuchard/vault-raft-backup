package storage

import (
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// snapshot upload to aws s3
func snapshotS3Upload(s3Bucket string, snapshotFile io.Reader, snapshotName string) error { //(*s3manager.UploadOutput, error) {
	// aws session with configuration populated automatically
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// initialize an s3 uploader with the session and default options
	uploader := s3manager.NewUploader(awsSession)

	// upload the snapshot file to the s3 bucket at specified key
	uploadResult, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(snapshotName),
		Body:   snapshotFile,
	})
	if err != nil {
		log.Printf("Vault backup failed to upload to S3 bucket %s", s3Bucket)
		return err
	}

	// output s3 uploader location info
	log.Printf("Vault Raft snapshot uploaded to %s", aws.StringValue(&uploadResult.Location))

	return err
	// return uploadResult, err
}
