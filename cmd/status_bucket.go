package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// statusBucketCmd represents the statusBucket command
var statusBucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "fetches basic bucket info",
	Long: `Confirms the bucket exists and provides minimal info for each provider:

Example: 
  maker status bucket -p PROVIDER -n BUCKET-NAME`,
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
			fmt.Println("gcp")
		case "azure":
			fmt.Println("azure")
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusBucketCmd)

	// Local flags which will only run when this command
	statusBucketCmd.Flags().StringP("name", "n", "", "name of the object")
	statusBucketCmd.MarkFlagRequired("name")
}
