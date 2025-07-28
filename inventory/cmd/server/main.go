package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nexarise/microservices-factory/inventory/internal/service"
	inventoryv1 "github.com/nexarise/microservices-factory/shared/pkg/proto/inventory/v1"
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

	inventoryService := service.NewInventoryService()
	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryService)

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
