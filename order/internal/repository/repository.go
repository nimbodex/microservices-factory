package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/order/internal/model"
)

// OrderRepository defines the interface for order repository operations
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*model.Order, error)
}
