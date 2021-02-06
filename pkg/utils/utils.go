package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// HandleErr handles error checking and exiting upon failures
func HandleErr(message string, err error) {
	if err != nil {
		fmt.Println(message, err)
		os.Exit(1)
	}
}

// HomeDir stores the path of the current users Home directory
var HomeDir, _ = os.UserHomeDir()

// ConfigFolder is the name of Makers config folder stored in Home
var ConfigFolder = ".maker"

// ConfigFolderPath is the full path to the .maker directory
var ConfigFolderPath = filepath.Join(HomeDir, ConfigFolder)
