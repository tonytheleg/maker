package gcp

import (
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// CreateStorageClient creates a new client to interact with GCP
func CreateStorageClient(keyfile string) (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(
		ctx, option.WithCredentialsFile(keyfile))
	if err != nil {
		return nil, errors.Errorf("Failed to create client: ", err)
	}
	return client, nil
}

// CreateStorageBucket creates a storage bucket on GCP
func CreateStorageBucket(client *storage.Client, name, project string) error {
	ctx := context.Background()
	bkt := client.Bucket(name)

	err := bkt.Create(ctx, project, nil)
	if err != nil {
		return errors.Errorf("Failed to create bucket: ", err)
	}
	fmt.Println("Bucket", name, "created")
	return nil
}

// GetStorageBucketInfo outputs instance info
func GetStorageBucketInfo(client *storage.Client, name string) error {
	ctx := context.Background()
	bkt := client.Bucket(name)
	attrs, err := bkt.Attrs(ctx)
	if err != nil {
		return errors.Errorf("Failed to fetch bucket: ", err)
	}

	// get objects
	query := &storage.Query{Prefix: ""}
	var names []string
	obj := bkt.Objects(ctx, query)
	for {
		attrs, err := obj.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}

	fmt.Printf(
		"Name: %s\nCreated: %s\nLocation: %s\nStorage Class: %s\n\n",
		attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)
	fmt.Println("Bucket contents:")
	for _, file := range names {
		fmt.Printf(" - %s\n", file)
	}
	return nil
}

// DeleteStorageObjects empties a bucket for deletion
func DeleteStorageBucket(client *storage.Client, name, project string) error {
	ctx := context.Background()
	err := client.Bucket(name).Delete(ctx)
	if err != nil {
		return errors.Errorf("Failed to delete bucket:", err)
	}
	fmt.Println("Bucket", name, "deleted")
	return nil
}

// DeleteStorageBucket delets a Storage bucket from GCP
func DeleteStorageObjects(client *storage.Client, name, project string) error {
	// confirm that deleteing space will delete all files first
	var confirmation string
	fmt.Printf("\nWARNING: To delete a Storage bucket, all objects in that bucket must be deleted!\n")
	fmt.Print("Do you wish to continue? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		return errors.Errorf("Cannot proceed -- must delete files before deleting bucket")
	}
	ctx := context.Background()
	bucket := client.Bucket(name)
	item := bucket.Objects(ctx, nil)
	for {
		objAttrs, err := item.Next()
		if err != nil && err != iterator.Done {
			return errors.Errorf("Failed to fetch files from bucket")
		}
		if err == iterator.Done {
			break
		}
		if err := bucket.Object(objAttrs.Name).Delete(ctx); err != nil {
			return errors.Errorf("Failed to delete files from bucket")
		}
	}
	fmt.Println("Deleted all object items in the bucket specified.")
	return nil
}
