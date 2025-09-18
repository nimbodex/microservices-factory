package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/payment/internal/model"
)

// PaymentRepository defines the interface for payment repository operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Payment, error)
	GetByOrderUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Payment, error)
	GetByTransactionUUID(ctx context.Context, transactionUUID uuid.UUID) (*model.Payment, error)
	Update(ctx context.Context, payment *model.Payment) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
