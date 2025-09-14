package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clientmocks "github.com/nimbodex/microservices-factory/order/internal/client/mocks"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	repomocks "github.com/nimbodex/microservices-factory/order/internal/repository/mocks"
)

func (s *OrderServiceTestSuite) TestGetOrder_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	params := orderv1.GetOrderParams{
		OrderUUID: orderUUID,
	}

	expectedOrder := &model.Order{
		UUID:      orderUUID,
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{partUUID1, partUUID2},
		Status:    model.StatusPendingPayment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository
	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(expectedOrder, nil)

	// Mock clients (not used in get)
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.GetOrder(ctx, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	getResp, ok := result.(*orderv1.GetOrderResponse)
	s.True(ok)
	s.Equal(orderUUID, getResp.OrderUUID)
	s.Equal(userUUID, getResp.UserUUID)
	s.Len(getResp.PartUuids, 2)
	s.Equal(orderv1.OrderStatus(model.StatusPendingPayment), getResp.Status)
	s.Equal(0.0, getResp.TotalPrice) // Simplified calculation

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestGetOrder_NotFound() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()

	params := orderv1.GetOrderParams{
		OrderUUID: orderUUID,
	}

	// Mock repository to return error
	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(nil, assert.AnError)

	// Mock clients
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.GetOrder(ctx, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	notFoundErr, ok := result.(*orderv1.NotFoundError)
	s.True(ok)
	s.Equal("order_not_found", notFoundErr.Error)
	s.Equal("order not found", notFoundErr.Message)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}
