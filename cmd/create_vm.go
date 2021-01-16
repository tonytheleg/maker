package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// createVmCmd represents the create vm command
var createVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "creates a VM",
	Long: `Used to create a VM object on the specified provider:

Usage: maker create vm -s s-1vcpu-1gb -i ubuntu-16-04-x64 -n test -p do`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetString("size")
		image, _ := cmd.Flags().GetString("image")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			patToken, defaultRegion, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client := do.CreateDoClient(patToken, defaultRegion)
			err = do.Authenticate(client)
			utils.HandleErr("Failed to authenticate:", err)

			err = do.CreateDoDroplet(client, name, defaultRegion, size, image)
			utils.HandleErr("Failed to create droplet:", err)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(defaultRegion, aws.CredsPath)
			utils.HandleErr("Failed to setup AWS Session:", err)

			aws.CreateEc2Instance(session, name, defaultRegion, image, size)
			utils.HandleErr("Failed to create EC2 instance:", err)
		case "gcp":
			fmt.Println("gcp")
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(createVMCmd)

	// Local flags which will only run when this command
	createVMCmd.Flags().StringP("name", "n", "", "name of the object")
	createVMCmd.MarkFlagRequired("name")
	createVMCmd.Flags().StringP("size", "s", "", "sets the size of the object")
	createVMCmd.MarkFlagRequired("size")
	createVMCmd.Flags().StringP("image", "i", "", "sets the image slug")
	createVMCmd.MarkFlagRequired("image")
}
