package v1

import (
	"context"

	"github.com/nimbodex/microservices-factory/payment/internal/service"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

// APIHandler handles gRPC requests for payment API
type APIHandler struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(paymentService service.PaymentService) *APIHandler {
	return &APIHandler{
		paymentService: paymentService,
	}
}

func (h *APIHandler) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	return h.paymentService.PayOrder(ctx, req)
}
