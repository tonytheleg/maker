package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

	"github.com/spf13/cobra"
)

// statusBucketCmd represents the statusBucket command
var statusBucketCmd = &cobra.Command{
	Use:     "bucket",
	Short:   "fetches basic bucket info",
	Long:    `Confirms the bucket exists and provides minimal info for each provider`,
	Example: "maker status bucket --provider {do|aws|gcp} --name BUCKET-NAME",
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
			err = do.GetDoSpaceInfo(client, name)
			utils.HandleErr("Failed to fetch Space:", err)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := aws.CreateS3Client(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to create client", err)

			err = aws.GetS3BucketInfo(client, name)
			utils.HandleErr("Failed to fetch S3 bucket:", err)
		case "gcp":
			keyfile, _, _, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := gcp.CreateStorageClient(keyfile)
			utils.HandleErr("Failed to create a Storage client:", err)

			err = gcp.GetStorageBucketInfo(client, name)
			utils.HandleErr("Failed to create Storage bucket:", err)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	statusCmd.AddCommand(statusBucketCmd)

	statusBucketCmd.Flags().StringP("name", "n", "", "name of the bucket")
	statusBucketCmd.MarkFlagRequired("name")
}
