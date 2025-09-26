package converter

import (
	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/order/internal/model"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
)

// ToCreateOrderRequest converts OpenAPI request to service model
func ToCreateOrderRequest(req *orderv1.CreateOrderRequest) *model.CreateOrderRequest {
	partUUIDs := make([]uuid.UUID, len(req.PartUuids))
	copy(partUUIDs, req.PartUuids)

	return &model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUUIDs: partUUIDs,
	}
}

// ToCreateOrderResponse converts service model to OpenAPI response
func ToCreateOrderResponse(order *model.Order, totalPrice float64) *orderv1.CreateOrderResponse {
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: totalPrice,
	}
}

// ToGetOrderResponse converts service model to OpenAPI response
func ToGetOrderResponse(order *model.Order, totalPrice float64) *orderv1.GetOrderResponse {
	return &orderv1.GetOrderResponse{
		OrderUUID:  order.UUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: totalPrice,
		Status:     orderv1.OrderStatus(order.Status),
	}
}

// ToPayOrderRequest converts OpenAPI request to service model
func ToPayOrderRequest(req *orderv1.PayOrderRequest) *model.PayOrderRequest {
	var paymentMethod model.PaymentMethod
	switch req.PaymentMethod {
	case orderv1.PaymentMethodCARD:
		paymentMethod = model.PaymentMethodCard
	case orderv1.PaymentMethodSBP:
		paymentMethod = model.PaymentMethodSBP
	default:
		paymentMethod = model.PaymentMethodUnknown
	}

	return &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}
}

// ToPayOrderResponse converts service model to OpenAPI response
func ToPayOrderResponse(transactionUUID uuid.UUID) *orderv1.PayOrderResponse {
	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}
}

// ToOrderStatus converts service OrderStatus to OpenAPI OrderStatus
func ToOrderStatus(status model.OrderStatus) orderv1.OrderStatus {
	return orderv1.OrderStatus(status)
}
