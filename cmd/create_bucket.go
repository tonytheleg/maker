package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

	"github.com/spf13/cobra"
)

// createBucketCmd represents the createBucket command
var createBucketCmd = &cobra.Command{
	Use:     "bucket",
	Short:   "creates a storage bucket",
	Long:    `Used to create a storage bucket on the specified provider`,
	Example: "maker create bucket --provider {do|aws|gcp} --name BUCKET-NAME (Must be globally unique!)",
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
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := aws.CreateS3Client(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to create client", err)

			err = aws.CreateS3Bucket(client, name)
			utils.HandleErr("Failed to create S3 bucket:", err)
		case "gcp":
			keyfile, _, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := gcp.CreateStorageClient(keyfile)
			utils.HandleErr("Failed to create a Storage client:", err)

			err = gcp.CreateStorageBucket(client, name, gcpProject)
			utils.HandleErr("Failed to create Storage bucket:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	createCmd.AddCommand(createBucketCmd)

	// Local flags which will only run when this command
	createBucketCmd.Flags().StringP("name", "n", "", "name of the bucket")
	createBucketCmd.MarkFlagRequired("name")
}
