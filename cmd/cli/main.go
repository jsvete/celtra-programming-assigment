package main

import (
	"celtra-programming-assigment/cmd/cli/ui"
	"fmt"
	"log"
	"time"

	"github.com/gizak/termui/v3"
)

var (
	ids = []string{
		"ID1",
		"ID2",
		"ID3",
		"ID4",
		"ID5",
		"ID6",
		"ID7",
		"ID8",
		"ID9",
		"ID10",
	}
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize CLI client: %v", err)
	}
	defer termui.Close()

	// setup UI elements
	l := ui.NewList("ACCOUNT ID", ids)
	p := ui.NewPanel("DATA", "")

	l.SetRect(0, 0, 25, 20)
	p.SetRect(26, 0, 70, 20)

	termui.Render(l, p)

	previousKey := ""
	uiEvents := termui.PollEvents()
	for {
		e := <-uiEvents

		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		case "<Enter>":
			p.Update(fmt.Sprintf("<%s>: %s%s%s\n", time.Now().UTC().Format("15:04:05.000 UTC"), "DATA FOR ", ids[l.SelectedRow], " HERE"))
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		termui.Render(l, p)
	}
}
