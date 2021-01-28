package cmd

import (
	"fmt"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// createBucketCmd represents the createBucket command
var createBucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "creates a storage bucket",
	Long: `Used to create a storage bucket on the specified provider:

Example: 
  maker create bucket -p PROVIDER -n BUCKET-NAME (Must be globally unique!)`,
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
			err = do.CreateDoSpace(client, name)
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
	rootCmd.AddCommand(createBucketCmd)

	// Local flags which will only run when this command
	createBucketCmd.Flags().StringP("name", "n", "", "name of the object")
	createBucketCmd.MarkFlagRequired("name")
}
