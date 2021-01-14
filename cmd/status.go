package cmd

import (
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "gets the status of the specified object on the specified platform",
	Long: `Provides basic information about objects on the cloud provider specificed:

Options: vm, cluster, bucket, db
Example: maker status vm [flags]`,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	// adds the "vm" subcommand (maker status vm)
	statusCmd.AddCommand(statusVMCmd)
}
