package do

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
	"golang.org/x/crypto/ssh/terminal"
)

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() {
	// ask for PAT token
	fmt.Printf("Configuring access and defaults...\n\n")
	fmt.Print("Enter PAT Token: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	os.Setenv("DO_PAT_TOKEN", string(pass))
	println()

	// ask for default region
	var region string
	fmt.Print("Default Region: ")
	fmt.Scanln(&region)
	os.Setenv("DO_DEFAULT_REGION", string(region))
	println()
}

// CreateDoClient creates a new client to interact with Digital Ocean
func CreateDoClient(patToken string) *godo.Client {
	client := godo.NewFromToken(patToken)
	return client
}

// Authenticate gets account info and prints it
func Authenticate(client *godo.Client) {
	ctx := context.TODO()
	account, _, err := client.Account.Get(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfully Authenticated")
	fmt.Printf("Account Email: %s\nDefault Region: %s\n", account.Email, os.Getenv("DO_DEFAULT_REGION"))
}

// CreateDoDroplet creates a droplet with provided specs
func CreateDoDroplet(client godo.Client, name string, region string, sizeSlug string, imageSlug string) {
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
