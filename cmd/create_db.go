package cmd

import (
	"fmt"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// createDbCmd represents the createDb command
var createDbCmd = &cobra.Command{
	Use:   "db",
	Short: "creates a database",
	Long: `Used to create a Postgres database on the specified provider
Sizes and Image names are provider specific!`,
	Example: "maker create db --provider {do|aws|gcp} --size SIZE --name NAME",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetString("size")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion

			client := do.CreateDoClient(patToken, defaultRegion)

			err = do.CreateDoDatabase(client, name, size, defaultRegion)
			utils.HandleErr("Failed to create database:", err)
		case "aws":
			fmt.Println("create DB called", provider, name, size)
		case "gcp":
			fmt.Println("create DB called", provider, name, size)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	createCmd.AddCommand(createDbCmd)

	// Local flags which will only run when this command
	createDbCmd.Flags().StringP("name", "n", "", "name of the VM")
	createDbCmd.MarkFlagRequired("name")
	createDbCmd.Flags().StringP("size", "s", "", "sets the VM size/Instance type")
	createDbCmd.MarkFlagRequired("size")
}
