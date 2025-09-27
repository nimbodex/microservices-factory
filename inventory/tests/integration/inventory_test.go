//go:build integration

package integration

import (
	"time"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
	"github.com/nimbodex/microservices-factory/inventory/internal/repository"
	"github.com/nimbodex/microservices-factory/inventory/internal/service"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

func (suite *InventoryIntegrationTestSuite) TestInventoryService_GetPart() {
	// Arrange
	partRepo := repository.NewPartRepository(suite.env.Client(), "test_inventory")
	inventoryService := service.NewInventoryService(partRepo)

	// Create a test part first
	testPart := &model.Part{
		UUID:          uuid.New().String(),
		Name:          TestPartName,
		Description:   TestPartDescription,
		Price:         TestPartPrice,
		StockQuantity: int32(TestPartQuantity),
		Category:      inventoryv1.Category_CATEGORY_ENGINE,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := partRepo.Create(suite.ctx, testPart)
	suite.Require().NoError(err)

	// Act
	result, err := inventoryService.GetPart(suite.ctx, &inventoryv1.GetPartRequest{
		Uuid: testPart.UUID,
	})

	// Assert
	if err != nil {
		suite.Failf("GetPart failed", "Error: %v", err)
	}
	suite.NotNil(result)
	if result != nil {
		suite.Equal(testPart.UUID, result.Part.Uuid)
		suite.Equal(testPart.Name, result.Part.Name)
		suite.Equal(testPart.Description, result.Part.Description)
		suite.Equal(testPart.Price, result.Part.Price)
		suite.Equal(int64(testPart.StockQuantity), result.Part.StockQuantity)
	}
}

func (suite *InventoryIntegrationTestSuite) TestInventoryService_GetPart_NotFound() {
	// Arrange
	partRepo := repository.NewPartRepository(suite.env.Client(), "test_inventory")
	inventoryService := service.NewInventoryService(partRepo)

	nonExistentUUID := uuid.New().String()

	// Act
	result, err := inventoryService.GetPart(suite.ctx, &inventoryv1.GetPartRequest{
		Uuid: nonExistentUUID,
	})

	// Assert
	suite.Error(err)
	suite.Nil(result)
}

func (suite *InventoryIntegrationTestSuite) TestInventoryService_ListParts() {
	// Arrange
	partRepo := repository.NewPartRepository(suite.env.Client(), "test_inventory")
	inventoryService := service.NewInventoryService(partRepo)

	// Create test parts
	parts := []*model.Part{
		{
			UUID:          uuid.New().String(),
			Name:          "Part 1",
			Description:   "Description 1",
			Price:         100.0,
			StockQuantity: 10,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			UUID:          uuid.New().String(),
			Name:          "Part 2",
			Description:   "Description 2",
			Price:         200.0,
			StockQuantity: 20,
			Category:      inventoryv1.Category_CATEGORY_FUEL,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	for _, part := range parts {
		err := partRepo.Create(suite.ctx, part)
		suite.Require().NoError(err)
	}

	// Act
	result, err := inventoryService.ListParts(suite.ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{},
	})

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result.Parts, 2)

	// Verify all created parts are in the result
	partUUIDs := make(map[string]bool)
	for _, part := range result.Parts {
		partUUIDs[part.Uuid] = true
	}

	for _, originalPart := range parts {
		suite.True(partUUIDs[originalPart.UUID], "Part %s not found in result", originalPart.UUID)
	}
}

func (suite *InventoryIntegrationTestSuite) TestInventoryService_ListParts_EmptyResult() {
	// Arrange
	partRepo := repository.NewPartRepository(suite.env.Client(), "test_inventory")
	inventoryService := service.NewInventoryService(partRepo)

	// Act
	result, err := inventoryService.ListParts(suite.ctx, &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{},
	})

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Len(result.Parts, 0)
}
