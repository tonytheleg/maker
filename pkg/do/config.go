package do

import (
	"fmt"
	"io/ioutil"
	"maker/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// ConfigFile makes up the required settings in a DO config file
type ConfigFile struct {
	PatToken              string
	DefaultRegion         string
	SpacesAccessKey       string
	SpacesSecretKey       string
	SpacesDefaultEndpoint string
}

// ConfigName is the name of the config file used by Maker
var ConfigName = "do_config"

// ConfigPath is the full path to the ConfigFile
var ConfigPath = filepath.Join(utils.HomeDir, utils.ConfigFolder, ConfigName)

// SetupConfig setups config directory and file
func SetupConfig() error {
	// check that .maker exists
	_, err := os.Stat(utils.ConfigFolder)
	if os.IsExist(err) {
		err := os.Mkdir(filepath.Join(utils.HomeDir, utils.ConfigFolder), 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config directory %s:", utils.ConfigFolder)
		}
	}

	// check if config exists to create, or to verify
	_, err = os.Stat(ConfigPath)
	if os.IsNotExist(err) {
		Configure()
	}
	ConfirmCurrentConfig()
	return nil
}

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() error {
	task := GetConfigTasks()
	config, _ := LoadConfig()

	switch task {
	case "1":
		err := CreateDropletConfigFile(&config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
	case "2":
		err := CreateSpacesConfigFile(&config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
	default:
		err := CreateDropletConfigFile(&config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
		err = CreateSpacesConfigFile(&config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
	}
	fmt.Println("Config file generated at", ConfigPath)
	PrintCurrentConfig()
	return nil
}

// ConfirmCurrentConfig prints out the current config file
func ConfirmCurrentConfig() error {
	PrintCurrentConfig()

	// show config and confirm
	var confirmation string
	fmt.Printf("\nIs this info accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		err := Configure()
		if err != nil {
			return errors.Wrap(err, "Failed to configure")
		}
	}
	return nil
}

// GetConfigTasks determines what configs to set
func GetConfigTasks() string {
	var selection string
	fmt.Println("Select Configuration Option:")
	fmt.Printf("1. Set Droplet Config\n2. Set Spaces Config\n3. Set All\n")
	fmt.Printf("Selection?: ")
	fmt.Scanln(&selection)
	return string(selection)
}

// PrintCurrentConfig outputs the current config file
func PrintCurrentConfig() error {
	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s:", ConfigPath)
	}
	fmt.Printf("\nCurrent Config:\n\n%s", string(data))
	return nil
}

// CreateDropletConfigFile creates the config file to use in all DO commands
func CreateDropletConfigFile(config *ConfigFile) error {
	// ask for PAT token
	fmt.Println("Please authenticate using your Digital Ocean account...")
	fmt.Println("Tokens can be generated at https://cloud.digitalocean.com/account/api/tokens")
	fmt.Print("Enter PAT Token: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Errorf("Failed to capture password:", err)
	}
	config.PatToken = string(pass)
	println()

	// ask for default region
	var region string
	fmt.Print("Default Region: ")
	fmt.Scanln(&region)
	config.DefaultRegion = string(region)
	println()

	// create the config
	viper.SetConfigType("yaml")
	viper.Set("pat_token", config.PatToken)
	viper.Set("default_region", config.DefaultRegion)
	return viper.WriteConfigAs(ConfigPath)
}

// CreateSpacesConfigFile creates the config file to use in spaces creation
func CreateSpacesConfigFile(config *ConfigFile) error {
	// ask for access key
	var accessKey string
	fmt.Println("Please authenticate using your Digital Ocean account...")
	fmt.Println("Tokens can be generated at https://cloud.digitalocean.com/account/api/tokens")
	fmt.Print("Enter Access Key: ")
	fmt.Scanln(&accessKey)
	config.SpacesAccessKey = string(accessKey)
	println()

	// ask for secret key
	fmt.Print("Enter Secret Key: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Errorf("Failed to secret key:", err)
	}
	config.SpacesSecretKey = string(pass)
	println()

	// ask for default region
	var region string
	fmt.Print("Default Endpoint Region (ie, nyc3): ")
	fmt.Scanln(&region)
	config.SpacesDefaultEndpoint = string(region)
	println()

	// create the config
	viper.SetConfigType("yaml")
	viper.Set("spaces_access_key", config.SpacesAccessKey)
	viper.Set("spaces_secret_key", config.SpacesSecretKey)
	viper.Set("spaces_endpoint_region", config.SpacesDefaultEndpoint)
	return viper.WriteConfigAs(ConfigPath)
}

// LoadConfig parses the viper config file and loads into a struct
func LoadConfig() (ConfigFile, error) {
	conf := &ConfigFile{}
	viper.SetConfigFile(ConfigPath)
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		return ConfigFile{}, errors.Errorf("Error reading config file %s:", ConfigPath, err)
	}
	err = viper.Unmarshal(conf)
	if err != nil {
		return ConfigFile{}, errors.Errorf("Error reading config file %s:", ConfigPath, err)
	}
	return *conf, nil
}
