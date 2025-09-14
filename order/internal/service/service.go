package service

import (
	"context"

	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

// OrderService defines the interface for order service operations
type OrderService interface {
	CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error)
	GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error)
	PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error)
	CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error)
	NewError(ctx context.Context, err error) *orderv1.InternalServerErrorStatusCode
}
