package service

import (
	"context"

	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

// PaymentService defines the interface for payment service operations
type PaymentService interface {
	PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error)
}
