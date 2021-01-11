package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/digitalocean/godo"
)

// CreateDoClient creates a new client to interact with Digital Ocean
func CreateAwsClient(patToken, defaultRegion string) *godo.Client {
	client := godo.NewFromToken(patToken)
	return client
}

// Authenticate gets account info and prints it
func Authenticate(client *godo.Client) {
	ctx := context.TODO()
	_, _, err := client.Account.Get(ctx)
	if err != nil {
		panic(err)
	}
}

// CreateDoDroplet creates a droplet with provided specs
func CreateEc2Instance(client *godo.Client, name string, region string, sizeSlug string, imageSlug string) {
	ctx := context.TODO()
	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   sizeSlug,
		Image: godo.DropletCreateImage{
			Slug: imageSlug,
		},
	}

	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		fmt.Println(err)
	}
	fmt.Println(droplet.Name, "created")
}

// GetDoDroplet grabs the droplet ID with the provided name
func GetEc2Instance(client *godo.Client, name string) (int, error) {
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		panic(err)
	}
	for index := range droplets {
		if droplets[index].Name == name {
			return droplets[index].ID, nil
		}
	}
	return 1, errors.New("Droplet not found with that name")
}

// PrintDropletStatus outputs some droplet info
func PrintEc2Status(client *godo.Client, id int) {
	ctx := context.TODO()
	droplet, _, err := client.Droplets.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Name: %s\nUID: %d\nMemory: %d\nDisk: %d\n\nDistribution: %s\nVersion: %s\n\nPublic IP: %s\nRegion: %s\nStatus: %s\n",
		droplet.Name,
		droplet.ID,
		droplet.Memory,
		droplet.Disk,
		droplet.Image.Distribution,
		droplet.Image.Name,
		droplet.Networks.V4[1].IPAddress,
		droplet.Region.Slug,
		droplet.Status,
	)
}

// DeleteDoDroplet delets a droplet with the provided ID
func DeleteEc2Instance(client *godo.Client, id int) error {
	ctx := context.TODO()
	_, err := client.Droplets.Delete(ctx, id)
	return err
}
