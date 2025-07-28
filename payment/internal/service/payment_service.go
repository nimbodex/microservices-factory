package service

import (
	"context"
	"log"

	"github.com/google/uuid"

	paymentv1 "github.com/nexarise/microservices-factory/shared/pkg/proto/payment/v1"
)

type PaymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionUUID := uuid.New().String()

	log.Printf("Payment was successful, transaction_uuid: %s", transactionUUID)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
