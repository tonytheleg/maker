package do

import (
	"fmt"

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

// DeleteDoSpace deletes a Space bucket on DigitalOcean
func DeleteDoSpace(client *s3.S3, name string) error {
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err := client.DeleteBucket(input)
	if err != nil {
		return errors.Errorf("Failed to delete Space:", err)
	}
	fmt.Println("Space", name, "deleted")
	return nil
}
