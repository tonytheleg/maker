package do

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// CreateDoSpacesClient creates a client to interact with DO Spaces API
func CreateDoSpacesClient(spacesKey, spacesSecret, defaultRegion string) *s3.S3 {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:    aws.String("https://" + defaultRegion + ".digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	return s3Client
}

// CreateDoSpace creats a Spaces bucket on DigitalOcean
func CreateDoSpace(client *s3.S3, name string) error {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := client.CreateBucket(params)
	if err != nil {
		return errors.Errorf("Failed to create Space:", err)
	}
	fmt.Println("Space", name, "created")
	return nil
}

// GetDoSpaceInfo may or may not get info about a space...
func GetDoSpaceInfo(client *s3.S3, name string) error {
	spaces, err := client.ListBuckets(nil)
	if err != nil {
		return errors.Errorf("Failed to list spaces:", err)
	}

	input := &s3.ListObjectsInput{Bucket: aws.String(name)}

	objects, err := client.ListObjects(input)
	if err != nil {
		return errors.Errorf("Failed to fetch objects in space:", err)
	}

	for _, bucket := range spaces.Buckets {
		if aws.StringValue(bucket.Name) == name {
			fmt.Printf("Name: %s\nCreation Date: %s\n\n", aws.StringValue(bucket.Name), bucket.CreationDate.Format(time.UnixDate))
		}
	}
	fmt.Println("Space contents:")
	for _, obj := range objects.Contents {
		fmt.Printf(" - %s\n", aws.StringValue(obj.Key))
	}
	return nil
}

// DeleteSpaceObjects removes all objects in a space to prep for deletion
func DeleteSpaceObjects(client *s3.S3, name string) {
	// confirm that deleteing space will delete all files first
	var confirmation string
	fmt.Printf("\nWARNING: To delete a Space, all objects in that space must be deleted!\n")
	fmt.Print("Do you wish to continue? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		err := Configure()
		if err != nil {
			return errors.Wrap(err, "Failed to configure")
		}
	}
	return nil

	// loop through all objects in bucket and delete first
	listInput := &s3.ListObjectsInput{Bucket: aws.String(name)}
	objects, err := client.ListObjects(listInput)
	if err != nil {
		return errors.Errorf("Failed to fetch objects in space:", err)
	}

	for _, obj := range objects.Contents {
		fmt.Printf(" - %s\n", aws.StringValue(obj.Key))

		input := &s3.DeleteObjectInput{
			Bucket: aws.String(name),
			Key:    aws.String(aws.StringValue(obj.Key)),
		}

		_, err := client.DeleteObject(input)
		if err != nil {
			return errors.Errorf("Failed to remove objects in space:", err)
		}

	}
}

// DeleteDoSpace deletes a Space bucket on DigitalOcean
func DeleteDoSpace(client *s3.S3, name string) error {
	// confirm that deleteing space will delete all files first

	// loop through all objects in bucket and delete first
	listInput := &s3.ListObjectsInput{Bucket: aws.String(name)}
	objects, err := client.ListObjects(listInput)
	if err != nil {
		return errors.Errorf("Failed to fetch objects in space:", err)
	}

	for _, obj := range objects.Contents {
		fmt.Printf(" - %s\n", aws.StringValue(obj.Key))

		input := &s3.DeleteObjectInput{
			Bucket: aws.String(name),
			Key:    aws.String(aws.StringValue(obj.Key)),
		}

		_, err := client.DeleteObject(input)
		if err != nil {
			return errors.Errorf("Failed to remove objects in space:", err)
		}

	}
	// delete bucket
	deleteInput := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = client.DeleteBucket(deleteInput)
	if err != nil {
		return errors.Errorf("Failed to delete Space:", err)
	}
	fmt.Println("Space", name, "deleted")
	return nil

}
