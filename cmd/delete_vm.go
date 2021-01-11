package cmd

import (
	"fmt"
	"maker/pkg/do"

	"github.com/spf13/cobra"
)

// deleteVmCmd represents the delete vm command
var deleteVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "deletes a VM",
	Long: `Used to delete a VM object on the specified provider:

Usage: maker delete vm -p do -n VM-NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")

		if provider == "do" {
			patToken, defaultRegion := do.LoadConfig()
			client := do.CreateDoClient(patToken, defaultRegion)
			do.Authenticate(client)

			dropletId, err := do.GetDoDroplet(client, name)
			if err != nil {
				panic(err)
			}
			err = do.DeleteDoDroplet(client, dropletId)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("Droplet", name, "deleted")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteVmCmd)

	// Local flags which will only run when this command
	deleteVmCmd.Flags().StringP("name", "n", "", "name of the object")
	deleteVmCmd.MarkFlagRequired("name")
}
