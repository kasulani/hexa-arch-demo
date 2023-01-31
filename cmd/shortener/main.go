package main

import (
	"log"
	"os"

	"github.com/kasulani/hexa-arch-demo/internal/app"
)

func main() {
	container := app.Container()

	defer container.Cleanup()

	if err := container.Invoke(app.Run); err != nil {
		log.Printf("failed to start application: %q\n", err)
		os.Exit(1)
	}

	if err := container.Invoke(app.TerminateConnections);  err != nil {
		log.Printf("failed to terminate connections: %q\n", err)
		os.Exit(1)
	}
}
