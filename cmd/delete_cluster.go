package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// deleteClusterCmd represents the deleteCluster command
var deleteClusterCmd = &cobra.Command{
	Use:     "cluster",
	Short:   "deletes a Kubernetes cluster",
	Long:    `Deletes a Kubernetes cluster on the specified provider`,
	Example: "maker status cluster --provider {do|aws|gcp} --name CLUSTER-NAME",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion
			client := do.CreateDoClient(patToken, defaultRegion)
			utils.HandleErr("Failed to authenticate:", err)

			clusterID, err := do.GetDoCluster(client, name)
			err = do.DeleteDoCluster(client, clusterID, name)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			err = aws.DeleteEksCluster(session, name)
			utils.HandleErr("Failed to delete the cluster", err)
		case "gcp":
			fmt.Println("create cluster gcp called", name)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)

	deleteClusterCmd.Flags().StringP("name", "n", "", "name of the cluster")
	deleteClusterCmd.MarkFlagRequired("name")
}
