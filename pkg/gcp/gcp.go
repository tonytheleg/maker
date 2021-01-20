package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// CreateGceService creates a new client to interact with GCP
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
	publicNic := &compute.AccessConfig{
		Name:        "External NAT",
		NetworkTier: "PREMIUM",
		Type:        "ONE_TO_ONE_NAT",
	}
	publicNics := []*compute.AccessConfig{publicNic}

	nic := &compute.NetworkInterface{
		AccessConfigs: publicNics,
	}
	nics := []*compute.NetworkInterface{nic}

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
	os := strings.Split(resp.Disks[0].Licenses[0], "/")
	currentZone := strings.Split(resp.Zone, "/")
	fmt.Printf(
		"Name: %s\nDistribution: %s\n\nPublic IP: %s\nZone: %s\nStatus: %s\n",
		resp.Name,
		os[len(os)-1],
		resp.NetworkInterfaces[0].AccessConfigs[0].NatIP,
		currentZone[len(currentZone)-1],
		resp.Status,
	)
	return nil
}

// DeleteGceInstance delets a droplet with the provided ID
func DeleteGceInstance(computeService *compute.Service, name, zone, project string) error {
	ctx := context.Background()

	_, err := computeService.Instances.Delete(project, zone, name).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to delete GCE Instance %s: ", name, err)
	}
	fmt.Printf("Instance %s has been deleted\n", name)
	return nil
}
