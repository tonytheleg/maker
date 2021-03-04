package cmd

import (
	"fmt"
	"maker/internal/aws"
	"maker/internal/do"
	"maker/internal/gcp"
	"maker/internal/utils"

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

			err = aws.DeleteEksNodeGroup(session, name, name+"-nodegroup")
			utils.HandleErr("Failed to delete the node group:", err)

			err = aws.DeleteEksCluster(session, name, name+"-nodegroup")
			utils.HandleErr("Failed to delete the cluster", err)
		case "gcp":
			keyfile, defaultZone, gcpProject, err := gcp.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			client, err := gcp.CreateGkeClient(keyfile)
			utils.HandleErr("Failed to create a Compute Service:", err)

			err = gcp.DeleteGkeCluster(client, name, gcpProject, defaultZone)
			utils.HandleErr("Failed to create GCE instance:", err)
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

/*
[tony@batcave ~]$ gcloud auth activate-service-account maker-sa@review-287714.iam.gserviceaccount.com  --key-file=/home/tony/.maker/makersa.json
Activated service account credentials for: [maker-sa@review-287714.iam.gserviceaccount.com]
[tony@batcave ~]$ gcloud auth --help
[tony@batcave ~]$ gcloud auth print-access-token
ya29.c.KqQB9AcGt7gpAbFXPJfB-rp3o7O0WmGSI_TSrFdzu2Xm8kKckPQb7eOV14p4BNvC84X_cH9Vo4rgussbYiiR3Ers21Z-lyjPpYzX1qzGCVXIMYzt8mtmmp9dJ_EeNcxYtquoVNfjTdaVDek9A7QmkGpO6JfPnmnZML8JxtPmkPZMQ2kqb2ed_hJSsjeQ3UZ84yN5LqH_pNuw5tFV1vJbIp5NEBK1AkI
*/
