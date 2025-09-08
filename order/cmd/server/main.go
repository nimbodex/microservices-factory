package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nimbodex/microservices-factory/order/internal/service"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

const (
	port              = ":8080"
	readHeaderTimeout = 30 * time.Second
)

func main() {
	log.Println("Starting Order Service...")

	orderService, err := service.NewOrderService()
	if err != nil {
		log.Fatalf("Failed to create order service: %v", err)
	}

	server, err := orderv1.NewServer(orderService)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	httpServer := &http.Server{
		Addr:              port,
		Handler:           server,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down Order Service...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Order Service listening on %s", port)
	log.Println("Available endpoints:")
	log.Println("\t - POST /api/v1/orders: create order")
	log.Println("\t - GET /api/v1/orders/{uuid}: get order")
	log.Println("\t - POST /api/v1/orders/{uuid}/pay: pay order")
	log.Println("\t - POST /api/v1/orders/{uuid}/cancel: cancel order")

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Order Service stopped")
}
