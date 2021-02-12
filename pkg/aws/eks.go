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
		RoleName: aws.String("EKSClusterRole"),
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
	//policy := `{ "Version": "2012-10-17", "Statement": [{ "Effect": "Allow", "Principal": { "AWS": "arn:aws:iam::898425707596:root" }, "Action": "sts:AssumeRole" }]}`
	policy := `{ "Version": "2012-10-17", "Statement": [{ "Effect": "Allow", "Principal": { "AWS": "arn:aws:iam::898425707596:root" }, "Action": "sts:AssumeRole" }, { "Effect": "Allow", "Principal": { "Service": "ec2.amazonaws.com" }, "Action": "sts:AssumeRole" }]}`
	rolesvc := iam.New(sess)
	roleInput := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(policy),
		Description:              aws.String("EKS Role with Cluster and Node policies"),
		RoleName:                 aws.String("EKSClusterRole"),
	}

	roleResult, err := rolesvc.CreateRole(roleInput)
	if err != nil {
		return "", errors.Errorf("Failed to create role", err)
	}

	// add policy
	policysvc := iam.New(sess)
	policies := []*iam.AttachRolePolicyInput{
		{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
			RoleName:  aws.String("EKSClusterRole"),
		},
		{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"),
			RoleName:  aws.String("EKSClusteEKSClusterRolerRoleTwo"),
		},
		{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
			RoleName:  aws.String("EKSClusterRole"),
		},
	}

	for _, policyInput := range policies {
		_, err = policysvc.AttachRolePolicy(policyInput)
		if err != nil {
			return "", errors.Errorf("Failed to add policy %s to role:", *policyInput.PolicyArn, err)
		}
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

// CreateEksNodeGroup creates workers for the EKS cluster just created
func CreateEksNodeGroup(sess *session.Session, name, arn string, nodeCount int, subnets []string) error {
	digest := md5.New()
	digest.Write([]byte(name))
	hash := digest.Sum(nil)

	svc := eks.New(sess)
	input := &eks.CreateNodegroupInput{
		CapacityType:       aws.String("On-Demand"),
		ClientRequestToken: aws.String(name + "-nodegroup-" + string(hash)),
		ClusterName:        aws.String(name),
		InstanceTypes:      aws.StringSlice([]string{"t3.micro"}),
		NodeRole:           aws.String("EKSClusterRole"),
		NodegroupName:      aws.String(name + "-nodegroups"),
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(int64(nodeCount)),
			MaxSize:     aws.Int64(int64(nodeCount)),
			MinSize:     aws.Int64(int64(nodeCount)),
		},
		Subnets: aws.StringSlice(subnets),
	}

	result, err := svc.CreateNodegroup(input)
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

// GetEksKubeconfig creates a kubeconfig needed to access the cluster
func GetEksKubeconfig() {
	/* Kubeconfig
	   https://pkg.go.dev/k8s.io/kops@v1.19.0/pkg/kubeconfig#KubeconfigBuilder
	   https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html
	   Replace the <endpoint-url> with the endpoint URL that was created for your cluster.

	   Replace the <base64-encoded-ca-cert> with the certificateAuthority.data that was created for your cluster.

	   Replace the <cluster-name> with your cluster name.
	   *Get this info from describe cluster output

	   *apiVersion: v1
	   clusters:
	   - cluster:
	       server: <endpoint-url>
	       certificate-authority-data: <base64-encoded-ca-cert>
	     name: kubernetes
	   contexts:
	   - context:
	       cluster: kubernetes
	       user: aws
	     name: aws
	   *current-context: aws
	   *kind: Config
	   preferences: {}
	   users:
	   - name: aws
	     user:
	       exec:
	         apiVersion: client.authentication.k8s.io/v1alpha1
	         command: aws
	         args:
	           - "eks"
	           - "get-token"
	           - "--cluster-name"
	           - "<cluster-name>"
	           # - "--role"
	           # - "<role-arn>"
	         # env:
	           # - name: AWS_PROFILE
	           #   value: "<aws-profile>"
	   }

	   type KubectlConfig struct {
	   	Kind           string                    `json:"kind"`
	   	ApiVersion     string                    `json:"apiVersion"`
	   	CurrentContext string                    `json:"current-context"`
	   	Clusters       []*KubectlClusterWithName `json:"clusters"`
	   	Contexts       []*KubectlContextWithName `json:"contexts"`
	   	Users          []*KubectlUserWithName    `json:"users"`
	   }

	   type KubectlClusterWithName struct {
	   	Name    string         `json:"name"`
	   	Cluster KubectlCluster `json:"cluster"`
	   }

	   type KubectlCluster struct {
	   	Server                   string `json:"server,omitempty"`
	   	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
	   }

	   type KubectlContextWithName struct {
	   	Name    string         `json:"name"`
	   	Context KubectlContext `json:"context"`
	   }

	   type KubectlContext struct {
	   	Cluster string `json:"cluster"`
	   	User    string `json:"user"`
	   }

	   type KubectlUserWithName struct {
	   	Name string      `json:"name"`
	   	User KubectlUser `json:"user"`
	   }

	   type KubectlUser struct {
	   	ClientCertificateData []byte `json:"client-certificate-data,omitempty"`
	   	ClientKeyData         []byte `json:"client-key-data,omitempty"`
	   	Password              string `json:"password,omitempty"`
	   	Username              string `json:"username,omitempty"`
	   	Token                 string `json:"token,omitempty"`
	   }

	*/
}

// DeleteEksCluster destroys an EKS cluster
func DeleteEksCluster() {

}

// DeleteEksClusterRole removes the service role used for EKS
func DeleteEksClusterRole() {

}
