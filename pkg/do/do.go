package do

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

// CreateDoClient creates a new client to interact with Digital Ocean
func CreateDoClient(patToken, defaultRegion string) *godo.Client {
	client := godo.NewFromToken(patToken)
	return client
}

// Authenticate gets account info and prints it
func Authenticate(client *godo.Client) error {
	ctx := context.TODO()
	_, _, err := client.Account.Get(ctx)
	if err != nil {
		return errors.Errorf("Failed to create context -- ", err)
	}
	return nil
}

// CreateDoDroplet creates a droplet with provided specs
func CreateDoDroplet(client *godo.Client, name string, region string, sizeSlug string, imageSlug string) error {
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
		return errors.Errorf("Failed to create droplet -- ", err)
	}
	fmt.Println(droplet.Name, "created")
	return nil
}

// GetDoDroplet grabs the droplet ID with the provided name
func GetDoDroplet(client *godo.Client, name string) (int, error) {
	var dropletID int
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	droplets, _, err := client.Droplets.List(ctx, opt)
	if err != nil {
		return 1, errors.Wrapf(err, "Could not list droplets to search for %s", name)
	}
	for index := range droplets {
		if droplets[index].Name == name {
			dropletID = droplets[index].ID
		}
	}
	if dropletID != 0 {
		return dropletID, nil
	}
	return 1, errors.Wrapf(err, "Could not find droplet with name %s", name)

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
func DeleteDoDroplet(client *godo.Client, id int, name string) error {
	ctx := context.TODO()
	_, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		return errors.Errorf("Deleting droplet failed -- ", err)
	}
	fmt.Println("Droplet", name, "deleted")
	return nil
}
