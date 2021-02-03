package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/gcp"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// deleteBucketCmd represents the deleteBucket command
var deleteBucketCmd = &cobra.Command{
	Use:     "bucket",
	Short:   "deletes a storage bucket",
	Long:    `Used to delete a storage bucket on the specified provider`,
	Example: "maker delete bucket --provider {do|aws|gcp} --name BUCKET-NAME",
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
			err = do.DeleteSpaceObjects(client, name)
			utils.HandleErr("Failed to delete Space:", err)

			err = do.DeleteDoSpace(client, name)
			utils.HandleErr("Failed to delete Space:", err)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := aws.CreateS3Client(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to create sessions", err)

			err = aws.DeleteS3Objects(client, name)
			utils.HandleErr("Failed to delete bucket", err)

			err = aws.DeleteS3Bucket(client, name)
			utils.HandleErr("Failed to deletes S3 bucket:", err)
		case "gcp":
			keyfile, _, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := gcp.CreateStorageClient(keyfile)
			utils.HandleErr("Failed to create a Storage client:", err)

			err = gcp.DeleteStorageObjects(client, name, gcpProject)
			utils.HandleErr("Failed to delete objects in bucket", err)

			err = gcp.DeleteStorageBucket(client, name, gcpProject)
			utils.HandleErr("Failed to create Storage bucket:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteBucketCmd)

	deleteBucketCmd.Flags().StringP("name", "n", "", "name of the bucket")
	deleteBucketCmd.MarkFlagRequired("name")
}
