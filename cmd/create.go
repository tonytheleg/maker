package cmd

import (
	"fmt"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 && args[1] != "vm" {
			return fmt.Errorf("vm required -- need to know object to be created")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetString("size")
		image, _ := cmd.Flags().GetString("image")

		// need a switch but for now just look at DO
		if provider == "do" {
			patToken, defaultRegion := do.LoadConfig()
			client := do.CreateDoClient(patToken, defaultRegion)
			do.Authenticate(client)
			do.CreateDoDroplet(client, name, defaultRegion, size, image)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	createCmd.PersistentFlags().StringP("provider", "p", "", "sets the cloud provider")
	createCmd.MarkFlagRequired("provider")
	createCmd.PersistentFlags().StringP("name", "n", "", "name of the object")
	createCmd.MarkFlagRequired("name")
	createCmd.PersistentFlags().StringP("size", "s", "", "sets the size of the object")
	createCmd.MarkFlagRequired("size")
	createCmd.PersistentFlags().StringP("image", "i", "", "sets the image slug")
	createCmd.MarkFlagRequired("image")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//createCmd.Flags().StringP("provider", "p", "", "sets the cloud provider")
	createCmd.Flags().Bool("vm", false, "creates a vm")

}
