package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"os"

	"github.com/spf13/cobra"
)

// deleteVmCmd represents the delete vm command
var deleteVMCmd = &cobra.Command{
	Use:   "vm",
	Short: "deletes a VM",
	Long: `Used to delete a VM object on the specified provider:

Usage: maker delete vm -p do -n VM-NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			patToken, defaultRegion, err := do.LoadConfig()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			client := do.CreateDoClient(patToken, defaultRegion)
			err = do.Authenticate(client)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			dropletID, err := do.GetDoDroplet(client, name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = do.DeleteDoDroplet(client, dropletID, name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case "aws":
			defaultRegion := aws.LoadConfig()
			session := aws.CreateAwsSession(defaultRegion, aws.CredsPath)
			instanceID := aws.GetInstanceID(session, name)
			aws.DeleteEc2Instance(session, instanceID)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteVMCmd)

	// Local flags which will only run when this command
	deleteVMCmd.Flags().StringP("name", "n", "", "name of the object")
	deleteVMCmd.MarkFlagRequired("name")
}
