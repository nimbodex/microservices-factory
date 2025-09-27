package service

import (
	"context"

	"github.com/nimbodex/microservices-factory/inventory/internal/repository"
	"github.com/nimbodex/microservices-factory/inventory/internal/service/inventory"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// InventoryService defines the interface for inventory service operations
type InventoryService interface {
	GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error)
	ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error)
}

// NewInventoryService creates a new inventory service instance
func NewInventoryService(partRepo repository.PartRepository) InventoryService {
	return inventory.NewInventoryService(partRepo)
}
