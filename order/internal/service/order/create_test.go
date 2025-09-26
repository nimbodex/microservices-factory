package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nimbodex/microservices-factory/order/internal/client"
	clientmocks "github.com/nimbodex/microservices-factory/order/internal/client/mocks"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	repomocks "github.com/nimbodex/microservices-factory/order/internal/repository/mocks"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

func (s *OrderServiceTestSuite) TestCreateOrder_Success() {
	ctx := context.Background()
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	req := &orderv1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID1, partUUID2},
	}

	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUUIDs) == 2 &&
			order.Status == model.StatusPendingPayment
	})).Return(nil)

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockInventoryClient.On("GetPart", mock.Anything, partUUID1).Return(&client.Part{
		UUID:  partUUID1,
		Name:  "Part 1",
		Price: 100.0,
	}, nil)
	mockInventoryClient.On("GetPart", mock.Anything, partUUID2).Return(&client.Part{
		UUID:  partUUID2,
		Name:  "Part 2",
		Price: 200.0,
	}, nil)

	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.CreateOrder(ctx, req)

	s.NoError(err)
	s.NotNil(result)

	createResp, ok := result.(*orderv1.CreateOrderResponse)
	s.True(ok)
	s.NotEmpty(createResp.OrderUUID)
	s.Equal(0.0, createResp.TotalPrice)

	mockRepo.AssertExpectations(s.T())
	mockInventoryClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCreateOrder_PartNotFound() {
	ctx := context.Background()
	userUUID := uuid.New()
	partUUID := uuid.New()

	req := &orderv1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID},
	}

	mockRepo := repomocks.NewOrderRepository(s.T())

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockInventoryClient.On("GetPart", mock.Anything, partUUID).Return(nil, assert.AnError)

	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.CreateOrder(ctx, req)

	s.NoError(err)
	s.NotNil(result)

	badReqErr, ok := result.(*orderv1.BadRequestError)
	s.True(ok)
	s.Equal("part_not_found", badReqErr.Error)
	s.Contains(badReqErr.Message, "part")

	mockRepo.AssertExpectations(s.T())
	mockInventoryClient.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestCreateOrder_RepositoryError() {
	ctx := context.Background()
	userUUID := uuid.New()
	partUUID := uuid.New()

	req := &orderv1.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{partUUID},
	}

	mockRepo := repomocks.NewOrderRepository(s.T())
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	mockInventoryClient := clientmocks.NewInventoryClient(s.T())
	mockInventoryClient.On("GetPart", mock.Anything, partUUID).Return(&client.Part{
		UUID:  partUUID,
		Name:  "Part 1",
		Price: 100.0,
	}, nil)

	mockPaymentClient := clientmocks.NewPaymentClient(s.T())

	service := NewOrderService(mockRepo, mockInventoryClient, mockPaymentClient)

	result, err := service.CreateOrder(ctx, req)

	s.NoError(err)
	s.NotNil(result)

	internalErr, ok := result.(*orderv1.InternalServerError)
	s.True(ok)
	s.Equal("creation_failed", internalErr.Error)
	s.Equal("failed to create order", internalErr.Message)

	mockRepo.AssertExpectations(s.T())
	mockInventoryClient.AssertExpectations(s.T())
}
