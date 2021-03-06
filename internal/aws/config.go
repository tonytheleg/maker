package aws

import (
	"fmt"
	"io/ioutil"
	"maker/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// CredsFile makes up the required settings in a AWS config file
type CredsFile struct {
	AccessKeyID     string
	SecretAccessKey string
	DefaultRegion   string
}

// CredsName is the name of the config file used by Maker
var CredsName = "aws_credentials"

// CredsPath is the full path to the ConfigFile
var CredsPath = filepath.Join(utils.ConfigFolderPath, CredsName)

// Configure sets up the aws credentials file needed to auth with AWS
func Configure() error {
	// check that .maker exists
	_, err := os.Stat(utils.ConfigFolderPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(utils.ConfigFolderPath, 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to creds folder %s:", utils.ConfigFolder)
		}
	}

	// check if config exists to create, or to verify
	_, err = os.Stat(CredsPath)
	if os.IsNotExist(err) {
		creds := &CredsFile{}
		err = CreateCredsFile(creds)
		if err != nil {
			return errors.Wrapf(err, "Failed to create creds file %s:", CredsName)
		}
	}
	fmt.Println("Creds file generated at", CredsPath)
	ShowCurrentCreds()
	return nil
}

// ShowCurrentCreds prints out the current credentials file
func ShowCurrentCreds() error {
	data, err := ioutil.ReadFile(CredsPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read file %s:", CredsPath)
	}
	fmt.Printf("\nCurrent Credentials:\n\n%s", string(data))

	var confirmation string
	fmt.Print("Is this info accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		creds := &CredsFile{}
		err = CreateCredsFile(creds)
		if err != nil {
			return errors.Wrapf(err, "Failed to create creds file %s:", CredsName)
		}
		fmt.Println("Credentials file generated at", CredsPath)
	}
	return nil
}

// CreateCredsFile makes the credentials file to use in all AWS commands
func CreateCredsFile(creds *CredsFile) error {
	// ask for access key
	var accessKey string
	fmt.Print("Enter AWS Access Key ID: ")
	fmt.Scanln(&accessKey)
	creds.AccessKeyID = string(accessKey)

	fmt.Print("Enter AWS Secret Key ID: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return errors.Errorf("Failed to capture password:", err)
	}
	creds.SecretAccessKey = string(pass)
	println()

	var region string
	fmt.Print("Enter Default Region: ")
	fmt.Scanln(&region)
	creds.DefaultRegion = string(region)

	viper.SetConfigType("toml")
	viper.Set("default.aws_access_key_id", creds.AccessKeyID)
	viper.Set("default.aws_secret_access_key", creds.SecretAccessKey)
	viper.Set("region.default_region", creds.DefaultRegion)
	return viper.WriteConfigAs(CredsPath)
}

// LoadConfig parses the viper config file and loads into a struct
func LoadConfig() (string, error) {
	// This all needs to be cleaned up using vars but for now...
	viper.SetConfigFile(CredsPath)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		return "", errors.Errorf("Error reading creds file %s:", CredsPath, err)
	}

	return viper.GetString("region.default_region"), nil
}
