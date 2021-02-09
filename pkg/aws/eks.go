package aws

import (
	"crypto/md5"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

// GetExistingRoleARN checks if a role exists for AWSServiceRoleForAmazonEKS and returns its ARN
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
			return "", errors.Errorf("Failed", err)
		}
	}
	return string(*eksArn.Role.Arn), nil
}

// CreateEksClusterRole does stuff
func CreateEksClusterRole(sess *session.Session) (string, error) {
	rolesvc := iam.New(sess)
	roleInput := &iam.CreateServiceLinkedRoleInput{
		AWSServiceName: aws.String("eks.amazonaws.com"),
		Description:    aws.String("AWSServiceRoleForAmazonEKS"),
	}

	roleResult, err := rolesvc.CreateServiceLinkedRole(roleInput)
	if err != nil {
		return "", errors.Errorf("Failed to create service linked role", err)
	}

	// add policy
	policysvc := iam.New(sess)
	policyInput := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
		RoleName:  aws.String("AWSServiceRoleForAmazonEKS"),
	}

	_, err = policysvc.AttachRolePolicy(policyInput)
	if err != nil {
		return "", errors.Errorf("Failed to add policy to role:", err)
	}
	return string(*roleResult.Role.Arn), nil
}

// CreateEksCluster creates an EKS cluster with provided specs
func CreateEksCluster(sess *session.Session, name, arn string, subnets []string) error {
	digest := md5.New()
	digest.Write([]byte(name))
	hash := digest.Sum(nil)

	svc := eks.New(sess)
	input := &eks.CreateClusterInput{
		ClientRequestToken: aws.String(name + "-" + string(hash)),
		Name:               aws.String(name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds: []*string{
				aws.String(subnets[0]),
				aws.String(subnets[1]),
			},
		},
		RoleArn: aws.String(arn),
	}

	result, err := svc.CreateCluster(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case eks.ErrCodeResourceInUseException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeResourceInUseException, aerr.Error())
			case eks.ErrCodeResourceLimitExceededException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeResourceLimitExceededException, aerr.Error())
			case eks.ErrCodeInvalidParameterException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeInvalidParameterException, aerr.Error())
			case eks.ErrCodeClientException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeClientException, aerr.Error())
			case eks.ErrCodeServerException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeServerException, aerr.Error())
			case eks.ErrCodeServiceUnavailableException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeServiceUnavailableException, aerr.Error())
			case eks.ErrCodeUnsupportedAvailabilityZoneException:
				return errors.Errorf("Failed to create cluster:", eks.ErrCodeUnsupportedAvailabilityZoneException, aerr.Error())
			default:
				return errors.Errorf("Failed to create cluster:", aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return errors.Errorf("Failed to create cluser:", err)
		}
	}
	fmt.Println(result)
	return nil
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
