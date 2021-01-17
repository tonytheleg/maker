package test

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	jsonPath := "/home/tony/Downloads/review-287714-4df47aacb805.json"
	ctx := context.Background()

	fmt.Println("next is creating compute service")
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}
	// Project ID for this request.
	project := "review-287714" // TODO: Update placeholder value.

	// The name of the zone for this request.
	zone := "us-east1-b" // TODO: Update placeholder value.

	image := compute.AttachedDiskInitializeParams{
		SourceImage: "projects/ubuntu-os-cloud/global/images/ubuntu-1604-xenial-v20210112",
	}

	nic := &compute.NetworkInterface{}

	nics := []*compute.NetworkInterface{nic}

	disk := &compute.AttachedDisk{
		Boot:             true,
		InitializeParams: &image,
		DiskSizeGb:       10,
	}

	disks := []*compute.AttachedDisk{disk}

	rb := &compute.Instance{
		MachineType:       "projects/review-287714/zones/us-east1-b/machineTypes/e2-micro",
		Disks:             disks,
		Name:              "test-vm",
		NetworkInterfaces: nics,
	}

	fmt.Println("next is creating instance")

	resp, err := computeService.Instances.Insert(project, zone, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%#v\n", resp)
}

//	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
//	if err != nil {
//			log.Fatal(err)
//	}
//
//	computeService, err := compute.New(c)
//	if err != nil {
//			log.Fatal(err)
//	}

//	// Project ID for this request.
//	project := "review-287714" // TODO: Update placeholder value.
//
//	// The name of the zone for this request.
//	zone := "us-east1-b" // TODO: Update placeholder value.
//
//	rb := &compute.Instance{
//			// TODO: Add desired fields of the request body.
//	}
//
//	resp, err := computeService.Instances.Insert(project, zone, rb).Context(ctx).Do()
//	if err != nil {
//			log.Fatal(err)
//	}
//
//	// TODO: Change code below to process the `resp` object:
//	fmt.Printf("%#v\n", resp)
