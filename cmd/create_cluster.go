package cmd

import (
	"fmt"
	"maker/pkg/aws"
	"maker/pkg/do"
	"maker/pkg/utils"

	"github.com/spf13/cobra"
)

// createClusterCmd represents the createCluster command
var createClusterCmd = &cobra.Command{
	Use:     "cluster",
	Short:   "creates a Kubernetes cluster",
	Long:    `Used to create a Kubernetes cluster on the specified provider`,
	Example: "maker create cluster --provider {do|aws|gcp} --size SIZE --name CLUSTER-NAME",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		nodeSize, _ := cmd.Flags().GetString("node-size")
		nodeCount, _ := cmd.Flags().GetInt("node-count")
		version, _ := cmd.Flags().GetString("version")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion
			client := do.CreateDoClient(patToken, defaultRegion)
			utils.HandleErr("Failed to authenticate:", err)

			clusterID, err := do.CreateDoCluster(client, name, defaultRegion, nodeSize, version, nodeCount)
			utils.HandleErr("Failed to create cluster:", err)

			err = do.FetchDoKubeConfig(client, clusterID)
			utils.HandleErr("Failed to create cluster:", err)

		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			arn, err := aws.GetExistingRoleARN(session)
			fmt.Println(arn)
			if arn == "" {
				arn, err = aws.CreateEksClusterRole(session)
				utils.HandleErr("Failed to create EKS service linked role:", err)
			}
		case "gcp":
			fmt.Println("create cluster gcp called", name, nodeSize, nodeCount, version)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)

	createClusterCmd.Flags().StringP("name", "n", "", "name of the cluster")
	createClusterCmd.MarkFlagRequired("name")
	createClusterCmd.Flags().StringP("node-size", "s", "", "sets the node VM size/Instance type")
	createClusterCmd.MarkFlagRequired("node-size")
	createClusterCmd.Flags().IntP("node-count", "c", 2, "sets the node pool size (default 1)")
	createClusterCmd.Flags().StringP("version", "v", "", "sets the Kubernetes/Vendor version")
	createClusterCmd.MarkFlagRequired("version")
}
