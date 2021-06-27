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

// CreateDoDroplet creates a droplet with provided specs
func CreateDoDroplet(client *godo.Client, name string, region string, sizeSlug string, imageSlug string) error {
	ctx := context.TODO()
	dropletKey := &godo.DropletCreateSSHKey{}
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	keys, _, err := client.Keys.List(ctx, opt)
	if err != nil {
		return errors.Errorf("Failed to create droplet:", err)
	}
	if len(keys) < 1 {
		fmt.Println("To access a DO Droplet an SSH Key is required")
		fmt.Println("Create an SSH Key and Upload and try again")
		fmt.Println("https://docs.digitalocean.com/products/droplets/how-to/add-ssh-keys/to-account/")
		return errors.Errorf("failed to create droplet: SSH Key required and none are avaiable")
	}
	if len(keys) > 1 {
		var sshkeyID int
		fmt.Println("Multiple SSH Keys found -- Which would you like to use?")
		for _, sshkeys := range keys {
			fmt.Printf("Name: %s  ID: %d\n", sshkeys.Name, sshkeys.ID)
		}
		fmt.Printf("Enter a Key ID (not name): ")
		fmt.Scanln(&sshkeyID)
		*dropletKey = godo.DropletCreateSSHKey{ID: sshkeyID}
	} else {
		fmt.Printf("Using SSH Key %s\n", keys[0].Name)
		sshkeyID := keys[0].ID
		*dropletKey = godo.DropletCreateSSHKey{ID: sshkeyID}
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   sizeSlug,
		Image: godo.DropletCreateImage{
			Slug: imageSlug,
		},
		SSHKeys: []godo.DropletCreateSSHKey{*dropletKey},
	}

	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return errors.Errorf("Failed to create droplet:", err)
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
		return 1, errors.Wrapf(err, "Could not list droplets to search for %s:", name)
	}
	for index := range droplets {
		if droplets[index].Name == name {
			dropletID = droplets[index].ID
		}
	}
	if dropletID != 0 {
		return dropletID, nil
	}
	return 1, errors.Wrapf(err, "Could not find droplet with name %s:", name)

}

// PrintDropletStatus outputs some droplet info
func PrintDropletStatus(client *godo.Client, id int) {
	ctx := context.TODO()
	droplet, _, err := client.Droplets.Get(ctx, id)
	if err != nil {
		fmt.Println("Could not fetch droplet status:", err)
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
		return errors.Errorf("Deleting droplet failed:", err)
	}
	fmt.Println("Droplet", name, "deleted")
	return nil
}
