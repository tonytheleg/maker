package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/gcp"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// statusVmCmd represents the status vm command
var statusVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "gets the status of a VM",
	Long: `Provides basic info and resource ID's for a VM

Usage: maker status vm -p PROVIDER -n VM-NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			patToken, defaultRegion, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client := do.CreateDoClient(patToken, defaultRegion)
			err = do.Authenticate(client)
			utils.HandleErr("Failed to authenticate:", err)

			dropletID, err := do.GetDoDroplet(client, name)
			utils.HandleErr("Faiiled to fetch droplet ID:", err)

			do.PrintDropletStatus(client, dropletID)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(defaultRegion, aws.CredsPath)
			utils.HandleErr("Failed to setup AWS Session:", err)

			aws.PrintEc2Status(session, name)
		case "gcp":
			keyfile, defaultRegion, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			service, err := gcp.CreateGceService(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			err = gcp.PrintInstanceStatus(service, name, defaultRegion, gcpProject)
			utils.HandleErr("Failed to create GCE instance:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusVMCmd)

	// Local flags which will only run when this command
	statusVMCmd.Flags().StringP("name", "n", "", "name of the object")
	statusVMCmd.MarkFlagRequired("name")
}
