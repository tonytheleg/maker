package gcp

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// CreateGceService creates a new client to interact with Digital Ocean
func CreateGceService(keyfile string) (*compute.Service, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(
		ctx, option.WithCredentialsFile(keyfile))
	if err != nil {
		return nil, errors.Errorf("Failed to create client: ", err)
	}
	return computeService, nil
}

// CreateGceInstance creates a droplet with provided specs
func CreateGceInstance(computeService *compute.Service, name, zone, project, machineType, diskImage string) error {
	// image list https://console.cloud.google.com/compute/images
	ctx := context.Background()
	image := compute.AttachedDiskInitializeParams{SourceImage: diskImage}
	machineTypePath := fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", project, zone, machineType)
	fmt.Println(machineTypePath)
	nics := []*compute.NetworkInterface{new(compute.NetworkInterface)}

	disk := &compute.AttachedDisk{
		Boot:             true,
		InitializeParams: &image,
		DiskSizeGb:       10,
	}
	disks := []*compute.AttachedDisk{disk}

	rb := &compute.Instance{
		MachineType:       machineTypePath,
		Disks:             disks,
		Name:              name,
		NetworkInterfaces: nics,
	}

	resp, err := computeService.Instances.Insert(project, zone, rb).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to create GCE Instance: ", err)
	}
	fmt.Printf("%#v\n", resp)
	return nil
}

// PrintInstanceStatus outputs some droplet info
func PrintInstanceStatus(computeService *compute.Service, name, zone, project string) error {
	ctx := context.Background()

	resp, err := computeService.Instances.Get(project, zone, name).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to retreive GCE Instance %s: ", name, err)
	}
	fmt.Printf(
		"Name: %#v\nDistribution: %#v\n\nPublic IP: %#v\nRegion: %#v\nStatus: %#v\n",
		string(resp.Name),
		string(resp.Disks[0].InitializeParams.SourceImage),
		string(resp.NetworkInterfaces[0].NetworkIP),
		string(resp.Zone),
		string(resp.Status),
	)
	return nil
}

// DeleteGceInstance delets a droplet with the provided ID
func DeleteGceInstance(client *godo.Client, id int, name string) error {
	return nil
}
