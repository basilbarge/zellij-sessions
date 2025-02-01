package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	DirList     list.Model
	SelectedDir string
}

func NewApp() App {
	return App{DirList: list.New([]list.Item{}, list.NewDefaultDelegate(), 40, 60)}
}

func (a App) View() string {
	return a.DirList.View()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return a, tea.Quit

		case "enter":
			//if we press enter to apply a filter, ignore the press
			if a.DirList.FilterState() == list.Filtering {
				break
			}

			//otherwise want to store the selected list item
			//in the application state
			i, ok := a.DirList.SelectedItem().(DirListItem)

			if ok {
				//storing description because this is
				//where I save the full directory path
				a.SelectedDir = i.Description()
			}

			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.DirList.SetHeight(msg.Height)
		a.DirList.SetWidth(msg.Width)
	}

	var cmd tea.Cmd
	a.DirList, cmd = a.DirList.Update(msg)

	return a, cmd
}

func (a App) Init() tea.Cmd {
	return nil
}
