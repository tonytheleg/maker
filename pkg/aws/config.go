package aws

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

// CredsFile makes up the required settings in a AWS config file
type CredsFile struct {
	AccessKeyId     string
	SecretAccessKey string
	DefaultRegion   string
}

var HomeDir, _ = os.UserHomeDir()
var CredsFolder = ".maker"
var CredsName = "aws_credentials"
var CredsPath = filepath.Join(HomeDir, CredsFolder, CredsName)

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() {
	// check that .maker exists
	_, err := os.Stat(CredsFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(CredsFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// check if config exists to create, or to verify
	_, err = os.Stat(CredsPath)
	if os.IsNotExist(err) {
		creds := &CredsFile{}
		err = CreateCredsFile(creds)
		if err != nil {
			panic(err)
		}
		fmt.Println("Creds file generated at", CredsPath)
	} else {
		ShowCurrentCreds()
	}
}

// ShowCurrentCreds prints out the current config file
func ShowCurrentCreds() {
	data, err := ioutil.ReadFile(CredsPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nCurrent Credentials:\n\n%s", string(data))

	var confirmation string
	fmt.Print("Is this info still accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()

	if confirmation != "y" {
		creds := &CredsFile{}
		err = CreateCredsFile(creds)
		if err != nil {
			panic(err)
		}
		fmt.Println("Credentials file generated at", CredsPath)
	}
}

// CreateConfigFile makes the config file to use in all DO commands
func CreateCredsFile(creds *CredsFile) error {
	// ask for access key
	var accessKey string
	fmt.Print("Enter AWS Access Key ID: ")
	fmt.Scanln(&accessKey)
	creds.AccessKeyId = string(accessKey)

	fmt.Print("Enter AWS Secret Key ID: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	creds.SecretAccessKey = string(pass)
	println()

	var region string
	fmt.Print("Enter Default Region: ")
	fmt.Scanln(&region)
	creds.DefaultRegion = string(region)

	viper.SetConfigType("toml")
	viper.Set("default.aws_access_key_id", creds.AccessKeyId)
	viper.Set("default.aws_secret_access_key", creds.SecretAccessKey)
	viper.Set("region.default_region", creds.DefaultRegion)
	return viper.WriteConfigAs(CredsPath)
}

// LoadConfig parses the viper config file and loads into a struct
func LoadConfig() string {
	// This all needs to be cleaned up using vars but for now...
	viper.SetConfigFile(CredsPath)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	return viper.GetString("region.default_region")
}
