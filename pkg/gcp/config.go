package gcp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ConfigFile makes up the required settings in a DO config file
type ConfigFile struct {
	DefaultRegion string
	GcpProject    string
	Keyfile       string
}

// HomeDir stores the path of the current users Home directory
var HomeDir, _ = os.UserHomeDir()

// ConfigFolder is the name of Makers config folder stored in Home
var ConfigFolder = ".maker"

// ConfigName is the name of the config file used by Maker
var ConfigName = "gcp_config"

// ConfigPath is the full path to the ConfigFile
var ConfigPath = filepath.Join(HomeDir, ConfigFolder, ConfigName)

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() error {
	// check that .maker exists
	_, err := os.Stat(ConfigFolder)
	if os.IsExist(err) {
		err := os.Mkdir(filepath.Join(HomeDir, ConfigFolder), 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config directory %s:", ConfigFolder)
		}
	}

	// check if config exists to create, or to verify
	_, err = os.Stat(ConfigPath)
	if os.IsNotExist(err) {
		config := &ConfigFile{}
		err = CreateConfigFile(config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
		fmt.Println("Config file generated at", ConfigPath)
	}
	ShowCurrentConfig()
	return nil
}

// ShowCurrentConfig prints out the current config file
func ShowCurrentConfig() error {
	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s:", ConfigPath)
	}
	fmt.Printf("\nCurrent Config:\n\n%s", string(data))

	var confirmation string
	fmt.Printf("\nIs this info accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		config := &ConfigFile{}
		err = CreateConfigFile(config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
		fmt.Println("Config file generated at", ConfigPath)
	}
	return nil
}

// CreateConfigFile makes the config file to use in all DO commands
func CreateConfigFile(config *ConfigFile) error {
	// ask for default region
	var region string
	fmt.Print("Enter Default Compute Zone (ie: us-east1-b): ")
	fmt.Scanln(&region)
	config.DefaultRegion = string(region)
	println()

	// ask for default GCP Project
	var project string
	fmt.Print("Enter Target GCP Project: ")
	fmt.Scanln(&project)
	config.GcpProject = string(project)
	println()

	// ask for path to key json
	var keyfile string
	fmt.Println("GCP requires a Service Account Key file to make requests")
	fmt.Println("See 'https://cloud.google.com/iam/docs/creating-managing-service-account-keys#iam-service-account-keys-create-console' for help")
	fmt.Printf("\nEnter path to your key file (full path): ")
	fmt.Scanln(&keyfile)
	config.Keyfile = string(keyfile)
	println()

	viper.SetConfigType("yaml")
	viper.Set("keyfile", config.Keyfile)
	viper.Set("default_region", config.DefaultRegion)
	viper.Set("gcp_project", config.GcpProject)
	return viper.WriteConfigAs(ConfigPath)
}

// LoadConfig parses the viper config file and loads into a struct
func LoadConfig() (string, string, string, error) {
	viper.SetConfigFile(ConfigPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return "", "", "", errors.Errorf("Error reading config file %s:", ConfigPath, err)
	}
	return viper.GetString("keyfile"), viper.GetString("default_region"), viper.GetString("gcp_project"), nil
}
