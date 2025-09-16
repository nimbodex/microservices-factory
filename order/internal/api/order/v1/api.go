package v1

import (
	"context"

	"github.com/nimbodex/microservices-factory/order/internal/service"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

// APIHandler handles HTTP requests for order API
type APIHandler struct {
	orderService service.OrderService
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(orderService service.OrderService) *APIHandler {
	return &APIHandler{
		orderService: orderService,
	}
}

// CreateOrder handles POST /orders requests
func (h *APIHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	return h.orderService.CreateOrder(ctx, req)
}

// GetOrder handles GET /orders/{order_uuid} requests
func (h *APIHandler) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	return h.orderService.GetOrder(ctx, params)
}

// PayOrder handles POST /orders/{order_uuid}/pay requests
func (h *APIHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	return h.orderService.PayOrder(ctx, req, params)
}

// CancelOrder handles POST /orders/{order_uuid}/cancel requests
func (h *APIHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	return h.orderService.CancelOrder(ctx, params)
}

// NewError handles internal server errors
func (h *APIHandler) NewError(ctx context.Context, err error) *orderv1.InternalServerErrorStatusCode {
	return h.orderService.NewError(ctx, err)
}
