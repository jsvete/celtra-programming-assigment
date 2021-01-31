package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// List is a wrapper around widgets.List and is used to display all items.
type List struct {
	widgets.List
}

// Update updates the rows in the list with new values and re-renders it
func (l *List) Update(data ...string) {
	if len(data) > 0 {
		l.Rows = append(l.Rows, data...)
	}
}

// NewList creates a new list with initial data and set coordinates
func NewList(title string, data []string) *List {
	list := &List{
		List: *widgets.NewList(),
	}

	list.Title = title
	list.Rows = data
	list.TextStyle = termui.NewStyle(termui.ColorYellow)
	list.WrapText = false

	return list
}
