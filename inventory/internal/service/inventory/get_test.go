package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nimbodex/microservices-factory/inventory/internal/model"
	repomocks "github.com/nimbodex/microservices-factory/inventory/internal/repository/mocks"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
)

func (s *InventoryServiceTestSuite) TestGetPart_Success() {
	// Arrange
	ctx := context.Background()
	partUUID := uuid.New()

	req := &inventoryv1.GetPartRequest{
		Uuid: partUUID.String(),
	}

	expectedPart := &model.Part{
		UUID:          partUUID,
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.0,
		StockQuantity: 10,
		Category:      inventoryv1.Category_CATEGORY_ENGINE,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Mock repository
	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, partUUID).Return(expectedPart, nil)

	// Create service
	service := NewInventoryService(mockRepo)

	// Act
	result, err := service.GetPart(ctx, req)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.Part)
	s.Equal(partUUID.String(), result.Part.Uuid)
	s.Equal("Test Part", result.Part.Name)
	s.Equal("Test Description", result.Part.Description)
	s.Equal(100.0, result.Part.Price)
	s.Equal(int64(10), result.Part.StockQuantity)
	s.Equal(inventoryv1.Category_CATEGORY_ENGINE, result.Part.Category)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestGetPart_EmptyUUID() {
	// Arrange
	ctx := context.Background()

	req := &inventoryv1.GetPartRequest{
		Uuid: "",
	}

	// Mock repository (should not be called)
	mockRepo := repomocks.NewPartRepository(s.T())

	// Create service
	service := NewInventoryService(mockRepo)

	// Act
	result, err := service.GetPart(ctx, req)

	// Assert
	s.Error(err)
	s.Nil(result)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestGetPart_InvalidUUID() {
	// Arrange
	ctx := context.Background()

	req := &inventoryv1.GetPartRequest{
		Uuid: "invalid-uuid",
	}

	// Mock repository (should not be called)
	mockRepo := repomocks.NewPartRepository(s.T())

	// Create service
	service := NewInventoryService(mockRepo)

	// Act
	result, err := service.GetPart(ctx, req)

	// Assert
	s.Error(err)
	s.Nil(result)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestGetPart_NotFound() {
	// Arrange
	ctx := context.Background()
	partUUID := uuid.New()

	req := &inventoryv1.GetPartRequest{
		Uuid: partUUID.String(),
	}

	// Mock repository to return error
	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, partUUID).Return(nil, assert.AnError)

	// Create service
	service := NewInventoryService(mockRepo)

	// Act
	result, err := service.GetPart(ctx, req)

	// Assert
	s.Error(err)
	s.Nil(result)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}
