package main

import (
	"log"
	"net"

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// Initialize repository
	partRepo := part.NewMemoryPartRepository()

	// Initialize service layer
	inventoryService := inventoryservice.NewInventoryService(partRepo)

	// Initialize API handler
	apiHandler := v1.NewAPIHandler(inventoryService)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, apiHandler)

	reflection.Register(grpcServer)

	log.Printf("Inventory Service listening on %s", port)
	log.Println("Available methods:")
	log.Println("\t - GetPart: getting a detail by UUID")
	log.Println("\t - ListParts: getting parts list with filtering")
	log.Println("For testing use grpcurl or any gRPC client")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
