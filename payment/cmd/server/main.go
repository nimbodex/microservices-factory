package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/nexarise/microservices-factory/payment/internal/service"
	paymentv1 "github.com/nexarise/microservices-factory/shared/pkg/proto/payment/v1"
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

	paymentService := service.NewPaymentService()
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentService)

	reflection.Register(grpcServer)

	log.Printf("Payment Service listening on %s", port)
	log.Println("Available methods:")
	log.Println("\t - PayOrder: processing the order payment command")
	log.Println("For testing use grpcurl or any gRPC client")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
