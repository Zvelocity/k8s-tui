// ui/app.go
package ui

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zvelocity/k8s-tui/k8s"
)

// App represents our TUI application
type App struct {
	app         *tview.Application
	pages       *tview.Pages
	mainFlex    *tview.Flex
	k8sClient   *k8s.Client
	podListView *PodListView
}

// NewApp creates a new TUI application
func NewApp() *App {
	return &App{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
	}
}

// Initialize sets up the initial UI layout
func (a *App) Initialize() {
	// Create header
	header := tview.NewTextView().
		SetText("Kubernetes TUI").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.PrimaryTextColor)

	// Create footer with help text
	footer := tview.NewTextView().
		SetText("Press 'q' to quit | Use arrow keys to navigate").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.SecondaryTextColor)

	// Create main content area
	content := tview.NewTextView().
		SetText("Connecting to Kubernetes cluster...").
		SetTextAlign(tview.AlignCenter)

	// Create the main layout
	a.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(content, 0, 1, true).
		AddItem(footer, 1, 0, false)

	a.pages.AddPage("main", a.mainFlex, true, true)

	// Set up key bindings
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			a.app.Stop()
		}
		return event
	})

	// Initialize Kubernetes client
	go func() {
		client, err := k8s.NewClient()
		if err != nil {
			a.app.QueueUpdateDraw(func() {
				content.SetText(fmt.Sprintf("Failed to connect to Kubernetes: %v", err))
			})
			return
		}

		// Test the connection
		ctx := context.Background()
		if err := client.TestConnection(ctx); err != nil {
			a.app.QueueUpdateDraw(func() {
				content.SetText(fmt.Sprintf("Failed to test connection: %v", err))
			})
			return
		}

		a.k8sClient = client

		// Set up pod list view
		a.app.QueueUpdateDraw(func() {
			a.podListView = NewPodListView(a)

			// Update the main flex layout
			a.mainFlex.Clear()
			a.mainFlex.
				AddItem(header, 1, 0, false).
				AddItem(a.podListView.GetTable(), 0, 1, true).
				AddItem(footer, 1, 0, false)

			// Start refreshing the pod list
			a.podListView.Refresh()
		})
	}()
}

// Run starts the application
func (a *App) Run() error {
	a.Initialize()
	return a.app.SetRoot(a.pages, true).EnableMouse(true).Run()
}
