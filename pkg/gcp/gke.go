package gcp

import (
	"context"
	"fmt"

	container "cloud.google.com/go/container/apiv1"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

// CreateGkeClient returns a client needed to interact with GKE
func CreateGkeClient(keyfile string) (*container.ClusterManagerClient, error) {
	ctx := context.Background()
	client, err := container.NewClusterManagerClient(
		ctx, option.WithCredentialsFile(keyfile))

	if err != nil {
		return nil, errors.Errorf("Failed to create cluster manager client:", err)
	}
	return client, nil
}

// CreateGkeCluster creates an EKS cluster with provided specs
func CreateGkeCluster(client *container.ClusterManagerClient, name, project, zone, nodeSize string, nodeCount int) error {
	ctx := context.Background()
	parent := "projects/" + project + "/locations/" + zone
	req := &containerpb.CreateClusterRequest{
		Cluster: &containerpb.Cluster{
			Name:        name,
			Description: "cluster created by Maker",
			NodePools: []*containerpb.NodePool{{
				Name: name + "-nodepool",
				Config: &containerpb.NodeConfig{
					MachineType: nodeSize,
				},
				InitialNodeCount: int32(nodeCount),
			}},
		},
		Parent: parent,
		// requires a cluster definition type Cluster
		// reguires a CreateClusterRequest
		// Think I also need a NodePool struct in here
		// TODO: Fill request struct fields.
	}
	_, err := client.CreateCluster(ctx, req)
	if err != nil {
		return errors.Errorf("Failed to create cluster:", err)
	}
	fmt.Println("GKE Cluster", name, "creating")
	return nil
}

// CreateGkeNodeGroup creates workers for the EKS cluster just created
func CreateGkeNodeGroup() {
}

// CreateKubeconfig creates a kubeconfig needed to access the cluster
/*func CreateKubeconfig(endpoint, caData, name string) error {
	configString := `apiVersion: v1
clusters:
- cluster:
    server: ` + endpoint + `
    certificate-authority-data: ` + caData + `
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
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
        - "` + name + `"
`
	linesToWrite := configString
	err := ioutil.WriteFile(utils.ConfigFolderPath+"/aws_kubeconfig", []byte(linesToWrite), 0755)
	if err != nil {
		return errors.Errorf("Failed to write kubeconfig:", err)
	}
	return nil
}
*/

// GetCluster describes the cluster and returns cluster details needed for kubeconfig
func GetCluster() {
}

// GetClusterStatus checks the state of the EKS Cluster before creating a node group
func GetClusterStatus() {
}

// GetNodeGroupStatus grabs the current state of the node group
func GetNodeGroupStatus() {
}

// PrintGkeClusterStatus outputs EKS cluster info
func PrintGkeClusterStatus() {
}

// DeleteGkeNodeGroup deletes the node group before deleting the cluster
func DeleteGkeNodeGroup() {
}

// DeleteGkeCluster destroys an EKS cluster
func DeleteGkeCluster() {
}
