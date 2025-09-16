package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clientmocks "github.com/nimbodex/microservices-factory/order/internal/client/mocks"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	repomocks "github.com/nimbodex/microservices-factory/order/internal/repository/mocks"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

func (s *OrderServiceTestSuite) TestCancelOrder_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	params := orderv1.CancelOrderParams{
		OrderUUID: orderUUID,
	}

	existingOrder := &model.Order{
		UUID:      orderUUID,
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{uuid.New()},
		Status:    model.StatusPendingPayment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository
	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.UUID == orderUUID && order.Status == model.StatusCancelled
	})).Return(nil)

	// Mock clients (not used in cancel)
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.CancelOrder(ctx, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	_, ok := result.(*orderv1.CancelOrderNoContent)
	s.True(ok)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCancelOrder_OrderNotFound() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()

	params := orderv1.CancelOrderParams{
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
	result, err := service.CancelOrder(ctx, params)

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

func (s *OrderServiceTestSuite) TestCancelOrder_InvalidStatus() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	params := orderv1.CancelOrderParams{
		OrderUUID: orderUUID,
	}

	// Order already paid
	existingOrder := &model.Order{
		UUID:      orderUUID,
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{uuid.New()},
		Status:    model.StatusPaid, // Already paid
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository
	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)

	// Mock clients
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.CancelOrder(ctx, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	conflictErr, ok := result.(*orderv1.ConflictError)
	s.True(ok)
	s.Equal("invalid_status", conflictErr.Error)
	s.Equal("order cannot be cancelled", conflictErr.Message)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCancelOrder_UpdateFailed() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	params := orderv1.CancelOrderParams{
		OrderUUID: orderUUID,
	}

	existingOrder := &model.Order{
		UUID:      orderUUID,
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{uuid.New()},
		Status:    model.StatusPendingPayment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock repository
	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(assert.AnError)

	// Mock clients
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.CancelOrder(ctx, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	internalErr, ok := result.(*orderv1.InternalServerError)
	s.True(ok)
	s.Equal("update_failed", internalErr.Error)
	s.Equal("failed to update order status", internalErr.Message)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}
