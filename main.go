package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type model struct {
	pods    []string
	message string
}

func (m model) Init() tea.Cmd {
	return fetchPods
}

func fetchPods() tea.Msg {
	kubeconfig := os.Getenv("KUBECONFIG") // or use the default path if you dont have env variable set
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config) // create a new clientset
	if err != nil {
		return err
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{}) // get all pods in the cluster
	if err != nil {
		return err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case []string:
		m.pods = msg
	case error:
		m.message = fmt.Sprintf("Error: %v", msg)
	}
	return m, nil
}

func (m model) View() string { // view function to render the UI
	view := "Hello, Kubernetes!\n\nPods:\n"
	for _, pod := range m.pods { // loop through the pods and display them
		view += fmt.Sprintf("- %s\n", pod)
	}
	view += "\nPress Ctrl+C or 'q' to quit.\n"
	return view
}

func main() {
	p := tea.NewProgram(model{message: "Welcome to your TUI!"})
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
		os.Exit(1)
	}
}
