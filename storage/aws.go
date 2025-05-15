package storage

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// snapshot upload to aws s3
func snapshotS3Upload(s3Bucket string, snapshotFile io.Reader, snapshotName string) error { //(*s3manager.UploadOutput, error) {
	// aws session with configuration populated automatically
	ctx := context.TODO()
	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Print("failed to load AWS shared configuration and/or credentials")
		return err
	}

	// initialize a s3 uploader from the s3 client from the shared configuration
	s3Client := s3.NewFromConfig(config)
	s3Uploader := manager.NewUploader(s3Client)

	// upload the snapshot file to the s3 bucket at specified key
	uploadResult, err := s3Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(snapshotName),
		Body:   snapshotFile,
	})
	if err != nil {
		log.Printf("Vault backup failed to upload snapshot file %s to S3 bucket %s", snapshotName, s3Bucket)
		return err
	}

	// output s3 uploader location info
	log.Printf("Vault Raft snapshot file %s uploaded to S3 bucket %s with key %s", snapshotName, uploadResult.Location, *uploadResult.Key)
	return err
}
