package cmd

import (
	"maker/pkg/do"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates the specified object on the specified platform",
	Long: `Used to create various objects on the cloud provider specificed:

create vm --provider [do, linode, aws, azure, gcp] --size small
create cluster --provider [do, linode, aws, azure, gcp] --cluster-size small`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")

		// need a switch but for now just look at DO
		if provider == "do" {
			patToken, defaultRegion := do.LoadConfig()
			client := do.CreateDoClient(patToken, defaultRegion)
			do.Authenticate(client)
			do.CreateDoDroplet(client, "test-vm", defaultRegion, "s-1vcpu-1gb", "ubuntu-16-04-x64")
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringP("provider", "p", "", "sets the cloud provider")

}
