package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/nimbodex/microservices-factory/inventory/internal/api/inventory/v1"
	"github.com/nimbodex/microservices-factory/inventory/internal/repository/part"
	inventoryservice "github.com/nimbodex/microservices-factory/inventory/internal/service/inventory"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

const (
	port = "localhost:50051"
)

func main() {
	log.Println("Starting Inventory Service...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer func() {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Printf("failed to disconnect from MongoDB: %v", disconnectErr)
		}
	}()

	if pingErr := client.Ping(ctx, nil); pingErr != nil {
		log.Printf("failed to ping MongoDB: %v\n", pingErr)
		return
	}

	db := client.Database("inventory_db")
	collection := db.Collection("parts")
	partRepo := part.NewMongoPartRepository(collection)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Failed to listen on port %s: %v", port, err)
		return
	}

	grpcServer := grpc.NewServer()

	inventoryService := inventoryservice.NewInventoryService(partRepo)

	apiHandler := v1.NewAPIHandler(inventoryService)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, apiHandler)

	reflection.Register(grpcServer)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down Inventory Service...")

		cancel()

		grpcServer.GracefulStop()
		log.Println("gRPC server stopped")
	}()

	log.Printf("Inventory Service listening on %s", port)
	log.Println("Available methods:")
	log.Println("\t - GetPart: getting a detail by UUID")
	log.Println("\t - ListParts: getting parts list with filtering")
	log.Println("For testing use grpcurl or any gRPC client")

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- grpcServer.Serve(lis)
	}()

	select {
	case serveErr := <-serverErrChan:
		if serveErr != nil {
			log.Printf("gRPC server error: %v", serveErr)
			return
		}
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	log.Println("Inventory Service stopped")
}
