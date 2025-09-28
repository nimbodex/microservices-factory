package converter

import (
	"fmt"
	"math"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

func safeInt64ToInt32(value int64) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("integer overflow: value %d is out of int32 range", value)
	}

	return int32(value), nil
}

// ToServicePart converts protobuf part to service model
func ToServicePart(protoPart *inventoryv1.Part) (*model.Part, error) {
	if protoPart == nil {
		return nil, fmt.Errorf("protoPart cannot be nil")
	}

	_, err := uuid.Parse(protoPart.Uuid)
	if err != nil {
		return nil, err
	}

	var dimensions *model.Dimensions
	if protoPart.Dimensions != nil {
		dimensions = &model.Dimensions{
			Length: protoPart.Dimensions.Length,
			Width:  protoPart.Dimensions.Width,
			Height: protoPart.Dimensions.Height,
			Weight: protoPart.Dimensions.Weight,
		}
	}

	var manufacturer *model.Manufacturer
	if protoPart.Manufacturer != nil {
		manufacturer = &model.Manufacturer{
			Name:    protoPart.Manufacturer.Name,
			Country: protoPart.Manufacturer.Country,
			Website: protoPart.Manufacturer.Website,
		}
	}

	metadata := make(map[string]interface{})
	for k, v := range protoPart.Metadata {
		metadata[k] = v.AsInterface()
	}

	stockQuantity, err := safeInt64ToInt32(protoPart.StockQuantity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert stock quantity: %w", err)
	}

	return &model.Part{
		UUID:          protoPart.Uuid,
		Name:          protoPart.Name,
		Description:   protoPart.Description,
		Price:         protoPart.Price,
		StockQuantity: stockQuantity,
		Category:      protoPart.Category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          protoPart.Tags,
		Metadata:      metadata,
		CreatedAt:     protoPart.CreatedAt.AsTime(),
		UpdatedAt:     protoPart.UpdatedAt.AsTime(),
	}, nil
}

// ToProtoPart converts service model to protobuf part
func ToProtoPart(servicePart *model.Part) *inventoryv1.Part {
	if servicePart == nil {
		return nil
	}

	var dimensions *inventoryv1.Dimensions
	if servicePart.Dimensions != nil {
		dimensions = &inventoryv1.Dimensions{
			Length: servicePart.Dimensions.Length,
			Width:  servicePart.Dimensions.Width,
			Height: servicePart.Dimensions.Height,
			Weight: servicePart.Dimensions.Weight,
		}
	}

	var manufacturer *inventoryv1.Manufacturer
	if servicePart.Manufacturer != nil {
		manufacturer = &inventoryv1.Manufacturer{
			Name:    servicePart.Manufacturer.Name,
			Country: servicePart.Manufacturer.Country,
			Website: servicePart.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*structpb.Value)
	if servicePart.Metadata != nil {
		for k, v := range servicePart.Metadata {
			if val, err := structpb.NewValue(v); err == nil {
				metadata[k] = val
			}
		}
	}

	return &inventoryv1.Part{
		Uuid:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: int64(servicePart.StockQuantity),
		Category:      servicePart.Category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          servicePart.Tags,
		Metadata:      metadata,
		CreatedAt:     timestamppb.New(servicePart.CreatedAt),
		UpdatedAt:     timestamppb.New(servicePart.UpdatedAt),
	}
}

// ToServiceFilter converts protobuf filter to service model
func ToServiceFilter(protoFilter *inventoryv1.PartsFilter) *model.PartsFilter {
	if protoFilter == nil {
		return &model.PartsFilter{}
	}

	return &model.PartsFilter{
		UUIDs:                 protoFilter.Uuids,
		Names:                 protoFilter.Names,
		Categories:            protoFilter.Categories,
		ManufacturerCountries: protoFilter.ManufacturerCountries,
		Tags:                  protoFilter.Tags,
	}
}

// ToProtoFilter converts service model to protobuf filter
func ToProtoFilter(serviceFilter *model.PartsFilter) *inventoryv1.PartsFilter {
	if serviceFilter == nil {
		return &inventoryv1.PartsFilter{}
	}

	return &inventoryv1.PartsFilter{
		Uuids:                 serviceFilter.UUIDs,
		Names:                 serviceFilter.Names,
		Categories:            serviceFilter.Categories,
		ManufacturerCountries: serviceFilter.ManufacturerCountries,
		Tags:                  serviceFilter.Tags,
	}
}
