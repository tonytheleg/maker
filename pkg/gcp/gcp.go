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

// CreateGceInstance creates a compute instance with provided specs
func CreateGceInstance(computeService *compute.Service, name, project, zone, machineType, diskImage string) error {
	// make sure image is provided in proper format for GCP
	imageCheck := strings.Contains(diskImage, "/")
	if !imageCheck {
		err := errors.New("\nExample: 'ubuntu-os-cloud/ubuntu-1604-xenial-v20210119'")
		err = errors.Wrapf(err, "\nImage name must be provided in 'project/name' format")
		return err
	}
	s := strings.Split(diskImage, "/")
	imageProject, imageName := s[0], s[1]
	sourceImage := fmt.Sprintf("projects/%s/global/images/%s", imageProject, imageName)

	ctx := context.Background()
	image := compute.AttachedDiskInitializeParams{SourceImage: sourceImage}
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

	_, err := computeService.Instances.Insert(project, zone, rb).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to create GCE Instance: ", err)
	}
	fmt.Printf("Compute Instance %s is being created\n", name)
	return nil
}

// PrintInstanceStatus outputs instance info
func PrintInstanceStatus(computeService *compute.Service, name, project, zone string) error {
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
func DeleteGceInstance(computeService *compute.Service, name, project, zone string) error {
	ctx := context.Background()

	_, err := computeService.Instances.Delete(project, zone, name).Context(ctx).Do()
	if err != nil {
		return errors.Errorf("Failed to delete GCE Instance %s: ", name, err)
	}
	fmt.Printf("Instance %s has been deleted\n", name)
	return nil
}
