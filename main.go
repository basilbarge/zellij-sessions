package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
)


func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)
	var subdirs []fs.DirEntry

	subdirs, err := fs.ReadDir(fileSystem, "Documents")

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range subdirs {
		fmt.Println(path)
		command := exec.Command("tmux", "new-session", "-A", "-s", path.Name())
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err := command.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}

