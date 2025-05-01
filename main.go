package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create a simple text view to start
	textView := tview.NewTextView().
		SetText("K8s TUI!\n").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Create a flex layout to hold our components
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(textView, 0, 1, true)

	// Run the application
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
