package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/zellijsessions/utils"
	"github.com/zellijsessions/zellij-session"
	"os"
)

type App struct {
	dirList     list.Model
	selectedDir string
}

func NewApp() App {
	return App{dirList: list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 60)}
}

func (a App) View() string {
	return a.dirList.View()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return a, tea.Quit

		case "enter":
			//if we press enter to apply a filter, ignore the press
			if a.dirList.FilterState() == list.Filtering {
				break
			}

			//otherwise want to store the selected list item
			//in the application state
			i, ok := a.dirList.SelectedItem().(DirListItem)

			if ok {
				//storing description because this is
				//where I save the full directory path
				a.selectedDir = i.Description()
			}

			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.dirList.SetHeight(msg.Height)
		a.dirList.SetWidth(msg.Width)
	}

	var cmd tea.Cmd
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

	listItems := []list.Item{}
	for _, dir := range zellijSession.ProjectDirs {
		listItems = append(listItems, NewDirListItem(dir.Info.Name(), dir.AbsPath))
	}

	app := NewApp()

	app.dirList.SetItems(listItems)

	p := tea.NewProgram(app)

	finalApp, err := p.Run()

	if err != nil {
		utils.LogError(fmt.Sprintf("There was an error running the app, %v", err))
	}

	selectedDir := finalApp.(App).selectedDir

	zellijSession.StartSession(selectedDir)

}
