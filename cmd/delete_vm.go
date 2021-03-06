package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

	"github.com/spf13/cobra"
)

// deleteVmCmd represents the delete vm command
var deleteVMCmd = &cobra.Command{
	Use:     "vm",
	Short:   "deletes a VM",
	Long:    `Used to delete a VM object on the specified provider`,
	Example: "maker delete vm --provider {do|aws|gcp} --name NAME",
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

			err = do.DeleteDoDroplet(client, dropletID, name)
			utils.HandleErr("Failed to delete droplet:", err)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			instanceID, err := aws.GetInstanceID(session, name)
			utils.HandleErr("Failed to fetch EC2 instance ID:", err)

			err = aws.DeleteEc2Instance(session, instanceID)
			utils.HandleErr("Failed to delete EC2 instance ID:", err)
		case "gcp":
			keyfile, defaultZone, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			service, err := gcp.CreateGceService(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			err = gcp.DeleteGceInstance(service, name, gcpProject, defaultZone)
			utils.HandleErr("Failed to create GCE instance:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteVMCmd)

	deleteVMCmd.Flags().StringP("name", "n", "", "name of the VM")
	deleteVMCmd.MarkFlagRequired("name")
}
