package payment

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/payment/internal/converter"
	"github.com/nimbodex/microservices-factory/payment/internal/model"
	"github.com/nimbodex/microservices-factory/payment/internal/repository"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

// PaymentServiceImpl implements PaymentService interface
type PaymentServiceImpl struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentRepo repository.PaymentRepository
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(paymentRepo repository.PaymentRepository) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		paymentRepo: paymentRepo,
	}
}

func (s *PaymentServiceImpl) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	log.Printf("Processing payment for order %s with method %s", req.OrderUuid, req.PaymentMethod)

	payReq, err := converter.ToServicePayOrderRequest(req)
	if err != nil {
		log.Printf("Failed to convert payment request: %v", err)
		return nil, err
	}

	if payReq.PaymentMethod == model.PaymentMethodUnknown {
		log.Printf("Invalid payment method: %s", req.PaymentMethod)
		return nil, model.NewInvalidPaymentMethodError(payReq.PaymentMethod)
	}

	if payReq.Amount < 0 {
		log.Printf("Invalid amount: %f", payReq.Amount)
		return nil, model.NewInvalidAmountError(payReq.Amount)
	}

	transactionUUID := uuid.New()
	paymentUUID := uuid.New()

	payment := &model.Payment{
		UUID:            paymentUUID,
		OrderUUID:       payReq.OrderUUID,
		PaymentMethod:   payReq.PaymentMethod,
		Amount:          payReq.Amount,
		Status:          model.PaymentStatusCompleted, // Simplified - always successful for now
		TransactionUUID: transactionUUID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		log.Printf("Failed to create payment: %v", err)
		return nil, model.NewInternalError(err)
	}

	log.Printf("Payment was successful, transaction_uuid: %s", transactionUUID)

	return converter.ToProtoPayOrderResponse(transactionUUID), nil
}
