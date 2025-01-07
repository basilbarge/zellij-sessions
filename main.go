package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	//"log"
	"os"
	//"os/exec"
)

type Config struct {
	Dirs []string `json:"dirs"`
}

func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)

	AddDir(fileSystem, "testDir")

	config := GetConfig(fileSystem)

	fmt.Println(config.Dirs)
}

func AddDir(filesystem fs.FS, pathToAdd string) {
	config := GetConfig(filesystem)

	fmt.Println(config.Dirs)

	config.Dirs = append(config.Dirs, pathToAdd)

	fmt.Println(config.Dirs)

	marshaledConfig, err := json.Marshal(config)

	if err != nil {
		fmt.Printf("There was an error marshaling new config to json. %s\n", err)
	}

	err = os.WriteFile("/home/basilbarge/Documents/Projects/tmux-sessions/config.json", marshaledConfig, 0770)

	if err != nil {
		fmt.Printf("There was an error writing the new config. %s\n", err)
	}
}

func GetConfig(fileSystem fs.FS) Config {
	var paths Config

	data, err := fs.ReadFile(fileSystem, "Documents/Projects/tmux-sessions/config.json")

	if err != nil {
		fmt.Printf("There was an error reading the configuration file. %s\n", err)
	}

	err = json.Unmarshal(data, &paths)

	if err != nil {
		fmt.Printf("There was an error unmarshalling json. %s\n", err)
	}

	return paths
}
