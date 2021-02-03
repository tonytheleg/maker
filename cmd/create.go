package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [object]",
	Short: "creates the specified object on the specified platform",
	Long:  `Used to create various objects on the cloud provider specified`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
