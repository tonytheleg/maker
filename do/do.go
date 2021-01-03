package do

import (
	"fmt"
	"os"

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

	// test
	fmt.Printf("Region: %s  Password: %s\n", region, pass)
	fmt.Printf("DO_DEFAULT_REGION: %s  DO_PAT_TOKEN: %s\n", os.Getenv("DO_DEFAULT_REGION"), os.Getenv("DO_PAT_TOKEN"))
}

// CreateDoClient creates a new client to interact with Digital Ocean
func CreateDoClient(patToken string) {
	fmt.Println("Testing call to CreateDoClient")
	fmt.Println("PAT Token:", patToken)
	//	client := godo.NewFromToken(os.Getenv("PAT_TOKEN"))
	//
	//	ctx := context.TODO()
	//	createRequest := &godo.DropletCreateRequest{
	//		Name:   "do-api-gotest",
	//		Region: "nyc3",
	//		Size:   "s-1vcpu-1gb",
	//		Image: godo.DropletCreateImage{
	//			Slug: "ubuntu-16-04-x64",
	//		},
	//	}
	//
	//	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	//	if err != nil {
	//		fmt.Printf("Something bad happened: %s\n\n", err)
	//		fmt.Println(err)
	//	}
	//	fmt.Println(droplet.Name, "created")
}
