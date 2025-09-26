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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	v1 "github.com/nimbodex/microservices-factory/order/internal/api/order/v1"
	"github.com/nimbodex/microservices-factory/order/internal/client/grpc"
	"github.com/nimbodex/microservices-factory/order/internal/migrator"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("DB_URI")

	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		closeErr := con.Close(ctx)
		if closeErr != nil {
			log.Printf("failed to close connection: %v\n", closeErr)
		}
	}()

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("Database is unavailable: %v\n", err)
		return
	}

	db := stdlib.OpenDB(*con.Config().Copy())
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(db, migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Database migration error: %v\n", err)
		return
	}

	orderRepo := orderrepo.NewPostgresOrderRepository(db)

	inventoryClient, err := grpc.NewGRPCInventoryClient()
	if err != nil {
		log.Printf("Failed to create inventory client: %v", err)
		return
	}

	paymentClient, err := grpc.NewGRPCPaymentClient()
	if err != nil {
		log.Printf("Failed to create payment client: %v", err)
		return
	}

	orderService := orderservice.NewOrderService(orderRepo, inventoryClient, paymentClient)

	apiHandler := v1.NewAPIHandler(orderService)

	server, err := orderv1.NewServer(apiHandler)
	if err != nil {
		log.Printf("Failed to create server: %v", err)
		return
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

		cancel()

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

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
			log.Printf("Server shutdown error: %v", shutdownErr)
		}
	}()

	log.Printf("Order Service listening on %s", port)
	log.Println("Available endpoints:")
	log.Println("\t - POST /api/v1/orders: create order")
	log.Println("\t - GET /api/v1/orders/{uuid}: get order")
	log.Println("\t - POST /api/v1/orders/{uuid}/pay: pay order")
	log.Println("\t - POST /api/v1/orders/{uuid}/cancel: cancel order")

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- httpServer.ListenAndServe()
	}()

	select {
	case serveErr := <-serverErrChan:
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Printf("Server error: %v", serveErr)
			return
		}
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	log.Println("Order Service stopped")
}
