package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"

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
			patToken, defaultRegion := do.LoadConfig()
			client := do.CreateDoClient(patToken, defaultRegion)
			do.Authenticate(client)
			dropletID := do.GetDoDroplet(client, name)
			do.PrintDropletStatus(client, dropletID)
		case "aws":
			defaultRegion := aws.LoadConfig()
			session := aws.CreateAwsSession(defaultRegion, aws.CredsPath)
			aws.PrintEc2Status(session, name)
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
