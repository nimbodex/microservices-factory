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

	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.UUID == orderUUID && order.Status == model.StatusPaid
	})).Return(nil)

	mockPaymentClient := clientmocks.NewPaymentClient(s.T())
	mockPaymentClient.On("PayOrder", mock.Anything, orderUUID, client.PaymentMethodCard).Return(&client.PaymentResult{
		TransactionUUID: transactionUUID,
		Success:         true,
	}, nil)

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.PayOrder(ctx, req, params)

	s.NoError(err)
	s.NotNil(result)

	payResp, ok := result.(*orderv1.PayOrderResponse)
	s.True(ok)
	s.Equal(transactionUUID, payResp.TransactionUUID)

	mockRepo.AssertExpectations(s.T())
	mockPaymentClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_OrderNotFound() {
	ctx := context.Background()
	orderUUID := uuid.New()

	params := orderv1.PayOrderParams{
		OrderUUID: orderUUID,
	}

	req := &orderv1.PayOrderRequest{
		PaymentMethod: orderv1.PaymentMethodCARD,
	}

	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(nil, assert.AnError)

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.PayOrder(ctx, req, params)

	s.NoError(err)
	s.NotNil(result)

	notFoundErr, ok := result.(*orderv1.NotFoundError)
	s.True(ok)
	s.Equal("order_not_found", notFoundErr.Error)
	s.Equal("order not found", notFoundErr.Message)

	mockRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_InvalidStatus() {
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
		Status:    model.StatusPaid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.PayOrder(ctx, req, params)

	s.NoError(err)
	s.NotNil(result)

	conflictErr, ok := result.(*orderv1.ConflictError)
	s.True(ok)
	s.Equal("invalid_status", conflictErr.Error)
	s.Equal("order cannot be paid", conflictErr.Message)

	mockRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestPayOrder_PaymentFailed() {
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

	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("GetByUUID", mock.Anything, orderUUID).Return(existingOrder, nil)

	mockPaymentClient := clientmocks.NewPaymentClient(s.T())
	mockPaymentClient.On("PayOrder", mock.Anything, orderUUID, client.PaymentMethodCard).Return(nil, assert.AnError)

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.PayOrder(ctx, req, params)

	s.NoError(err)
	s.NotNil(result)

	internalErr, ok := result.(*orderv1.InternalServerError)
	s.True(ok)
	s.Equal("payment_failed", internalErr.Error)
	s.Equal("payment processing failed", internalErr.Message)

	mockRepo.AssertExpectations(s.T())
	mockPaymentClient.AssertExpectations(s.T())
}
