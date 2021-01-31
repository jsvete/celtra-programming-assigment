package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Panel is a wrapper around widgets.Paragraph and is used to display information about a selected item in a ui.List.
type Panel struct {
	widgets.Paragraph
}

// Update displays the text in the panel and updates the component.
func (p *Panel) Update(text string) {
	p.Text = text
}

// NewPanel creates a new panel where detailed data about a selected item from ui.List can be displayed.
func NewPanel(title string, data string) *Panel {
	panel := &Panel{
		Paragraph: *widgets.NewParagraph(),
	}

	panel.Title = title
	panel.Text = data
	panel.BorderStyle.Fg = termui.ColorBlue

	return panel
}
