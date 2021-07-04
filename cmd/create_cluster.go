package cmd

import (
	"errors"
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

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
			subnets, _ := cmd.Flags().GetStringSlice("subnets")
			if len(subnets) < 2 {
				err := errors.New("Must provide two subnets to create cluster (-s)")
				utils.HandleErr("Failed to initiate:", err)
			}
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			arn, err := aws.GetExistingRoleARN(session)
			if arn == "" {
				arn, err = aws.CreateEksClusterRole(session)
				utils.HandleErr("Failed to create EKS service linked role:", err)
			}
			err = aws.CreateEksCluster(session, name, arn, version, subnets)
			utils.HandleErr("Failed to create EKS cluster:", err)

			err = aws.CreateEksNodeGroup(session, name, arn, nodeSize, nodeCount, subnets)
			utils.HandleErr("Failed to create EKS node group:", err)

			result, err := aws.GetCluster(session, name)
			utils.HandleErr("Failed to grab cluster info:", err)

			err = aws.CreateKubeconfig(*result.Cluster.Endpoint, *result.Cluster.CertificateAuthority.Data, name)
			utils.HandleErr("Failed to grab cluster info:", err)
		case "gcp":
			keyfile, defaultZone, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := gcp.CreateGkeClient(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			err = gcp.CreateGkeCluster(client, name, gcpProject, defaultZone, nodeSize, nodeCount)
			utils.HandleErr("Failed to create GKE Cluster:", err)

			accessToken, err := gcp.FetchAccessToken(keyfile)
			utils.HandleErr("Failed to fetch token", err)

			err = gcp.CreateKubeconfig(client, name, gcpProject, defaultZone, accessToken)
			utils.HandleErr("Failed to create kubeconfig", err)
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
	createClusterCmd.Flags().IntP("node-count", "c", 2, "sets the node pool size")
	createClusterCmd.Flags().StringP("version", "v", "", "sets the Kubernetes/Vendor version")
	createClusterCmd.MarkFlagRequired("version")
	createClusterCmd.Flags().StringSliceP("subnets", "b", nil, "comma separated list of 2 subnets to deploy to (AWS Requuired Only)")
}
