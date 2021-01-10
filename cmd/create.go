package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [object]",
	Short: "creates the specified object on the specified platform",
	Long: `Used to create various objects on the cloud provider specificed:

Options: vm, cluster, bucket, db
Example: maker create vm [flags]`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	// adds the "vm" subcommand (maker create vm)
	createCmd.AddCommand(createVmCmd)
}
