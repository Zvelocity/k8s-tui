package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	clientset, err := initK8sClient()
	if err != nil {
		fmt.Printf("Error initializing Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	namespace := "default"
	if len(os.Args) > 1 {
		namespace = os.Args[1]
	}

	model := Model{
		clientset:           clientset,
		namespace:           namespace,
		resourceTypes:       []ResourceType{Pods, Deployments, Services, ConfigMaps, Logs},
		currentResourceType: 0,
		cursor:              0,
		loading:             true,
		mode:                normalMode,
		namespaces:          []string{},
		namespaceCursor:     0,
		kubectlOutput:       "",
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		fmt.Print("\033[H\033[2J")
	}
}
