package main

import (
	"fmt"

	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/zellijsessions/utils"
	"github.com/zellijsessions/zellij-session"
)

//Integrate bubbletea
// need model with View, Init, Update functions

type App struct {
	dirList list.Model
}

func NewApp() App {
	return App{dirList: list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 60)}
}

func (a App) View() string {
	return a.dirList.View()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return a, tea.Quit
		}
	}

	a.dirList, cmd = a.dirList.Update(msg)

	return a, cmd
}

func (a App) Init() tea.Cmd {
	return nil
}

type DirListItem struct {
	title       string
	description string
}

func (d DirListItem) FilterValue() string {
	return d.title
}

func (d DirListItem) Title() string {
	return d.title
}

func (d DirListItem) Description() string {
	return d.description
}

func NewDirListItem(title, description string) DirListItem {
	return DirListItem{title: title, description: description}
}

func main() {
	root := "/home/basilbarge"
	fileSystem := os.DirFS(root)
	zellijSession := zellijSession.NewZellijSession(fileSystem)

	var findDirs []string

	for _, dir := range zellijSession.ProjectDirs {
		findDirs = append(findDirs, dir.AbsPath)
	}

	listItems := []list.Item{}

	for _, dirString := range findDirs {
		listItems = append(listItems, NewDirListItem(dirString, ""))
	}

	app := NewApp()

	app.dirList.SetItems(listItems)

	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		utils.LogError(fmt.Sprintf("There was an error running the app, %v", err))
	}

	//dirBuilder := utils.ExecCommand("fzf", []string{}, *strings.NewReader(findStdOut.String()))

	//chosenDir := strings.TrimSpace(dirBuilder.String())

	//zellijSession.StartSession(chosenDir)

}
