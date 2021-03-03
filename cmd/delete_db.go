package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/utils"

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
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion

			client := do.CreateDoClient(patToken, defaultRegion)

			databaseID, err := do.GetDoDatabase(client, name)
			utils.HandleErr("Faiiled to fetch droplet ID:", err)

			err = do.DeleteDoDatabase(client, databaseID, name)
			utils.HandleErr("Failed to delete droplet:", err)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			err = aws.DeleteRdsInstance(session, name)
			utils.HandleErr("Failed to delete EC2 instance ID:", err)
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
