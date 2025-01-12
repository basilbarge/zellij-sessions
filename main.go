package main

import (
	"encoding/json"
	"fmt"
	"github.com/zellijsessions/utils"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Config struct {
	Dirs []string `json:"dirs"`
}

func newConfig(fileSystem fs.FS, configPath string) *Config {
	var paths *Config

	data, err := fs.ReadFile(fileSystem, configPath)

	if err != nil {
		fmt.Printf("There was an error reading the configuration file. %s\n", err)
	}

	err = json.Unmarshal(data, &paths)

	if err != nil {
		fmt.Printf("There was an error unmarshalling json. %s\n", err)
	}

	return paths
}

func (config *Config) RemoveDir(filesystem fs.FS, pathToRemove string) {
	if !slices.Contains(config.Dirs, pathToRemove) {
		fmt.Println(fmt.Errorf("The current configuration does not contain %s as a directory so it cannot be removed", pathToRemove))
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
		fmt.Printf("There was an error marshaling new config to json. %s\n", err)
	}

	err = os.WriteFile("/home/basilbarge/Documents/Projects/tmux-sessions/config.json", marshaledConfig, 0770)

	if err != nil {
		fmt.Printf("There was an error writing the new config. %s\n", err)
	}
}

func (config *Config) AddDir(filesystem fs.FS, pathToAdd string) {
	if _, err := os.Stat(pathToAdd); err != nil {

		if os.IsNotExist(err) {
			fmt.Println("The directory you want to add does not exist on your machine. Check your spelling or try a different one!")
		} else {
			fmt.Printf("An error occured when searching for the directory to add. %s\n", err)
		}
	}

	config.Dirs = append(config.Dirs, pathToAdd)

	marshaledConfig, err := json.MarshalIndent(config, "", "	")

	if err != nil {
		fmt.Printf("There was an error marshaling new config to json. %s\n", err)
	}

	err = os.WriteFile("/home/basilbarge/Documents/Projects/tmux-sessions/config.json", marshaledConfig, 0770)

	if err != nil {
		fmt.Printf("There was an error writing the new config. %s\n", err)
	}
}

type ZellijSession struct {
	Config     *Config
	Filesystem fs.FS
}

func newZellijSession(filesystem fs.FS) *ZellijSession {
	return &ZellijSession{
		Config:     newConfig(filesystem, "Documents/Projects/zellij-sessions/config.json"),
		Filesystem: filesystem,
	}
}

func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)

	zellijSession := newZellijSession(fileSystem)

	findArgs := append(zellijSession.Config.Dirs, "-type", "d", "-maxdepth", "1")

	var findStdIn strings.Reader
	findStdOut := utils.ExecCommand("find", findArgs, findStdIn)

	dirBuilder := utils.ExecCommand("fzf", []string{}, *strings.NewReader(findStdOut.String()))

	chosenDir := strings.TrimSpace(dirBuilder.String())

	utils.RunShellCommand("cd", []string{chosenDir})

	utils.RunShellCommand("zellij", []string{"--session", filepath.Base(chosenDir)})
}
