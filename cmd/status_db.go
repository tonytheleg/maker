package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/gcp"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// createDbCmd represents the createDb command
var statusDbCmd = &cobra.Command{
	Use:     "db",
	Short:   "gets the status of a database",
	Long:    `Used to get the status of a Postgres database on the specified provider`,
	Example: "maker status db --provider {do|aws|gcp} --name NAME",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion

			client := do.CreateDoClient(patToken, defaultRegion)

			dropletID, err := do.GetDoDatabase(client, name)
			utils.HandleErr("Faiiled to fetch droplet ID:", err)

			do.PrintDatabaseStatus(client, dropletID)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			aws.PrintRdsStatus(session, name)
		case "gcp":
			keyfile, defaultZone, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			service, err := gcp.CreateSQLService(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			err = gcp.PrintSQLDbStatus(service, name, gcpProject, defaultZone)
			utils.HandleErr("Failed to fetch GCE instance:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	statusCmd.AddCommand(statusDbCmd)

	// Local flags which will only run when this command
	statusDbCmd.Flags().StringP("name", "n", "", "name of the VM")
	statusDbCmd.MarkFlagRequired("name")
}
