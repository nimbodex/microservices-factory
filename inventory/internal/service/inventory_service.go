package service

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nexarise/microservices-factory/inventory/internal/storage"
	inventoryv1 "github.com/nexarise/microservices-factory/shared/pkg/proto/inventory/v1"
)

type InventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
	storage *storage.MemoryStorage
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		storage: storage.NewMemoryStorage(),
	}
}

func (s *InventoryService) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	log.Printf("GetPart request received for UUID: %s", req.Uuid)

	if req.Uuid == "" {
		return nil, status.Error(codes.InvalidArgument, "UUID cannot be empty")
	}

	part, err := s.storage.GetPart(req.Uuid)
	if err != nil {
		log.Printf("Part not found for UUID: %s, error: %v", req.Uuid, err)
		return nil, status.Error(codes.NotFound, "part not found")
	}

	log.Printf("Part found: %s (%s)", part.Name, part.Uuid)

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	log.Printf("ListParts request received with filter: %+v", req.Filter)

	parts, err := s.storage.ListParts(req.Filter)
	if err != nil {
		log.Printf("Error retrieving parts: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	log.Printf("Found %d parts matching the filter", len(parts))

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}
