package do

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// ConfigFile makes up the required settings in a DO config file
type ConfigFile struct {
	PatToken      string `json:"pat_token"`
	DefaultRegion string `json:"region"`
}

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFolder := ".maker"
	configName := "do_config"
	configPath := filepath.Join(homeDir, configFolder, configName)

	// check that .maker exists
	_, err = os.Stat(configFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(configFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// check if config exists to create, or to verify
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		config := &ConfigFile{}
		err = CreateConfigFile(config, configPath)
		if err != nil {
			panic(err)
		}
		fmt.Println("Config file generated at", configPath)
	} else {
		ShowCurrentConfig(configPath)
	}
}

// ShowCurrentConfig prints out the current config file
func ShowCurrentConfig(configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nCurrent Config:\n\n%s", string(data))

	var confirmation string
	fmt.Print("Is this info still accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		config := &ConfigFile{}
		err = CreateConfigFile(config, configPath)
		if err != nil {
			panic(err)
		}
		fmt.Println("Config file generated at", configPath)
	}
}

// CreateConfigFile makes the config file to use in all DO commands
func CreateConfigFile(config *ConfigFile, configPath string) error {
	// ask for PAT token
	fmt.Println("Please authenticate using your Digital Ocean account...")
	fmt.Println("Tokens can be generated at https://cloud.digitalocean.com/account/api/tokens")
	fmt.Print("Enter PAT Token: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	config.PatToken = string(pass)
	println()

	// ask for default region
	var region string
	fmt.Print("Default Region: ")
	fmt.Scanln(&region)
	config.DefaultRegion = string(region)
	println()

	viper.SetConfigType("yaml")
	viper.Set("pat_token", config.PatToken)
	viper.Set("default_region", config.DefaultRegion)
	return viper.WriteConfigAs(configPath)
}
