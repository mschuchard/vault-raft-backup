package storage

import (
	"context"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// snapshot upload to azure blob storage
func snapshotBlobUpload(container string, snapshotFile io.Reader, snapshotName string) error {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {

	}

	client, err := azblob.NewClient("", credential, nil)
	if err != nil {

	}

	_, err = client.UploadStream(context.TODO(), container, snapshotName, snapshotFile, nil)
	if err != nil {
		log.Printf("Vault backup failed to upload snapshot file %s to blob container %s", snapshotName, container)
	}

	return nil
}
