package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nimbodex/microservices-factory/order/internal/client"
	clientmocks "github.com/nimbodex/microservices-factory/order/internal/client/mocks"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	repomocks "github.com/nimbodex/microservices-factory/order/internal/repository/mocks"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

func (s *OrderServiceTestSuite) TestPayOrder_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	params := orderv1.PayOrderParams{
		OrderUUID: orderUUID,
	}

	req := &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
	}

	existingOrder := &model.Order{
		UUID:      orderUUID,
		UserUUID:  userUUID,
		PartUUIDs: []uuid.UUID{uuid.New()},
		Status:    model.StatusPendingPayment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedOrder := *existingOrder
	updatedOrder.Status = model.StatusPaid
	updatedOrder.UpdatedAt = time.Now()

	// Mock repository
	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.UUID == orderUUID && order.Status == model.StatusPaid
	})).Return(nil)

	// Mock payment client
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())
	mockPaymentClient.On("PayOrder", mock.Anything, orderUUID, client.PaymentMethodCard).Return(&client.PaymentResult{
		TransactionUUID: transactionUUID,
		Success:         true,
	}, nil)

	// Mock inventory client (not used in pay)
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.PayOrder(ctx, req, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	payResp, ok := result.(*orderv1.PayOrderResponse)
	s.True(ok)
	s.Equal(transactionUUID, payResp.TransactionUUID)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
	mockPaymentClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_OrderNotFound() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()

	params := orderv1.PayOrderParams{
		OrderUUID: orderUUID,
	}

	req := &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
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
	result, err := service.PayOrder(ctx, req, params)

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

func (s *OrderServiceTestSuite) TestPayOrder_InvalidStatus() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	params := orderv1.PayOrderParams{
		OrderUUID: orderUUID,
	}

	req := &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
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
	result, err := service.PayOrder(ctx, req, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	conflictErr, ok := result.(*orderv1.ConflictError)
	s.True(ok)
	s.Equal("invalid_status", conflictErr.Error)
	s.Equal("order cannot be paid", conflictErr.Message)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_PaymentFailed() {
	// Arrange
	ctx := context.Background()
	orderUUID := uuid.New()
	userUUID := uuid.New()

	params := orderv1.PayOrderParams{
		OrderUUID: orderUUID,
	}

	req := &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
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

	// Mock payment client to return error
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())
	mockPaymentClient.On("PayOrder", mock.Anything, orderUUID, client.PaymentMethodCard).Return(nil, assert.AnError)

	// Mock inventory client
	mockInventoryClient := clientmocks.NewInventoryClient(s.T())

	// Create service
	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	// Act
	result, err := service.PayOrder(ctx, req, params)

	// Assert
	s.NoError(err)
	s.NotNil(result)

	internalErr, ok := result.(*orderv1.InternalServerError)
	s.True(ok)
	s.Equal("payment_failed", internalErr.Error)
	s.Equal("payment processing failed", internalErr.Message)

	// Verify mocks
	mockRepo.AssertExpectations(s.T())
	mockPaymentClient.AssertExpectations(s.T())
}
