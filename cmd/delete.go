package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "deletes the specified object on the specified platform",
	Long: `Used to delete various objects on the cloud provider specificed:

Example: maker delete vm [flags]`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	// adds the "vm" subcommand (maker delete vm)
	deleteCmd.AddCommand(deleteVMCmd)
	deleteCmd.AddCommand(deleteBucketCmd)
}
