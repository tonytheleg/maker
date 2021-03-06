package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

	"github.com/spf13/cobra"
)

// statusVmCmd represents the status vm command
var statusVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "gets the status of a VM",
	Long: `Provides basic info and resource ID's for a VM

Example: 
  maker status vm -p PROVIDER -n NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion

			client := do.CreateDoClient(patToken, defaultRegion)

			dropletID, err := do.GetDoDroplet(client, name)
			utils.HandleErr("Faiiled to fetch droplet ID:", err)

			do.PrintDropletStatus(client, dropletID)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			aws.PrintEc2Status(session, name)
		case "gcp":
			keyfile, defaultZone, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			service, err := gcp.CreateGceService(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			err = gcp.PrintInstanceStatus(service, name, gcpProject, defaultZone)
			utils.HandleErr("Failed to fetch GCE instance:", err)
		case "azure":
			fmt.Println("azure called")
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	statusCmd.AddCommand(statusVMCmd)

	// Local flags which will only run when this command
	statusVMCmd.Flags().StringP("name", "n", "", "name of the VM")
	statusVMCmd.MarkFlagRequired("name")
}
