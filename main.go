package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"unsafe"
)

type Config struct {
	Dirs []string `json:"dirs"`
}

func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)
	config := GetConfig(fileSystem)
	
	findArgs := config.Dirs
	findArgs = append(findArgs, "-type", "d", "-maxdepth", "1")

	find := exec.Command("find", findArgs...)

	var findResult strings.Builder
	find.Stdout = &findResult

	if err := find.Run(); err != nil {
		fmt.Printf("There was a problem running command %s. %s\n", find.Path, err)
	}

	fzf := exec.Command("fzf")

	fzf.Stdin = strings.NewReader(findResult.String())

	var dirBuilder strings.Builder
	fzf.Stdout = &dirBuilder

	if err := fzf.Run(); err != nil {
		fmt.Printf("There was a problem running command %s. %s\n", fzf.Path, err)
	}

	chosenDir := strings.TrimSpace(dirBuilder.String())

	if err := os.Chdir(chosenDir); err != nil {
		fmt.Printf("Could not go to directory %s. %s\n", chosenDir, err)
	}

	zellijCmdString := fmt.Sprintf("zellij --session %s\n", filepath.Base(chosenDir))

    zellijCmdBytes, err := syscall.ByteSliceFromString(zellijCmdString)
    if err != nil {
        log.Fatalln(err)
    }

    var eno syscall.Errno
    for _, c := range zellijCmdBytes {
        _, _, eno = syscall.Syscall(syscall.SYS_IOCTL,
			0,
            syscall.TIOCSTI,
            uintptr(unsafe.Pointer(&c)),
        )
        if eno != 0 {
            log.Fatalln(eno)
        }
    }
}

func RemoveDir(filesystem fs.FS, pathToRemove string) {
	config := GetConfig(filesystem)

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

	fmt.Println(config.Dirs)

	marshaledConfig, err := json.MarshalIndent(config, "", "	")

	if err != nil {
		fmt.Printf("There was an error marshaling new config to json. %s\n", err)
	}

	err = os.WriteFile("/home/basilbarge/Documents/Projects/tmux-sessions/config.json", marshaledConfig, 0770)

	if err != nil {
		fmt.Printf("There was an error writing the new config. %s\n", err)
	}
}

func AddDir(filesystem fs.FS, pathToAdd string) {
	if _, err := os.Stat(pathToAdd); err != nil {

		if os.IsNotExist(err) {
			fmt.Println("The directory you want to add does not exist on your machine. Check your spelling or try a different one!")
		} else {
			fmt.Printf("An error occured when searching for the directory to add. %s\n", err)
		}
	}

	config := GetConfig(filesystem)

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
