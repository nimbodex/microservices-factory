package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nimbodex/microservices-factory/notification/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаем приложение
	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Запускаем приложение
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}

	// Ожидаем сигнал завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down...")

	// Останавливаем приложение
	if err := application.Stop(ctx); err != nil {
		log.Printf("Error stopping app: %v", err)
	}
}
