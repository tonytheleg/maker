package do

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
