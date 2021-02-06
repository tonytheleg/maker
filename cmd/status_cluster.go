package cmd

import (
	"fmt"
	"maker/pkg/do"
	"maker/pkg/utils"

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

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			config, err := do.LoadConfig()
			utils.HandleErr("Failed to load config:", err)

			patToken, defaultRegion := config.PatToken, config.DefaultRegion
			client := do.CreateDoClient(patToken, defaultRegion)
			utils.HandleErr("Failed to authenticate:", err)

			clusterID, err := do.GetDoCluster(client, name)
			do.PrintClusterStatus(client, clusterID)
			utils.HandleErr("Failed to create cluster:", err)
		case "aws":
			fmt.Println("create cluster aws called", name)
		case "gcp":
			fmt.Println("create cluster gcp called", name)
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	statusCmd.AddCommand(statusClusterCmd)

	statusClusterCmd.Flags().StringP("name", "n", "", "name of the cluster")
	statusClusterCmd.MarkFlagRequired("name")
}
