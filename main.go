// main.go
package main

import (
	"github.com/zvelocity/k8s-tui/ui"
	"log"
)

func main() {
	log.Println("Starting k8s-tui...")

	app := ui.NewApp()

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
