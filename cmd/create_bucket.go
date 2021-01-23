package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createBucketCmd represents the createBucket command
var createBucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		size, _ := cmd.Flags().GetString("size")
		image, _ := cmd.Flags().GetString("image")

		switch provider, _ := cmd.Flags().GetString("provider"); provider {
		case "do":
			fmt.Println("do called")
		case "aws":
			fmt.Println("aws called")
		case "gcp":
			fmt.Println("gcp called")
		case "azure":
			fmt.Println("azure called")
		default:
			fmt.Printf("Unknown Provder -- %s", provider)
		}
	},
}

func init() {
	rootCmd.AddCommand(createBucketCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createBucketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createBucketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
