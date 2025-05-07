// ui/podlist.go
package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/zvelocity/k8s-tui/k8s"
	corev1 "k8s.io/api/core/v1"
)

// PodListView represents the pod list view
type PodListView struct {
	table *tview.Table
	app   *App
	pods  []corev1.Pod
}

// NewPodListView creates a new pod list view
func NewPodListView(app *App) *PodListView {
	view := &PodListView{
		table: tview.NewTable().SetBorders(true).SetFixed(1, 0),
		app:   app,
	}

	// Set up the table
	view.table.SetSelectable(true, false)
	view.setupHeaders()

	return view
}

// setupHeaders sets up the table headers
func (v *PodListView) setupHeaders() {
	headers := []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE"}

	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false)
		v.table.SetCell(0, col, cell)
	}
}

// Refresh refreshes the pod list
func (v *PodListView) Refresh() {
	if v.app.k8sClient == nil {
		return
	}

	ctx := context.Background()
	pods, err := v.app.k8sClient.GetPods(ctx)
	if err != nil {
		// Handle error (you might want to show this in the UI)
		return
	}

	v.pods = pods
	v.updateTable()
}

// updateTable updates the table with current pod data
func (v *PodListView) updateTable() {
	// Clear existing rows (except header)
	for row := 1; row < v.table.GetRowCount(); row++ {
		v.table.RemoveRow(row)
	}

	// Add pod data
	for row, pod := range v.pods {
		actualRow := row + 1

		// Name
		v.table.SetCell(actualRow, 0, tview.NewTableCell(pod.Name))

		// Ready
		ready := fmt.Sprintf("%d/%d", getReadyContainers(pod), len(pod.Spec.Containers))
		v.table.SetCell(actualRow, 1, tview.NewTableCell(ready))

		// Status
		status := k8s.GetPodStatus(pod)
		v.table.SetCell(actualRow, 2, tview.NewTableCell(status))

		// Restarts
		restarts := k8s.GetPodRestarts(pod)
		v.table.SetCell(actualRow, 3, tview.NewTableCell(fmt.Sprintf("%d", restarts)))

		// Age
		age := time.Since(pod.CreationTimestamp.Time).Round(time.Second)
		v.table.SetCell(actualRow, 4, tview.NewTableCell(formatAge(age)))
	}
}

// getReadyContainers returns the number of ready containers in a pod
func getReadyContainers(pod corev1.Pod) int {
	ready := 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			ready++
		}
	}
	return ready
}

// formatAge formats a duration into a human-readable age string
func formatAge(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// GetTable returns the underlying table
func (v *PodListView) GetTable() *tview.Table {
	return v.table
}
