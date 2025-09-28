package main

import (
	"log"

	"github.com/nimbodex/microservices-factory/inventory/internal/app"
)

func main() {
	log.Println("Starting Inventory Service...")

	container := app.NewContainer()
	inventoryApp := container.BuildApp()

	if err := inventoryApp.Run(); err != nil {
		log.Fatalf("Failed to run inventory service: %v", err)
	}
}
