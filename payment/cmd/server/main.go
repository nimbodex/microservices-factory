package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/nimbodex/microservices-factory/payment/internal/api/payment/v1"
	"github.com/nimbodex/microservices-factory/payment/internal/repository/payment"
	paymentservice "github.com/nimbodex/microservices-factory/payment/internal/service/payment"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

const (
	port = "localhost:50052"
)

func main() {
	log.Println("Starting Payment Service...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	// Initialize repository
	paymentRepo := payment.NewMemoryPaymentRepository()

	// Initialize service layer
	paymentService := paymentservice.NewPaymentService(paymentRepo)

	// Initialize API handler
	apiHandler := v1.NewAPIHandler(paymentService)

	paymentv1.RegisterPaymentServiceServer(grpcServer, apiHandler)

	reflection.Register(grpcServer)

	log.Printf("Payment Service listening on %s", port)
	log.Println("Available methods:")
	log.Println("\t - PayOrder: processing the order payment command")
	log.Println("For testing use grpcurl or any gRPC client")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
