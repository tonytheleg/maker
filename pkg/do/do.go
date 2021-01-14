package do

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
)

// CreateDoClient creates a new client to interact with Digital Ocean
func CreateDoClient(patToken, defaultRegion string) *godo.Client {
	client := godo.NewFromToken(patToken)
	return client
}

// Authenticate gets account info and prints it
func Authenticate(client *godo.Client) {
	ctx := context.TODO()
	_, _, err := client.Account.Get(ctx)
	if err != nil {
		fmt.Println("Failed to create context -- ", err)
	}
}

// CreateDoDroplet creates a droplet with provided specs
func CreateDoDroplet(client *godo.Client, name string, region string, sizeSlug string, imageSlug string) {
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
		fmt.Println("Failed to create droplet -- ", err)
	}
	fmt.Println(droplet.Name, "created")
}

// GetDoDroplet grabs the droplet ID with the provided name
func GetDoDroplet(client *godo.Client, name string) int {
	var dropletID int
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		fmt.Println("Failed to fetch droplet id -- ", err)
		os.Exit(1)
	}
	for index := range droplets {
		if droplets[index].Name == name {
			dropletID = droplets[index].ID
		}
	}
	if dropletID != 0 {
		return dropletID
	}
	fmt.Println("Could not find droplet with that ID")
	os.Exit(1)
	return 0
}

// PrintDropletStatus outputs some droplet info
func PrintDropletStatus(client *godo.Client, id int) {
	ctx := context.TODO()
	droplet, _, err := client.Droplets.Get(ctx, id)
	if err != nil {
		fmt.Println("Could not fetch droplet status -- ", err)
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
func DeleteDoDroplet(client *godo.Client, id int, name string) {
	ctx := context.TODO()
	_, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		fmt.Println("Deleting droplet failed -- ", err)
		os.Exit(1)
	} else {
		fmt.Println("Droplet", name, "deleted")
	}
}
