package azure

import (
	"fmt"
	"io/ioutil"
	"maker/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ConfigFile makes up the required settings in a DO config file
type ConfigFile struct {
}

// ConfigName is the name of the config file used by Maker
var ConfigName = "azure_config"

// ConfigPath is the full path to the ConfigFile
var ConfigPath = filepath.Join(utils.HomeDir, utils.ConfigFolder, ConfigName)

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() error {
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
		config := &ConfigFile{}
		err = CreateConfigFile(config)
		if err != nil {
			return errors.Wrapf(err, "Failed to create config file %s:", ConfigName)
		}
	}
	fmt.Println("Config file generated at", ConfigPath)
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
	// DO SAMPLE
	// ask for PAT token
	//fmt.Println("Please authenticate using your Digital Ocean account...")
	//fmt.Println("Tokens can be generated at https://cloud.digitalocean.com/account/api/tokens")
	//fmt.Print("Enter PAT Token: ")
	//pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	//if err != nil {
	//	return errors.Errorf("Failed to capture password:", err)
	//}
	//config.PatToken = string(pass)
	//println()
	//
	//// ask for default region
	//var region string
	//fmt.Print("Default Region: ")
	//fmt.Scanln(&region)
	//config.DefaultRegion = string(region)
	//println()
	//
	//viper.SetConfigType("yaml")
	//viper.Set("pat_token", config.PatToken)
	//viper.Set("default_region", config.DefaultRegion)
	//return viper.WriteConfigAs(ConfigPath)
	return nil
}

// LoadConfig parses the viper config file and loads into a struct
func LoadConfig() (string, string, error) {
	// This all needs to be cleaned up using vars but for now...
	viper.SetConfigFile(ConfigPath)
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		return "", "", errors.Errorf("Error reading config file %s:", ConfigPath, err)
	}

	// DO SAMPLE
	// return viper.GetString("pat_token"), viper.GetString("default_region"), nil
	return "", "", nil
}
