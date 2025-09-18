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

	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, partUUID).Return(expectedPart, nil)

	service := NewInventoryService(mockRepo)

	result, err := service.GetPart(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.NotNil(result.Part)
	s.Equal(partUUID.String(), result.Part.Uuid)
	s.Equal("Test Part", result.Part.Name)
	s.Equal("Test Description", result.Part.Description)
	s.Equal(100.0, result.Part.Price)
	s.Equal(int64(10), result.Part.StockQuantity)
	s.Equal(inventoryv1.Category_CATEGORY_ENGINE, result.Part.Category)

	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestGetPart_EmptyUUID() {
	ctx := context.Background()

	req := &inventoryv1.GetPartRequest{
		Uuid: "",
	}

	mockRepo := repomocks.NewPartRepository(s.T())

	service := NewInventoryService(mockRepo)

	result, err := service.GetPart(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestGetPart_InvalidUUID() {
	ctx := context.Background()

	req := &inventoryv1.GetPartRequest{
		Uuid: "invalid-uuid",
	}

	mockRepo := repomocks.NewPartRepository(s.T())

	service := NewInventoryService(mockRepo)

	result, err := service.GetPart(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestGetPart_NotFound() {
	ctx := context.Background()
	partUUID := uuid.New()

	req := &inventoryv1.GetPartRequest{
		Uuid: partUUID.String(),
	}

	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, partUUID).Return(nil, assert.AnError)

	service := NewInventoryService(mockRepo)

	result, err := service.GetPart(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}
