package storage

import (
	"context"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// snapshot upload to azure blob storage
func snapshotBlobUpload(container string, snapshotFile io.Reader, snapshotName string, accountURL string) error {
	// create token credential from ms entra id
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Print("unable to create microsoft entra credential")
		return err
	}

	// create service client with token credential
	client, err := azblob.NewClient(accountURL, credential, nil)
	if err != nil {
		log.Print("unable to authenticate with entra token credential")
		return err
	}

	// upload vault raft backup to azure blob storage
	_, err = client.UploadStream(context.TODO(), container, snapshotName, snapshotFile, nil)
	if err != nil {
		log.Printf("Vault backup failed to upload snapshot file %s to blob container %s", snapshotName, container)
	}

	return err
}
