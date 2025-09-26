package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nimbodex/microservices-factory/payment/internal/model"
	repomocks "github.com/nimbodex/microservices-factory/payment/internal/repository/mocks"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

func (s *PaymentServiceTestSuite) TestPayOrder_Success() {
	ctx := context.Background()
	orderUUID := uuid.New()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockRepo := repomocks.NewPaymentRepository(s.T())
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(payment *model.Payment) bool {
		return payment.OrderUUID == orderUUID &&
			payment.PaymentMethod == model.PaymentMethodCard &&
			payment.Status == model.PaymentStatusCompleted
	})).Return(nil)

	service := NewPaymentService(mockRepo)

	result, err := service.PayOrder(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.NotEmpty(result.TransactionUuid)

	_, err = uuid.Parse(result.TransactionUuid)
	s.NoError(err)

	mockRepo.AssertExpectations(s.T())
}

func (s *PaymentServiceTestSuite) TestPayOrder_SBP_Success() {
	ctx := context.Background()
	orderUUID := uuid.New()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
	}

	mockRepo := repomocks.NewPaymentRepository(s.T())
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(payment *model.Payment) bool {
		return payment.OrderUUID == orderUUID &&
			payment.PaymentMethod == model.PaymentMethodSBP &&
			payment.Status == model.PaymentStatusCompleted
	})).Return(nil)

	service := NewPaymentService(mockRepo)

	result, err := service.PayOrder(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.NotEmpty(result.TransactionUuid)

	mockRepo.AssertExpectations(s.T())
}

func (s *PaymentServiceTestSuite) TestPayOrder_InvalidOrderUUID() {
	ctx := context.Background()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     "invalid-uuid",
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockRepo := repomocks.NewPaymentRepository(s.T())

	service := NewPaymentService(mockRepo)

	result, err := service.PayOrder(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}

func (s *PaymentServiceTestSuite) TestPayOrder_UnknownPaymentMethod() {
	ctx := context.Background()
	orderUUID := uuid.New()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_UNKNOWN,
	}

	mockRepo := repomocks.NewPaymentRepository(s.T())

	service := NewPaymentService(mockRepo)

	result, err := service.PayOrder(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}

func (s *PaymentServiceTestSuite) TestPayOrder_RepositoryError() {
	ctx := context.Background()
	orderUUID := uuid.New()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockRepo := repomocks.NewPaymentRepository(s.T())
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(assert.AnError)

	service := NewPaymentService(mockRepo)

	result, err := service.PayOrder(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}

func (s *PaymentServiceTestSuite) TestPayOrder_EmptyOrderUUID() {
	ctx := context.Background()

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     "",
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	mockRepo := repomocks.NewPaymentRepository(s.T())

	service := NewPaymentService(mockRepo)

	result, err := service.PayOrder(ctx, req)

	s.Error(err)
	s.Nil(result)

	mockRepo.AssertExpectations(s.T())
}
