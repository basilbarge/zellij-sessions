package main

import (
	"os"
	"strings"

	"github.com/zellijsessions/utils"
	"github.com/zellijsessions/zellij-session"
)

func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)
	zellijSession := zellijSession.NewZellijSession(fileSystem)


	var findDirs []string

	for _, dir := range zellijSession.ProjectDirs {
		findDirs = append(findDirs, dir.AbsPath)
	}

	dirBuilder := utils.ExecCommand("fzf", []string{}, *strings.NewReader(strings.Join(findDirs, "\n")))

	chosenDir := strings.TrimSpace(dirBuilder.String())

	zellijSession.StartSession(chosenDir)

}
