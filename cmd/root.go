package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "maker",
	Short: "creates various cloud services on multiple cloud platforms",
	Long: `Maker can be used to create various types of services in various cloud providers such as VM's,
K8s clusters, storage buckets, etc. Its not meant to be a full replacement for each 
providers own CLI's or clients. Handy for spinning up and down infra for labs and devlopment work kinda thing.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("provider", "p", "", "sets the cloud provider")
	rootCmd.MarkPersistentFlagRequired("provider")

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
/*
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".maker" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".maker")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
*/
