package do

import (
	"fmt"
	"io/ioutil"
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

// HomeDir stores the path of the current users Home directory
var HomeDir, _ = os.UserHomeDir()

// ConfigFolder is the name of Makers config folder stored in Home
var ConfigFolder = ".maker"

// ConfigName is the name of the config file used by Maker
var ConfigName = "do_config"

// ConfigPath is the full path to the ConfigFile
var ConfigPath = filepath.Join(HomeDir, ConfigFolder, ConfigName)

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() {
	// check that .maker exists
	_, err := os.Stat(ConfigFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(ConfigFolder, 0755)
		if err != nil {
			fmt.Printf("Failed to create config directory %s -- %s\n", ConfigFolder, err)
		}
	}

	// check if config exists to create, or to verify
	_, err = os.Stat(ConfigPath)
	if os.IsNotExist(err) {
		config := &ConfigFile{}
		err = CreateConfigFile(config)
		if err != nil {
			fmt.Printf("Failed to create config file %s -- %s\n", ConfigName, err)
		}
		fmt.Println("Config file generated at", ConfigPath)
	} else {
		ShowCurrentConfig()
	}
}

// ShowCurrentConfig prints out the current config file
func ShowCurrentConfig() {
	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		fmt.Println("Failed to read file -- ", err)
	}
	fmt.Printf("\nCurrent Config:\n\n%s", string(data))

	var confirmation string
	fmt.Print("Is this info still accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		config := &ConfigFile{}
		err = CreateConfigFile(config)
		if err != nil {
			fmt.Printf("Failed to create config file %s -- %s\n", ConfigName, err)
		}
		fmt.Println("Config file generated at", ConfigPath)
	}
}

// CreateConfigFile makes the config file to use in all DO commands
func CreateConfigFile(config *ConfigFile) error {
	// ask for PAT token
	fmt.Println("Please authenticate using your Digital Ocean account...")
	fmt.Println("Tokens can be generated at https://cloud.digitalocean.com/account/api/tokens")
	fmt.Print("Enter PAT Token: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to capture password -- ", err)
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
	return viper.WriteConfigAs(ConfigPath)
}

// LoadConfig parses the viper config file and loads into a struct
func LoadConfig() (string, string) {
	// This all needs to be cleaned up using vars but for now...
	viper.SetConfigFile(ConfigPath)
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	return viper.GetString("pat_token"), viper.GetString("default_region")
}
