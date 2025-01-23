package main

import (
	"github.com/zellijsessions/utils"
	"github.com/zellijsessions/zellij-session"
	"os"
	"strings"
)

func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)

	zellijSession := zellijSession.NewZellijSession(fileSystem)

	findArgs := append(zellijSession.Config.Dirs, "-mindepth", "1", "-maxdepth", "1", "-type", "d")

	var findStdIn strings.Reader
	findStdOut := utils.ExecCommand("find", findArgs, findStdIn)

	dirBuilder := utils.ExecCommand("fzf", []string{}, *strings.NewReader(findStdOut.String()))

	chosenDir := strings.TrimSpace(dirBuilder.String())

	zellijSession.StartSession(chosenDir)

}
