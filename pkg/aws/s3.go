package aws

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

// CreateS3Sessions creates a session to connect to AWS
func CreateS3Client(credentialsFile, defaultRegion string) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(defaultRegion),
		Credentials: credentials.NewSharedCredentials(credentialsFile, "default")},
	)
	if err != nil {
		return nil, errors.Errorf("Failed to create session:", err)
	}
	s3Client := s3.New(sess)
	return s3Client, nil
}

// CreateS3Bucket creats an S3 bucket on AWS
func CreateS3Bucket(client *s3.S3, name string) error {
	params := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := client.CreateBucket(params)
	if err != nil {
		return errors.Errorf("Failed to create Bucket:", err)
	}
	fmt.Println("Bucket", name, "created")
	return nil
}

// GetS3BucketInfo may or may not get info about a space...
func GetS3BucketInfo(client *s3.S3, name string) error {
	spaces, err := client.ListBuckets(nil)
	if err != nil {
		return errors.Errorf("Failed to list buckets:", err)
	}

	input := &s3.ListObjectsInput{Bucket: aws.String(name)}

	objects, err := client.ListObjects(input)
	if err != nil {
		return errors.Errorf("Failed to fetch objects in bucket:", err)
	}

	for _, bucket := range spaces.Buckets {
		if aws.StringValue(bucket.Name) == name {
			fmt.Printf("Name: %s\nCreation Date: %s\n\n", aws.StringValue(bucket.Name), bucket.CreationDate.Format(time.UnixDate))
		}
	}
	fmt.Println("Bucket contents:")
	for _, obj := range objects.Contents {
		fmt.Printf(" - %s\n", aws.StringValue(obj.Key))
	}
	return nil
}

// DeleteS3Objects removes all objects in a bucket to prep for deletion
func DeleteS3Objects(client *s3.S3, name string) error {
	// confirm that deleteing space will delete all files first
	var confirmation string
	fmt.Printf("\nWARNING: To delete an S3 bucket, all objects in that bucket must be deleted!\n")
	fmt.Print("Do you wish to continue? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		return errors.Errorf("Cannot proceed -- must delete files before deleting bucket")
	} else {
		// loop through all objects in bucket and delete first
		listInput := &s3.ListObjectsInput{Bucket: aws.String(name)}
		objects, err := client.ListObjects(listInput)
		if err != nil {
			return errors.Errorf("Failed to fetch objects in bucket:", err)
		}

		for _, obj := range objects.Contents {
			input := &s3.DeleteObjectInput{
				Bucket: aws.String(name),
				Key:    aws.String(aws.StringValue(obj.Key)),
			}

			_, err := client.DeleteObject(input)
			if err != nil {
				return errors.Errorf("Failed to remove objects in bucket:", err)
			}
		}
		fmt.Println("All objects from", name, "deleted")
	}
	return nil
}

// DeleteS3Bucket deletes an S3 bucket on AWS
func DeleteS3Bucket(client *s3.S3, name string) error {
	deleteInput := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err := client.DeleteBucket(deleteInput)
	if err != nil {
		return errors.Errorf("Failed to delete bucket:", err)
	}
	fmt.Println("Bucket", name, "deleted")
	return nil
}
