package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"

	"github.com/spf13/cobra"
)

// statusVmCmd represents the status vm command
var statusVmCmd = &cobra.Command{
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

			dropletId, err := do.GetDoDroplet(client, name)
			if err != nil {
				panic(err)
			}
			do.PrintDropletStatus(client, dropletId)
		case "aws":
			aws.Configure()
		default:
			// freebsd, openbsd,
			// plan9, windows...
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusVmCmd)

	// Local flags which will only run when this command
	statusVmCmd.Flags().StringP("name", "n", "", "name of the object")
	statusVmCmd.MarkFlagRequired("name")
}
