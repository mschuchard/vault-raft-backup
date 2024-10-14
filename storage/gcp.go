package storage

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

// snapshot upload to gcp cloud storage
func snapshotCSUpload(csBucket string, snapshotFile io.Reader, snapshotName string) error {
	// gcp client
	ctx := context.Background()
	gcpClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Print("unable to initialize GCP client for transfer to cloud storage bucket")
		return err
	}
	defer gcpClient.Close()

	// initialize target in cloud storage bucket
	uploadTarget := gcpClient.Bucket(csBucket).Object(snapshotName)
	// only upload if does not exist (TODO feature all below)
	// uploadTarget = uploadTarget.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your csBucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	// 	return fmt.Errorf("object.Attrs: %w", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// write snapshotfile to upload target in cloud storage bucket
	uploadTransfer := uploadTarget.NewWriter(ctx)
	if _, err = io.Copy(uploadTransfer, snapshotFile); err != nil {
		log.Printf("failed to upload snapshot file %s to bucket %s", uploadTarget.ObjectName(), uploadTarget.BucketName())
		return err
	}
	if err := uploadTransfer.Close(); err != nil {
		log.Printf("failed to close snapshot file %s upload transfer to %s", uploadTarget.ObjectName(), uploadTarget.BucketName())
		return err
	}

	return nil
}
