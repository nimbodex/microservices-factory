package inventory

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nimbodex/microservices-factory/inventory/internal/converter"
	"github.com/nimbodex/microservices-factory/inventory/internal/repository"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

// InventoryServiceImpl implements InventoryService interface
type InventoryServiceImpl struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partRepo repository.PartRepository
}

// NewInventoryService creates a new inventory service instance
func NewInventoryService(partRepo repository.PartRepository) *InventoryServiceImpl {
	return &InventoryServiceImpl{
		partRepo: partRepo,
	}
}

func (s *InventoryServiceImpl) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	log.Printf("GetPart request received for UUID: %s", req.Uuid)

	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "UUID cannot be empty")
	}

	partUUID, err := uuid.Parse(req.Uuid)
	if err != nil {
		log.Printf("Invalid UUID format: %s, error: %v", req.Uuid, err)
		return nil, status.Error(codes.InvalidArgument, "invalid UUID format")
	}

	part, err := s.partRepo.GetByUUID(ctx, partUUID)
	if err != nil {
		log.Printf("Part not found for UUID: %s, error: %v", req.Uuid, err)
		return nil, status.Error(codes.NotFound, "part not found")
	}

	log.Printf("Part found: %s (%s)", part.Name, part.UUID)

	protoPart := converter.ToProtoPart(part)
	return &inventoryv1.GetPartResponse{
		Part: protoPart,
	}, nil
}

func (s *InventoryServiceImpl) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	log.Printf("ListParts request received with filter: %+v", req.Filter)

	filter := converter.ToServiceFilter(req.Filter)

	parts, err := s.partRepo.List(ctx, filter)
	if err != nil {
		log.Printf("Error retrieving parts: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	log.Printf("Found %d parts matching the filter", len(parts))

	protoParts := make([]*inventoryv1.Part, len(parts))
	for i, part := range parts {
		protoParts[i] = converter.ToProtoPart(part)
	}

	return &inventoryv1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
