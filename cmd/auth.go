package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth [flags]",
	Short: "configures authentication to specified provider",
	Long: `auth is used to set the required config files needed to
authenticate and communicate with a cloud provider
Required settings will be prompted based on provider`,
	Example: "maker auth --provider {do|aws|gcp}",
	Run: func(cmd *cobra.Command, args []string) {
		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			err := do.SetupConfig()
			utils.HandleErr("Failed to setup configuration files:", err)
		case "aws":
			err := aws.Configure()
			utils.HandleErr("Failed to setup configuration files:", err)
		case "gcp":
			err := gcp.Configure()
			utils.HandleErr("Failed to setup configuration files:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
