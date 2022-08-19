package ui

import (
	"github.com/rivo/tview"

	"github.com/ardnew/bases/ui/view/input"
	"github.com/ardnew/bases/ui/view/output"
)

type UI struct {
	*tview.Application
	root *tview.Flex
	in   *input.View
	out  *output.View
}

func New() *UI {
	// Allocate all primitives before configuring
	return (&UI{
		Application: tview.NewApplication(),
		root:        tview.NewFlex(),
		in:          input.New(),
		out:         output.New(),
	}).init()
}

func (ui *UI) init() *UI {
	ui.root.
		SetDirection(tview.FlexRow).
		AddItem(ui.out, 0, 1, false).
		AddItem(ui.in, 1, 0, true)
	ui.SetRoot(ui.root, true)
	return ui
}

func (ui *UI) Run() error {
	return ui.SetFocus(ui.root).Run()
}
