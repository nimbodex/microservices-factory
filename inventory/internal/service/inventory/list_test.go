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

func (s *InventoryServiceTestSuite) TestListParts_Success() {
	ctx := context.Background()

	req := &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Categories: []inventoryv1.Category{inventoryv1.Category_CATEGORY_ENGINE},
		},
	}

	expectedParts := []*model.Part{
		{
			UUID:          uuid.New().String(),
			Name:          "Engine Part 1",
			Description:   "Engine Description 1",
			Price:         100.0,
			StockQuantity: 10,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			UUID:          uuid.New().String(),
			Name:          "Engine Part 2",
			Description:   "Engine Description 2",
			Price:         200.0,
			StockQuantity: 5,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("List", mock.Anything, mock.MatchedBy(func(filter *model.PartsFilter) bool {
		return len(filter.Categories) == 1 && filter.Categories[0] == inventoryv1.Category_CATEGORY_ENGINE
	})).Return(expectedParts, nil)

	service := NewInventoryService(mockRepo)

	result, err := service.ListParts(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Parts, 2)
	s.Equal("Engine Part 1", result.Parts[0].Name)
	s.Equal("Engine Part 2", result.Parts[1].Name)

	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestListParts_EmptyFilter() {
	ctx := context.Background()

	req := &inventoryv1.ListPartsRequest{
		Filter: nil,
	}

	expectedParts := []*model.Part{
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
	}

	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("List", mock.Anything, mock.MatchedBy(func(filter *model.PartsFilter) bool {
		return filter != nil
	})).Return(expectedParts, nil)

	service := NewInventoryService(mockRepo)

	result, err := service.ListParts(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Parts, 1)
	s.Equal("Part 1", result.Parts[0].Name)

	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestListParts_RepositoryError() {
	ctx := context.Background()

	req := &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Categories: []inventoryv1.Category{inventoryv1.Category_CATEGORY_ENGINE},
		},
	}

	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("List", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	service := NewInventoryService(mockRepo)

	result, err := service.ListParts(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}

func (s *InventoryServiceTestSuite) TestListParts_EmptyResult() {
	ctx := context.Background()

	req := &inventoryv1.ListPartsRequest{
		Filter: &inventoryv1.PartsFilter{
			Categories: []inventoryv1.Category{inventoryv1.Category_CATEGORY_UNKNOWN},
		},
	}

	mockRepo := repomocks.NewPartRepository(s.T())
	mockRepo.On("List", mock.Anything, mock.Anything).Return([]*model.Part{}, nil)

	service := NewInventoryService(mockRepo)

	result, err := service.ListParts(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Parts, 0)

	mockRepo.AssertExpectations(s.T())
}
