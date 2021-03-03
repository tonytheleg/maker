package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createDbCmd represents the createDb command
var deleteDbCmd = &cobra.Command{
	Use:     "db",
	Short:   "deletes a database",
	Long:    `Used to delete a Postgres database on the specified provider`,
	Example: "maker delete db --provider {do|aws|gcp} --name NAME",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			fmt.Println("create DB called", provider, name)
		case "aws":
			fmt.Println("create DB called", provider, name)
		case "gcp":
			fmt.Println("create DB called", provider, name)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteDbCmd)

	// Local flags which will only run when this command
	deleteDbCmd.Flags().StringP("name", "n", "", "name of the VM")
	deleteDbCmd.MarkFlagRequired("name")
}
