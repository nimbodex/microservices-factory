package v1

import (
	"context"

	"github.com/nimbodex/microservices-factory/inventory/internal/service"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// APIHandler handles gRPC requests for inventory API
type APIHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	inventoryService service.InventoryService
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(inventoryService service.InventoryService) *APIHandler {
	return &APIHandler{
		inventoryService: inventoryService,
	}
}

func (h *APIHandler) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	return h.inventoryService.GetPart(ctx, req)
}

func (h *APIHandler) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	return h.inventoryService.ListParts(ctx, req)
}
