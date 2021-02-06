package do

import (
	"context"
	"fmt"
	"io/ioutil"
	"maker/pkg/utils"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

// CreateDoCluster creates a Kubernetes cluster on DigitalOcean
func CreateDoCluster(client *godo.Client, name, defaultRegion, nodeSize, version string, nodeCount int) (string, error) {
	ctx := context.TODO()
	req := &godo.KubernetesClusterCreateRequest{
		Name:        name,
		RegionSlug:  defaultRegion,
		VersionSlug: version,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{
				Name:  name + "-pool",
				Size:  nodeSize,
				Count: nodeCount,
			},
		},
	}
	cluster, _, err := client.Kubernetes.Create(ctx, req)
	if err != nil {
		return "", errors.Errorf("Creating cluster failed:", err)
	}
	fmt.Println("Cluster", name, "created")
	return cluster.ID, nil
}

// GetDoCluster grabs the droplet ID with the provided name
func GetDoCluster(client *godo.Client, name string) (string, error) {
	var clusterID string
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	clusters, _, err := client.Kubernetes.List(ctx, opt)
	if err != nil {
		return "", errors.Wrapf(err, "Could not list clusters to search for %s:", name)
	}
	for index := range clusters {
		if clusters[index].Name == name {
			clusterID = clusters[index].ID
		}
	}
	if clusterID != "" {
		return clusterID, nil
	}
	return "", errors.Wrapf(err, "Could not find cluster with name %s:", name)
}

// PrintClusterStatus outputs some droplet info
func PrintClusterStatus(client *godo.Client, id string) {
	ctx := context.TODO()
	cluster, _, err := client.Kubernetes.Get(ctx, id)
	if err != nil {
		fmt.Println("Could not fetch cluster status:", err)
	}
	fmt.Printf("Cluster Info\n------------\n")
	fmt.Printf(
		"Name: %s\nUID: %s\nRegion: %s\nVersion: %s\n\nCluster Subnet: %s\nService Subnet: %s\nPublic IP: %s\nEndpoint: %s\n\nNode Size: %s\nNode Count: %d\nCreated: %s\n\n",
		cluster.Name,
		cluster.ID,
		cluster.RegionSlug,
		cluster.VersionSlug,
		cluster.ClusterSubnet,
		cluster.ServiceSubnet,
		cluster.IPv4,
		cluster.Endpoint,
		cluster.NodePools[0].Size,
		cluster.NodePools[0].Count,
		cluster.CreatedAt,
	)
	fmt.Printf("Node Info\n---------\n")
	for i := 0; i < cluster.NodePools[0].Count; i++ {
		fmt.Printf(
			"Node Name: %s\nNode Status %s\n\n",
			cluster.NodePools[0].Nodes[i].Name,
			cluster.NodePools[0].Nodes[i].Status.State,
		)
	}
}

// DeleteDoCluster delets a droplet with the provided ID
func DeleteDoCluster(client *godo.Client, id string, name string) error {
	ctx := context.TODO()
	_, err := client.Kubernetes.Delete(ctx, id)
	if err != nil {
		return errors.Errorf("Deleting cluster failed:", err)
	}
	fmt.Println("Cluster", name, "deleted")
	return nil
}

// FetchDoKubeConfig fetches a kubeconfig file for the cluster and writes it to .makers config directory
func FetchDoKubeConfig(client *godo.Client, id string) error {
	ctx := context.TODO()
	config, _, err := client.Kubernetes.GetKubeConfig(ctx, id)
	if err != nil {
		return errors.Errorf("Fetching kubeconfig failed:", err)
	}
	kubeConfigFile := string(config.KubeconfigYAML)
	err = ioutil.WriteFile(utils.ConfigFolderPath+"/do_kubeconfig", []byte(kubeConfigFile), 0755)
	if err != nil {
		return errors.Errorf("Failed to write kubeconfig:", err)
	}
	fmt.Println("Kubeconfig file written to", utils.ConfigFolderPath)
	fmt.Printf("To use the kubeconfig, be sure to run 'export KUBECONFIG=%s/do_kubeconfig'\n", utils.ConfigFolderPath)
	return nil
}
