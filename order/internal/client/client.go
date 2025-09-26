package client

import (
	"context"

	"github.com/google/uuid"
)

// InventoryClient defines the interface for inventory service client
type InventoryClient interface {
	GetPart(ctx context.Context, partUUID uuid.UUID) (*Part, error)
	ListParts(ctx context.Context, limit, offset int) ([]*Part, error)
}

// PaymentClient defines the interface for payment service client
type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod PaymentMethod) (*PaymentResult, error)
}

// Part represents a part from inventory service
type Part struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

// PaymentMethod represents payment method
type PaymentMethod string

const (
	PaymentMethodCard    PaymentMethod = "CARD"
	PaymentMethodSBP     PaymentMethod = "SBP"
	PaymentMethodUnknown PaymentMethod = "UNKNOWN"
)

// PaymentResult represents the result of payment processing
type PaymentResult struct {
	TransactionUUID uuid.UUID `json:"transaction_uuid"`
	Success         bool      `json:"success"`
	Message         string    `json:"message,omitempty"`
}
