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

	config := GetConfig(fileSystem)

	fmt.Println(config.Dirs)
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
