package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/nimbodex/microservices-factory/order/internal/api/order/v1"
	"github.com/nimbodex/microservices-factory/order/internal/client/grpc"
	orderrepo "github.com/nimbodex/microservices-factory/order/internal/repository/order"
	orderservice "github.com/nimbodex/microservices-factory/order/internal/service/order"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

const (
	port              = ":8080"
	readHeaderTimeout = 30 * time.Second
)

func main() {
	log.Println("Starting Order Service...")

	orderRepo := orderrepo.NewMemoryOrderRepository()

	inventoryClient, err := grpc.NewGRPCInventoryClient()
	if err != nil {
		log.Fatalf("Failed to create inventory client: %v", err)
	}

	paymentClient, err := grpc.NewGRPCPaymentClient()
	if err != nil {
		log.Fatalf("Failed to create payment client: %v", err)
	}

	orderService := orderservice.NewOrderService(orderRepo, inventoryClient, paymentClient)

	apiHandler := v1.NewAPIHandler(orderService)

	server, err := orderv1.NewServer(apiHandler)
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

		if inventoryClient != nil {
			if closeErr := inventoryClient.Close(); closeErr != nil {
				log.Printf("Failed to close inventory client: %v", closeErr)
			}
		}
		if paymentClient != nil {
			if closeErr := paymentClient.Close(); closeErr != nil {
				log.Printf("Failed to close payment client: %v", closeErr)
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if shutdownErr := httpServer.Shutdown(ctx); shutdownErr != nil {
			log.Printf("Server shutdown error: %v", shutdownErr)
		}
	}()

	log.Printf("Order Service listening on %s", port)
	log.Println("Available endpoints:")
	log.Println("\t - POST /api/v1/orders: create order")
	log.Println("\t - GET /api/v1/orders/{uuid}: get order")
	log.Println("\t - POST /api/v1/orders/{uuid}/pay: pay order")
	log.Println("\t - POST /api/v1/orders/{uuid}/cancel: cancel order")

	if serveErr := httpServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", serveErr)
	}

	log.Println("Order Service stopped")
}
