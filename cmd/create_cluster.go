/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
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

			err = do.CreateDoCluster(client, name, defaultRegion, nodeSize, version, nodeCount)
			utils.HandleErr("Failed to create cluster:", err)
		case "aws":
			fmt.Println("create cluster aws called", name, nodeSize, nodeCount, version)
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
	createClusterCmd.Flags().IntP("node-count", "c", 1, "sets the node pool size (default 1)")
	createClusterCmd.Flags().StringP("version", "v", "", "sets the Kubernetes/Vendor version")
	createClusterCmd.MarkFlagRequired("version")
}
