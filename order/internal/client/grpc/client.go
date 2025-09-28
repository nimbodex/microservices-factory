package grpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nimbodex/microservices-factory/order/internal/client"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

const (
	InventoryServiceAddr = "localhost:50051"
	PaymentServiceAddr   = "localhost:50052"
)

// GRPCInventoryClient implements InventoryClient using gRPC
type GRPCInventoryClient struct {
	client inventoryv1.InventoryServiceClient
	conn   *grpc.ClientConn
}

// GRPCPaymentClient implements PaymentClient using gRPC
type GRPCPaymentClient struct {
	client paymentv1.PaymentServiceClient
	conn   *grpc.ClientConn
}

// NewGRPCInventoryClient creates a new gRPC inventory client
func NewGRPCInventoryClient() (*GRPCInventoryClient, error) {
	conn, err := grpc.NewClient(InventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	client := inventoryv1.NewInventoryServiceClient(conn)

	return &GRPCInventoryClient{
		client: client,
		conn:   conn,
	}, nil
}

// NewGRPCPaymentClient creates a new gRPC payment client
func NewGRPCPaymentClient() (*GRPCPaymentClient, error) {
	conn, err := grpc.NewClient(PaymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	client := paymentv1.NewPaymentServiceClient(conn)

	return &GRPCPaymentClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *GRPCInventoryClient) GetPart(ctx context.Context, partUUID uuid.UUID) (*client.Part, error) {
	resp, err := c.client.GetPart(ctx, &inventoryv1.GetPartRequest{
		Uuid: partUUID.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get part %s: %w", partUUID, err)
	}

	return &client.Part{
		UUID:  partUUID,
		Name:  resp.Part.Name,
		Price: resp.Part.Price,
	}, nil
}

func (c *GRPCInventoryClient) ListParts(ctx context.Context, limit, offset int) ([]*client.Part, error) {
	resp, err := c.client.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list parts: %w", err)
	}

	parts := make([]*client.Part, len(resp.Parts))
	for i, part := range resp.Parts {
		partUUID, err := uuid.Parse(part.Uuid)
		if err != nil {
			log.Printf("Failed to parse part UUID %s: %v", part.Uuid, err)
			continue
		}

		parts[i] = &client.Part{
			UUID:  partUUID,
			Name:  part.Name,
			Price: part.Price,
		}
	}

	return parts, nil
}

func (c *GRPCPaymentClient) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod client.PaymentMethod) (*client.PaymentResult, error) {
	var grpcPaymentMethod paymentv1.PaymentMethod
	switch paymentMethod {
	case client.PaymentMethodCard:
		grpcPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case client.PaymentMethodSBP:
		grpcPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		grpcPaymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}

	resp, err := c.client.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: grpcPaymentMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process payment for order %s: %w", orderUUID, err)
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction UUID %s: %w", resp.TransactionUuid, err)
	}

	return &client.PaymentResult{
		TransactionUUID: transactionUUID,
		Success:         true,
	}, nil
}

func (c *GRPCInventoryClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *GRPCPaymentClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
