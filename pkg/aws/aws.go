package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateAwsSession sets up a new session using the config file
func CreateAwsSession(defaultRegion, credentialsFile string) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(defaultRegion),
		Credentials: credentials.NewSharedCredentials(credentialsFile, "default")},
	)
	if err != nil {
		panic(err)
	}
	return sess

}

// CreateEc2Instance creates an ec2 instance with provided specs
func CreateEc2Instance(sess *session.Session, name string, region string, ami string, instanceType string) {
	// Create EC2 service client
	svc := ec2.New(sess)

	result, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(ami),
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}

	fmt.Println("Created instance", *result.Instances[0].InstanceId)

	// Add tags to the created instance
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{result.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
	})
	if err != nil {
		fmt.Println("Could not create tags for instance", result.Instances[0].InstanceId, err)
		return
	}

	fmt.Println("Successfully tagged instance")
}

// PrintEc2Status outputs ec2 instance info
func PrintEc2Status(sess *session.Session, name string) {
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(name),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

// DeleteEc2Instance destroys an instance
func DeleteEc2Instance() {
}
