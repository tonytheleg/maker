package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"maker/internal/utils"
	"time"

	container "cloud.google.com/go/container/apiv1"
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	credentialspb "google.golang.org/genproto/googleapis/iam/credentials/v1"
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

// CreateGkeCluster creates an GKE cluster with provided specs
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
	}
	_, err := client.CreateCluster(ctx, req)
	if err != nil {
		return errors.Errorf("Failed to create cluster:", err)
	}
	fmt.Println("GKE Cluster", name, "creating")
	return nil
}

// CreateKubeconfig creates a kubeconfig needed to access the cluster
func CreateKubeconfig(client *container.ClusterManagerClient, name, project, zone, accessToken string) error {
	// Cluster has to be finished before an Endpoint IP is available
	for {
		if cluster, err := GetCluster(client, name, project, zone); cluster.Status.String() != "ERROR" {
			if err != nil {
				return errors.Errorf("Failed to get cluster status:", err)
			}
			if cluster.Status.String() == "RUNNING" {
				if cluster.Endpoint == "" {
					fmt.Println("Cluster completed!...waiting for Endpoint to populate")
					fmt.Println("Current status:", cluster.Status.String())
				}
				fmt.Println("Cluster completed!")
				fmt.Println("Current status:", cluster.Status.String())
				break
			}
			fmt.Println("Waiting for cluster completion -- Current status:", cluster.Status.String())
			time.Sleep(1 * time.Minute)
		}
	}
	cluster, err := GetCluster(client, name, project, zone)
	contextString := "gke_" + project + "_" + zone + "_" + name
	configString := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ` + cluster.MasterAuth.ClusterCaCertificate + `
    server: https://` + cluster.Endpoint + `
  name: ` + contextString + `
contexts:
- context:
    cluster: ` + contextString + `
    user: ` + contextString + `
  name: ` + contextString + `
current-context: ` + contextString + `
kind: Config
preferences: {}
users:
- name: ` + contextString + `
  user:
    auth-provider:
      config:
        access-token: ` + accessToken + `
        cmd-args: config config-helper --format=json
        cmd-path: /usr/lib64/google-cloud-sdk/bin/gcloud
        expiry: "2021-02-28T15:31:34Z"
        expiry-key: '{.credential.token_expiry}'
        token-key: '{.credential.access_token}'
      name: gcp
`
	linesToWrite := configString
	err = ioutil.WriteFile(utils.ConfigFolderPath+"/gke_kubeconfig", []byte(linesToWrite), 0755)
	if err != nil {
		return errors.Errorf("Failed to write kubeconfig:", err)
	}
	fmt.Println("Kubeconfig created at", utils.ConfigFolderPath+"/gke_kubeconfig")
	fmt.Printf("To use the kubeconfig, be sure to run 'export KUBECONFIG=%s/gke_kubeconfig'\n", utils.ConfigFolderPath)
	return nil
}

// GetCluster describes the cluster and returns cluster details needed for kubeconfig and status
func GetCluster(client *container.ClusterManagerClient, name, project, zone string) (*containerpb.Cluster, error) {
	ctx := context.Background()

	req := &containerpb.GetClusterRequest{
		Name: "projects/" + project + "/locations/" + zone + "/clusters/" + name,
	}
	resp, err := client.GetCluster(ctx, req)
	if err != nil {
		return nil, errors.Errorf("Failed to fetch cluster data:", err)
	}
	return resp, nil
}

// PrintGkeClusterStatus outputs EKS cluster info
func PrintGkeClusterStatus(client *container.ClusterManagerClient, name, project, zone string) {
	cluster, _ := GetCluster(client, name, project, zone)
	fmt.Printf(
		"Name: %s\nVersion: %s\nEndpoint: %s\nNetwork: %s\nServices CIDR: %s\nZone: %s\n\nNode Pool Name: %s\nNode Count: %d\nMachine Type: %s\nImage: %s\nStatus: %s\nCreation Date: %s\n\n",
		cluster.Name,
		cluster.CurrentNodeVersion,
		cluster.Endpoint,
		cluster.Network,
		cluster.ServicesIpv4Cidr,
		cluster.Zone,
		cluster.NodePools[0].Name,
		cluster.CurrentNodeCount,
		cluster.NodeConfig.MachineType,
		cluster.NodeConfig.ImageType,
		cluster.Status,
		cluster.CreateTime,
	)
}

// DeleteGkeCluster destroys an EKS cluster
func DeleteGkeCluster(client *container.ClusterManagerClient, name, project, zone string) error {
	ctx := context.Background()

	req := &containerpb.DeleteClusterRequest{
		Name: "projects/" + project + "/locations/" + zone + "/clusters/" + name,
	}

	_, err := client.DeleteCluster(ctx, req)
	if err != nil {
		return errors.Errorf("Failed to delete cluster:", err)
	}
	fmt.Println("Cluster deleted")
	return nil
}

// FetchAccessToken prints an access token
func FetchAccessToken(keyfile string) (string, error) {
	ctx := context.Background()
	c, err := credentials.NewIamCredentialsClient(
		ctx, option.WithCredentialsFile(keyfile))
	if err != nil {
		return "", errors.Errorf("Failed to generate token:", err)
	}

	req := &credentialspb.GenerateAccessTokenRequest{
		Name:  "projects/-/serviceAccounts/maker-sa@review-287714.iam.gserviceaccount.com",
		Scope: []string{"https://www.googleapis.com/auth/cloud-platform"},
	}
	resp, err := c.GenerateAccessToken(ctx, req)
	if err != nil {
		return "", errors.Errorf("Failed to generate token:", err)
	}
	return resp.AccessToken, nil
}
