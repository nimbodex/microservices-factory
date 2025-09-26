package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
)

// PartRepository defines the interface for part repository operations
type PartRepository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Part, error)
	List(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
	Create(ctx context.Context, part *model.Part) error
	Update(ctx context.Context, part *model.Part) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}
