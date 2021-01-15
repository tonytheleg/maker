package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/utils"

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
		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			err := do.Configure()
			utils.HandleErr("Failed to setup configuration files", err)
		case "aws":
			err := aws.Configure()
			utils.HandleErr("Failed to setup configuration files", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// local flags
	authCmd.Flags().StringP("provider", "p", "", "sets the cloud provider")
	authCmd.MarkFlagRequired("provider")
}
