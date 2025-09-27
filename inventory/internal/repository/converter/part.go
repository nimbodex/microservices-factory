package converter

import (
	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
	repomodel "github.com/nimbodex/microservices-factory/inventory/internal/repository/model"
)

// ToRepoPart converts service model to repository model
func ToRepoPart(servicePart *model.Part) *repomodel.Part {
	var dimensions *repomodel.Dimensions
	if servicePart.Dimensions != nil {
		dimensions = &repomodel.Dimensions{
			Length: servicePart.Dimensions.Length,
			Width:  servicePart.Dimensions.Width,
			Height: servicePart.Dimensions.Height,
			Weight: servicePart.Dimensions.Weight,
		}
	}

	var manufacturer *repomodel.Manufacturer
	if servicePart.Manufacturer != nil {
		manufacturer = &repomodel.Manufacturer{
			Name:    servicePart.Manufacturer.Name,
			Country: servicePart.Manufacturer.Country,
			Website: servicePart.Manufacturer.Website,
		}
	}

	return &repomodel.Part{
		UUID:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: servicePart.StockQuantity,
		Category:      servicePart.Category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          servicePart.Tags,
		Metadata:      servicePart.Metadata,
		CreatedAt:     servicePart.CreatedAt,
		UpdatedAt:     servicePart.UpdatedAt,
	}
}

// FromRepoPart converts repository model to service model
func FromRepoPart(repoPart *repomodel.Part) (*model.Part, error) {
	// Validate UUID format
	_, err := uuid.Parse(repoPart.UUID)
	if err != nil {
		return nil, err
	}

	var dimensions *model.Dimensions
	if repoPart.Dimensions != nil {
		dimensions = &model.Dimensions{
			Length: repoPart.Dimensions.Length,
			Width:  repoPart.Dimensions.Width,
			Height: repoPart.Dimensions.Height,
			Weight: repoPart.Dimensions.Weight,
		}
	}

	var manufacturer *model.Manufacturer
	if repoPart.Manufacturer != nil {
		manufacturer = &model.Manufacturer{
			Name:    repoPart.Manufacturer.Name,
			Country: repoPart.Manufacturer.Country,
			Website: repoPart.Manufacturer.Website,
		}
	}

	return &model.Part{
		UUID:          repoPart.UUID,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      repoPart.Category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          repoPart.Tags,
		Metadata:      repoPart.Metadata,
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}, nil
}
