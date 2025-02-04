package main

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/zellijsessions/models"
	"github.com/zellijsessions/utils"
	"github.com/zellijsessions/zellij-session"
)

type ConfigArgs struct {
	Add    string
	Remove string
}

func main() {
	var args struct {
		Config *ConfigArgs `arg:"subcommand:config"`
	}

	parser := arg.MustParse(&args)

	switch {
	case args.Config != nil:
		if *args.Config == (ConfigArgs{}) {
			parser.FailSubcommand("expected add or remove", "config")
		}

		if args.Config.Add == "" {
			fmt.Println("No add")
		}

		if args.Config.Remove == "" {
			fmt.Println("No remove")
		}
	}

	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)
	zellijSession := zellijSession.NewZellijSession(fileSystem)

	listItems := []list.Item{}
	for _, dir := range zellijSession.ProjectDirs {
		listItems = append(listItems, models.NewDirListItem(dir.Info.Name(), dir.AbsPath))
	}

	app := models.NewApp()

	app.DirList.SetItems(listItems)

	p := tea.NewProgram(app)

	finalApp, err := p.Run()

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error running the app, %v", err))
	}

	selectedDir := finalApp.(models.App).SelectedDir

	if selectedDir == "" {
		utils.LogError("You must select a directory")
	}

	zellijSession.StartSession(selectedDir)

}
