package models


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
