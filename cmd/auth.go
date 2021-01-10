package cmd

import (
	"maker/pkg/do"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "configures authentication to specified provider",
	Long: `auth is used to set the required config files needed to
authenticate and communicate with a cloud provider

maker auth --provider [do, aws, azure, gcp]
Required settings will be prompted based on provider`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")

		// need a switch but for now just look at DO
		if provider == "do" {
			do.Configure()
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// local flags
	authCmd.Flags().StringP("provider", "p", "", "sets the cloud provider")
	authCmd.MarkFlagRequired("provider")
}
