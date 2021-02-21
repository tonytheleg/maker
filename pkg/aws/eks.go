package aws

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

// GetExistingRoleARN checks if a role exists for AWSServiceRoleForAmazonEKS and returns its ARN
func GetExistingRoleARN(sess *session.Session) (string, error) {
	svc := iam.New(sess)
	roleInput := &iam.GetRoleInput{
		RoleName: aws.String("EKSClusterRole"),
	}

	eksArn, err := svc.GetRole(roleInput)
	if err != nil {
		return "", errors.Errorf("Failed", err)
	}
	return string(*eksArn.Role.Arn), nil
}

// CreateEksClusterRole does stuff
func CreateEksClusterRole(sess *session.Session) (string, error) {
	policy := `{ "Version": "2012-10-17", "Statement": [{ "Effect": "Allow", "Principal": { "AWS": "arn:aws:iam::898425707596:root" }, "Action": "sts:AssumeRole" }, { "Effect": "Allow", "Principal": { "Service": "ec2.amazonaws.com", "Service": "eks.amazonaws.com" }, "Action": "sts:AssumeRole" }]}`
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
			RoleName:  aws.String("EKSClusterRole"),
		},
		{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
			RoleName:  aws.String("EKSClusterRole"),
		},
		{
			PolicyArn: aws.String("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"),
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

	_, err := svc.CreateCluster(input)
	if err != nil {
		return errors.Errorf("Failed to create cluser:", err)
	}
	fmt.Println("Creating", name, "cluster. This can take up to 10-15 minutes...")
	return nil
}

// CreateEksNodeGroup creates workers for the EKS cluster just created
func CreateEksNodeGroup(sess *session.Session, name, arn, nodeSize string, nodeCount int, subnets []string) error {
	// pre-check

	for {
		if status, err := GetClusterStatus(sess, name); status != "FAILED" {
			if err != nil {
				return errors.Errorf("Failed to get cluster status to create node group:", err)
			}
			if status == "ACTIVE" {
				fmt.Println("Cluster completed!")
				fmt.Println("Current status:", status)
				break
			}
			fmt.Println("Waiting for cluster completion to create node group -- Current status:", status)
			time.Sleep(1 * time.Minute)
		}
	}

	digest := md5.New()
	digest.Write([]byte(name))
	hash := digest.Sum(nil)

	svc := eks.New(sess)
	input := &eks.CreateNodegroupInput{
		CapacityType:       aws.String("ON_DEMAND"),
		ClientRequestToken: aws.String(name + "-nodegroup-" + string(hash)),
		ClusterName:        aws.String(name),
		InstanceTypes:      aws.StringSlice([]string{nodeSize}),
		NodeRole:           aws.String(arn),
		NodegroupName:      aws.String(name + "-nodegroup"),
		ScalingConfig: &eks.NodegroupScalingConfig{
			DesiredSize: aws.Int64(int64(nodeCount)),
			MaxSize:     aws.Int64(int64(nodeCount)),
			MinSize:     aws.Int64(int64(nodeCount)),
		},
		Subnets: aws.StringSlice(subnets),
	}

	_, err := svc.CreateNodegroup(input)
	if err != nil {
		return errors.Errorf("Failed to create cluser:", err)
	}
	fmt.Println("Node group", *input.NodegroupName, "creating")
	return nil
}

func newClientset(cluster *eks.Cluster) (*kubernetes.Clientset, error) {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: aws.StringValue(cluster.Name),
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(aws.StringValue(cluster.CertificateAuthority.Data))
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        aws.StringValue(cluster.Endpoint),
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// GetDeployStatus checks the state of the EKS Cluster before creating a node group
func GetClusterStatus(sess *session.Session, name string) (string, error) {
	svc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		return "", errors.Errorf("Failed to fetch cluster status:", err)
	}
	return *result.Cluster.Status, nil
}

// GetNodeGroupStatus grabs the current state of the node group
func GetNodeGroupStatus(sess *session.Session, clusterName, nodeGroupName string) (string, error) {
	svc := eks.New(sess)
	input := &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodeGroupName),
	}

	result, err := svc.DescribeNodegroup(input)
	if err != nil {
		return "", errors.Errorf("Failed to fetch cluster status:", err)
	}
	return *result.Nodegroup.Status, nil
}

// PrintEksClusterStatus outputs EKS cluster info
func PrintEksClusterStatus(sess *session.Session, name string) error {
	svc := eks.New(sess)
	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}

	result, err := svc.DescribeCluster(input)
	if err != nil {
		return errors.Errorf("Failed to fetch cluster status:", err)
	}

	fmt.Printf("\nName: %s\nARN: %s\n\nEndpoint: %s\nService IP: %s\n\nVersion: %s\nCreated: %s\nState: %s\n",
		*result.Cluster.Name,
		*result.Cluster.Arn,
		*result.Cluster.Endpoint,
		*result.Cluster.KubernetesNetworkConfig.ServiceIpv4Cidr,
		*result.Cluster.Version,
		*result.Cluster.CreatedAt,
		*result.Cluster.Status,
	)
	return nil
}

// DeleteEksNodeGroup deletes the node group before deleting the cluster
func DeleteEksNodeGroup(sess *session.Session, clusterName, nodeGroupName string) error {
	svc := eks.New(sess)
	input := &eks.DeleteNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodeGroupName),
	}

	_, err := svc.DeleteNodegroup(input)
	if err != nil {
		return errors.Errorf("Failed to delete the node group:", err)
	}
	fmt.Println("Node group", nodeGroupName, "deleted")
	return nil
}

// DeleteEksCluster destroys an EKS cluster
func DeleteEksCluster(sess *session.Session, name, nodeGroupName string) error {
	// pre-check
	for {
		status, err := GetNodeGroupStatus(sess, name, nodeGroupName)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				fmt.Println("err is:", err)
				fmt.Println("aerr is:", aerr)
				fmt.Println("aerrCode is:", aerr.Code())
				break
			}
			break
		}
		fmt.Println("Waiting for nodes to finish deleting -- Current status:", status)
		time.Sleep(1 * time.Minute)
	}

	svc := eks.New(sess)
	input := &eks.DeleteClusterInput{
		Name: aws.String(name),
	}

	_, err := svc.DeleteCluster(input)
	if err != nil {
		return errors.Errorf("Failed to delete the cluster:", err)
	}
	fmt.Println("Cluster", name, "deleted")
	return nil
}

func main() {
	name := "wonderful-outfit-1583362361"
	region := "us-east-2"
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	eksSvc := eks.New(sess)

	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}
	result, err := eksSvc.DescribeCluster(input)
	if err != nil {
		log.Fatalf("Error calling DescribeCluster: %v", err)
	}
	clientset, err := newClientset(result.Cluster)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error getting EKS nodes: %v", err)
	}
	log.Printf("There are %d nodes associated with cluster %s", len(nodes.Items), name)
}

// GetEksKubeconfig creates a kubeconfig needed to access the cluster
func GetEksKubeconfig(sess *session.Session, name string) {
	/* Kubeconfig
	   https://pkg.go.dev/k8s.io/kops@v1.19.0/pkg/kubeconfig#KubeconfigBuilder
	   https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html
	   Replace the <endpoint-url> with the endpoint URL that was created for your cluster.

	   Replace the <base64-encoded-ca-cert> with the certificateAuthority.data that was created for your cluster.

	   Replace the <cluster-name> with your cluster name.
	   *Get this info from describe cluster output
	*/
	/*
		cluster := KubectlClusterWithName{
			Name: name,
			Cluster: KubectlCluster{
				Server:                   "INSERT-ENDPOINT-HERE",
				CertificateAuthorityData: "INSERT-CADATA-HERE",
			},
		}

		context := KubectlContextWithName{
			Name: name,
			Context: KubectlContext{
				Cluster: name,
				User:    "maker",
			},
		}

		user := KubectlUserWithName{
			Name: "INSERT-NAME",
			User: KubectlUser{},
		}

		kubeConf := KubectlConfig{
			Kind:           "Config",
			APIVersion:     "v1",
			CurrentContext: name,
			Clusters:       cluster,
			Contexts:       context,
			Users:          "INSERT-USERS-OBJECT",
		}
	*/
	/*
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
	*/
}
