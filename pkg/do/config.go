package do

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// ConfigFile makes up the required settings in a DO config file
type ConfigFile struct {
	PatToken      string `json:"pat_token"`
	DefaultRegion string `json:"region"`
}

/* check if config exists in home dir
if doesnt, prompt to se config
if does, ask if want to  reconfig?
reconfig if yes
*/

// Configure sets the PAT token and default Region for Digital Ocean
func Configure() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configFolder := homeDir + "/.maker"
	_, err = os.Stat(configFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(configFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	configFilePath := configFolder + "/do_config"
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		CreateConfigFile()
	} else {
		ShowCurrentConfig(configFilePath)
	}
}

// ShowCurrentConfig prints out the current config file
func ShowCurrentConfig(configFilePath string) {
	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Print(file)

	var confirmation string
	fmt.Print("Is this info still accurate? (Y/n): ")
	fmt.Scanln(&confirmation)
	confirmation = strings.ToLower(string(confirmation))
	println()
	if confirmation != "n" {
		CreateConfigFile()
	}
}

// CreateConfigFile makes the config file to use in all DO commands
func CreateConfigFile() {
	// ask for PAT token
	fmt.Println("Please authenticate using your Digital Ocean account...")
	fmt.Println("Tokens can be generated at https://cloud.digitalocean.com/account/api/tokens")
	fmt.Print("Enter PAT Token: ")
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	os.Setenv("DO_PAT_TOKEN", string(pass))
	println()

	// ask for default region
	var region string
	fmt.Print("Default Region: ")
	fmt.Scanln(&region)
	os.Setenv("DO_DEFAULT_REGION", string(region))
	println()
}
