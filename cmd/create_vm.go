package cmd

import (
	"maker/pkg/do"

	"github.com/spf13/cobra"
)

// createVmCmd represents the createVm command
var createVmCmd = &cobra.Command{
	Use:   "vm",
	Short: "creates a VM",
	Long: `Used to create a VM object on the specified provider:

Usage: maker create vm -s s-1vcpu-1gb -i ubuntu-16-04-x64 -n test -p do`,
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetString("size")
		image, _ := cmd.Flags().GetString("image")

		if provider == "do" {
			patToken, defaultRegion := do.LoadConfig()
			client := do.CreateDoClient(patToken, defaultRegion)
			do.Authenticate(client)
			do.CreateDoDroplet(client, name, defaultRegion, size, image)
		}
	},
}

func init() {
	rootCmd.AddCommand(createVmCmd)

	// Local flags which will only run when this command
	createVmCmd.PersistentFlags().StringP("name", "n", "", "name of the object")
	createVmCmd.MarkPersistentFlagRequired("name")
	createVmCmd.Flags().StringP("size", "s", "", "sets the size of the object")
	createVmCmd.MarkFlagRequired("size")
	createVmCmd.Flags().StringP("image", "i", "", "sets the image slug")
	createVmCmd.MarkFlagRequired("image")
}
