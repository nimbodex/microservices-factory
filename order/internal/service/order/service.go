package order

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"

	"github.com/nimbodex/microservices-factory/order/internal/client"
	"github.com/nimbodex/microservices-factory/order/internal/converter"
	"github.com/nimbodex/microservices-factory/order/internal/model"
	"github.com/nimbodex/microservices-factory/order/internal/repository"
)

// OrderServiceImpl implements OrderService interface
type OrderServiceImpl struct {
	orderRepo       repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
}

// NewOrderService creates a new order service instance
func NewOrderService(
	orderRepo repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

// CreateOrder creates a new order with the specified parts for a user
func (s *OrderServiceImpl) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	log.Printf("Creating order for user %s with parts %v", req.UserUUID, req.PartUuids)

	createReq := converter.ToCreateOrderRequest(req)

	if s.inventoryClient != nil {
		for _, partUUID := range createReq.PartUUIDs {
			_, err := s.inventoryClient.GetPart(ctx, partUUID)
			if err != nil {
				log.Printf("Part %s not found in inventory: %v", partUUID, err)
				return &orderv1.BadRequestError{
					Error:   "part_not_found",
					Message: fmt.Sprintf("part %s not found", partUUID),
				}, nil
			}
		}
	}

	orderUUID := uuid.New()
	order := &model.Order{
		UUID:      orderUUID,
		UserUUID:  createReq.UserUUID,
		PartUUIDs: createReq.PartUUIDs,
		Status:    model.StatusPendingPayment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		log.Printf("Failed to create order: %v", err)
		return &orderv1.InternalServerError{
			Error:   "creation_failed",
			Message: "failed to create order",
		}, nil
	}

	log.Printf("Order %s created successfully", orderUUID)

	totalPrice := 0.0
	return converter.ToCreateOrderResponse(order, totalPrice), nil
}

// GetOrder retrieves an order by its UUID
func (s *OrderServiceImpl) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	log.Printf("Getting order %s", params.OrderUUID)

	order, err := s.orderRepo.GetByUUID(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Order %s not found: %v", params.OrderUUID, err)
		return &orderv1.NotFoundError{
			Error:   "order_not_found",
			Message: "order not found",
		}, nil
	}

	log.Printf("Order %s found with status %s", params.OrderUUID, order.Status)

	totalPrice := 0.0
	return converter.ToGetOrderResponse(order, totalPrice), nil
}

// PayOrder processes payment for an order using the specified payment method
func (s *OrderServiceImpl) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	log.Printf("Processing payment for order %s with method %s", params.OrderUUID, req.PaymentMethod)

	order, err := s.orderRepo.GetByUUID(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Order %s not found: %v", params.OrderUUID, err)
		return &orderv1.NotFoundError{
			Error:   "order_not_found",
			Message: "order not found",
		}, nil
	}

	if order.Status != model.StatusPendingPayment {
		log.Printf("Order %s cannot be paid, current status: %s", params.OrderUUID, order.Status)
		return &orderv1.ConflictError{
			Error:   "invalid_status",
			Message: "order cannot be paid",
		}, nil
	}

	var transactionUUID uuid.UUID

	if s.paymentClient != nil {
		payReq := converter.ToPayOrderRequest(req)

		paymentResult, err := s.paymentClient.PayOrder(ctx, params.OrderUUID, client.PaymentMethod(payReq.PaymentMethod))
		if err != nil {
			log.Printf("Payment failed for order %s: %v", params.OrderUUID, err)
			return &orderv1.InternalServerError{
				Error:   "payment_failed",
				Message: "payment processing failed",
			}, nil
		}

		transactionUUID = paymentResult.TransactionUUID
	} else {
		transactionUUID = uuid.New()
	}

	order.Status = model.StatusPaid
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		log.Printf("Failed to update order %s: %v", params.OrderUUID, err)
		return &orderv1.InternalServerError{
			Error:   "update_failed",
			Message: "failed to update order status",
		}, nil
	}

	log.Printf("Payment successful for order %s, transaction: %s", params.OrderUUID, transactionUUID)

	return converter.ToPayOrderResponse(transactionUUID), nil
}

// CancelOrder cancels an order if it is in pending payment status
func (s *OrderServiceImpl) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	log.Printf("Cancelling order %s", params.OrderUUID)

	order, err := s.orderRepo.GetByUUID(ctx, params.OrderUUID)
	if err != nil {
		log.Printf("Order %s not found: %v", params.OrderUUID, err)
		return &orderv1.NotFoundError{
			Error:   "order_not_found",
			Message: "order not found",
		}, nil
	}

	if order.Status != model.StatusPendingPayment {
		log.Printf("Order %s cannot be cancelled, current status: %s", params.OrderUUID, order.Status)
		return &orderv1.ConflictError{
			Error:   "invalid_status",
			Message: "order cannot be cancelled",
		}, nil
	}

	order.Status = model.StatusCancelled
	order.UpdatedAt = time.Now()

	if err := s.orderRepo.Update(ctx, order); err != nil {
		log.Printf("Failed to update order %s: %v", params.OrderUUID, err)
		return &orderv1.InternalServerError{
			Error:   "update_failed",
			Message: "failed to update order status",
		}, nil
	}

	log.Printf("Order %s cancelled successfully", params.OrderUUID)

	return &orderv1.CancelOrderNoContent{}, nil
}

// NewError creates a standardized internal server error response
func (s *OrderServiceImpl) NewError(ctx context.Context, err error) *orderv1.InternalServerErrorStatusCode {
	log.Printf("Internal error: %v", err)
	return &orderv1.InternalServerErrorStatusCode{
		StatusCode: 500,
		Response: orderv1.InternalServerError{
			Error:   "internal_error",
			Message: err.Error(),
		},
	}
}
