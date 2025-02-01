package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/zellijsessions/models"
	"github.com/zellijsessions/utils"
	"github.com/zellijsessions/zellij-session"
	"os"
)

func main() {
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

	zellijSession.StartSession(selectedDir)

}
