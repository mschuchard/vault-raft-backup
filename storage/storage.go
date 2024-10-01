package storage

import "github.com/aws/aws-sdk-go/service/s3/s3manager"

func StorageTransfer(config *awsConfig, snapshotPath string, cleanup bool) (*s3manager.UploadOutput, error) {
	return snapshotS3Upload(config, snapshotPath, cleanup)
}
