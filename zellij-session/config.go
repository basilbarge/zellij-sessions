package zellijSession

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/zellijsessions/utils"
)

type Config struct {
	Dirs []string `json:"dirs"`
}

func NewConfig(fileSystem fs.FS, configPath string) *Config {
	var paths *Config

	data, err := fs.ReadFile(fileSystem, configPath)

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error reading the configuration file. %s\n", err))
	}

	err = json.Unmarshal(data, &paths)

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error unmarshalling json. %s\n", err))
	}

	return paths
}

func (config *Config) RemoveDir(pathToRemove string) {
	if (strings.Contains(pathToRemove, "~")) {
		pathToRemove = parseTildeInPath(pathToRemove)
	}


	if !slices.Contains(config.Dirs, pathToRemove) {
		utils.LogError(fmt.Sprintln(fmt.Errorf("The current configuration does not contain %s as a directory so it cannot be removed", pathToRemove)))
		return
	}

	idxToRemove := 0
	for idx, path := range config.Dirs {
		if path == pathToRemove {
			idxToRemove = idx
			break
		}
	}

	config.Dirs = append(config.Dirs[:idxToRemove], config.Dirs[idxToRemove+1:]...)

	marshaledConfig, err := json.MarshalIndent(config, "", "	")

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error marshaling new config to json. %s\n", err))
	}

	err = os.WriteFile("/home/basilbarge/Documents/Projects/zellij-sessions/config.json", marshaledConfig, 0770)

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error writing the new config. %s\n", err))
	}
}

func (config *Config) AddDir(pathToAdd string) {
	if (strings.Contains(pathToAdd, "~")) {
		pathToAdd = parseTildeInPath(pathToAdd)
	}

	if _, err := os.Stat(pathToAdd); err != nil {

		if os.IsNotExist(err) {
			utils.LogError("The directory you want to add does not exist on your machine. Check your spelling or try a different one!")
		} else {
			utils.LogError(fmt.Sprintf("An error occured when searching for the directory to add. %s\n", err))
		}
	}

	config.Dirs = append(config.Dirs, pathToAdd)

	marshaledConfig, err := json.MarshalIndent(config, "", "	")

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error marshaling new config to json. %s\n", err))
	}

	err = os.WriteFile("/home/basilbarge/Documents/Projects/zellij-sessions/config.json", marshaledConfig, 0770)

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error writing the new config. %s\n", err))
	}
}

func parseTildeInPath(path string) string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error getting the user's home directory: %v", err))
	}

	return strings.Replace(path, "~", homeDir, 1)
}
