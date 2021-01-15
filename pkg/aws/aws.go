package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

// CreateAwsSession sets up a new session using the config file
func CreateAwsSession(defaultRegion, credentialsFile string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(defaultRegion),
		Credentials: credentials.NewSharedCredentials(credentialsFile, "default")},
	)
	if err != nil {
		return nil, errors.Errorf("Failed to create session", err)
	}
	return sess, nil
}

// CreateEc2Instance creates an ec2 instance with provided specs
func CreateEc2Instance(
	sess *session.Session, name string, region string, ami string, instanceType string) error {

	// Create the instance
	svc := ec2.New(sess)
	result, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String(ami),
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		return errors.Wrapf(err, "Failed to create EC2 instance %s", name)
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
		return errors.Wrapf(err, "Failed to tag instance %s with name %s",
			*result.Instances[0].InstanceId, name)
	}
	fmt.Println("Successfully tagged instance")
	return nil
}

// GetInstanceID fetches the EC2 Instance ID for status or deleting
func GetInstanceID(sess *session.Session, name string) (string, error) {
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
				return "", errors.Wrapf(
					errors.New(aerr.Error()), "Failed to describe instance %s", name)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return "", errors.Wrapf(
				errors.New(err.Error()), "Failed to describe instance %s", name)

		}
		return ""
	}
	return *result.Reservations[0].Instances[0].InstanceId
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

	fmt.Printf(
		"Name: %s\nID: %s\n\nAMI: %s\nInstance Type: %s\n\nPublic IP: %s\nPublic DNS: %s\nRegion: %s\nStatus: %s\n",
		*result.Reservations[0].Instances[0].Tags[0].Value,
		*result.Reservations[0].Instances[0].InstanceId,
		*result.Reservations[0].Instances[0].ImageId,
		*result.Reservations[0].Instances[0].InstanceType,
		*result.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicIp,
		*result.Reservations[0].Instances[0].NetworkInterfaces[0].Association.PublicDnsName,
		*result.Reservations[0].Instances[0].Placement.AvailabilityZone,
		*result.Reservations[0].Instances[0].State.Name,
	)
}

// DeleteEc2Instance destroys an instance
func DeleteEc2Instance(sess *session.Session, id string) {
	svc := ec2.New(sess)
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	result, err := svc.TerminateInstances(input)
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
	fmt.Printf("Success: %s is %s",
		*result.TerminatingInstances[0].InstanceId,
		*result.TerminatingInstances[0].CurrentState.Name,
	)
}
