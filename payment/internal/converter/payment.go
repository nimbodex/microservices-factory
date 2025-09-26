package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/payment/internal/model"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

// ToServicePaymentMethod converts protobuf PaymentMethod to service model
func ToServicePaymentMethod(protoMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch protoMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	default:
		return model.PaymentMethodUnknown
	}
}

// ToProtoPaymentMethod converts service model PaymentMethod to protobuf
func ToProtoPaymentMethod(serviceMethod model.PaymentMethod) paymentv1.PaymentMethod {
	switch serviceMethod {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

// ToServicePayOrderRequest converts protobuf request to service model
func ToServicePayOrderRequest(protoReq *paymentv1.PayOrderRequest) (*model.PayOrderRequest, error) {
	if protoReq == nil {
		return nil, fmt.Errorf("protoReq cannot be nil")
	}

	orderUUID, err := uuid.Parse(protoReq.OrderUuid)
	if err != nil {
		return nil, err
	}

	return &model.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: ToServicePaymentMethod(protoReq.PaymentMethod),
		Amount:        0.0, // Default amount, could be calculated from order
	}, nil
}

// ToProtoPayOrderResponse converts service model to protobuf response
func ToProtoPayOrderResponse(transactionUUID uuid.UUID) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}
}
