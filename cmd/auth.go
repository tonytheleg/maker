package cmd

import (
	"maker/do"
	"os"

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
			do.CreateDoClient(os.Getenv("DO_PAT_TOKEN"))
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	authCmd.Flags().StringP("provider", "p", "", "sets the cloud provider")
}
