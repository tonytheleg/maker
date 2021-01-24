package cmd

import (
	"fmt"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// deleteBucketCmd represents the deleteBucket command
var deleteBucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "deletes a storage bucket",
	Long: `Used to delete a storage bucket on the specified provider:

Usage: maker delete bucket -p PROVIDER -n BUCKET-NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)
			accessKey := config.SpacesAccessKey
			secretKey := config.SpacesSecretKey
			endpoint := config.SpacesDefaultEndpoint

			client := do.CreateDoSpacesClient(accessKey, secretKey, endpoint)
			err = do.DeleteDoSpace(client, name)
			utils.HandleErr("Failed to create Space:", err)
		case "aws":
			fmt.Println("aws", name)
		case "gcp":
			fmt.Println("gcp")
		case "azure":
			fmt.Println("azure")
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteBucketCmd)

	// Local flags which will only run when this command
	deleteBucketCmd.Flags().StringP("name", "n", "", "name of the object")
	deleteBucketCmd.MarkFlagRequired("name")
}
