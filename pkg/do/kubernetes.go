package do

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

// CreateDoCluster creates a Kubernetes cluster on DigitalOcean
func CreateDoCluster(client *godo.Client, name, defaultRegion, nodeSize, version string, nodeCount int) error {
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
	_, _, err := client.Kubernetes.Create(ctx, req)
	if err != nil {
		return errors.Errorf("Creating cluster failed:", err)
	}
	fmt.Println("Cluster", name, "created")
	return nil
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
	fmt.Println("Node Info:")
	for i := 0; i < cluster.NodePools[0].Count; i++ {
		fmt.Printf(
			"Node Name: %s\nNode Status %s\n\n",
			cluster.NodePools[0].Nodes[i].Name,
			cluster.NodePools[0].Nodes[i].Status,
		)
	}
}

// DeleteDoCluster delets a droplet with the provided ID
func DeleteDoCluster(client *godo.Client, id int, name string) error {
	ctx := context.TODO()
	_, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		return errors.Errorf("Deleting droplet failed:", err)
	}
	fmt.Println("Cluster", name, "deleted")
	return nil
}

func FetchDoKubeConfig() {

}
