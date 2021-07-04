package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

	"github.com/spf13/cobra"
)

// statusClusterCmd represents the statusCluster command
var statusClusterCmd = &cobra.Command{
	Use:     "cluster",
	Short:   "gets the status of a Kubernetes cluster",
	Long:    `Used to fetch information about a Kubernetes clusters on the specified provider`,
	Example: "maker status cluster --provider {do|aws|gcp} --name CLUSTER-NAME",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		getConfig, _ := cmd.Flags().GetBool("fetch-kubeconfig")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion
			client := do.CreateDoClient(patToken, defaultRegion)
			utils.HandleErr("Failed to authenticate:", err)

			clusterID, err := do.GetDoCluster(client, name)
			do.PrintClusterStatus(client, clusterID)
			if getConfig {
				do.FetchDoKubeConfig(client, clusterID)
			}
			utils.HandleErr("Failed to create cluster:", err)
		case "aws":
			defaultRegion, err := aws.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			session, err := aws.CreateAwsSession(aws.CredsPath, defaultRegion)
			utils.HandleErr("Failed to setup AWS Session:", err)

			err = aws.PrintEksClusterStatus(session, name, name+"-nodegroup")
			utils.HandleErr("Failed to get EKS cluster status:", err)
			if getConfig {
				result, err := aws.GetCluster(session, name)
				utils.HandleErr("Failed to grab cluster info:", err)

				err = aws.CreateKubeconfig(*result.Cluster.Endpoint, *result.Cluster.CertificateAuthority.Data, name)
				utils.HandleErr("Failed to grab cluster info:", err)
			}
		case "gcp":
			keyfile, defaultZone, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := gcp.CreateGkeClient(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			gcp.PrintGkeClusterStatus(client, name, gcpProject, defaultZone)
			utils.HandleErr("Failed to fetch GCE instance:", err)
			if getConfig {
				accessToken, err := gcp.FetchAccessToken(keyfile)
				utils.HandleErr("Failed to fetch token", err)

				err = gcp.CreateKubeconfig(client, name, gcpProject, defaultZone, accessToken)
				utils.HandleErr("Failed to create kubeconfig", err)
			}
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	statusCmd.AddCommand(statusClusterCmd)

	statusClusterCmd.Flags().StringP("name", "n", "", "name of the cluster")
	statusClusterCmd.MarkFlagRequired("name")
	statusClusterCmd.Flags().BoolP("fetch-kubeconfig", "k", false, "fetches the Kubeconfig while checking status")
}
