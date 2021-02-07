package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

// Required
//Name *string
//ResourcesVpcConfig *types.VpcConfigRequest
//RoleArn *string <--- https://docs.aws.amazon.com/eks/latest/userguide/service_IAM_role.html#create-service-role

// GetExistingRoleARN checks if a role exists for AWSServiceRoleForAmazonEKS
func GetExistingRoleARN(sess *session.Session) (string, error) {
	svc := iam.New(sess)
	roleInput := &iam.GetRoleInput{
		RoleName: aws.String("AWSServiceRoleForAmazonEKS"),
	}

	eksArn, err := svc.GetRole(roleInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				return "", errors.Errorf("Failed", iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				return "", errors.Errorf("Failed", iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				return "", errors.Errorf("Failed", err)
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return "", errors.Errorf("Failed", err)
		}
	}
	return string(*eksArn.Role.Arn), nil
}

// CreateEksClusterRole does stuff
func CreateEksClusterRole(sess *session.Session) (string, error) {
	svc := iam.New(sess)
	input := &iam.CreateServiceLinkedRoleInput{
		AWSServiceName: aws.String("eks.amazonaws.com"),
		Description:    aws.String("AWSServiceRoleForAmazonEKS"),
	}

	result, err := svc.CreateServiceLinkedRole(input)
	if err != nil {
		return "", errors.Errorf("Failed to create service linked role", err)
	}
	return string(*result.Role.Arn), nil
}

// CreateEksCluster creates an EKS cluster with provided specs
func CreateEksCluster() {
}

// GetClusterID fetches the EKS cluster ID for status or deleting
func GetClusterID() {

}

// PrintEksClusterStatus outputs EKS cluster info
func PrintEksClusterStatus() {

}

// DeleteEksCluster destroys an EKS cluster
func DeleteEksCluster() {

}

// DeleteEksClusterRole removes the service role used for EKS
func DeleteEksClusterRole() {

}
