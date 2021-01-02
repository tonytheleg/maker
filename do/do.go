package do

import "fmt"

// CreateDoClient creates a new client to interact with Digital Ocean
func CreateDoClient() {
	fmt.Println("Testing call to CreateDoClient")
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
