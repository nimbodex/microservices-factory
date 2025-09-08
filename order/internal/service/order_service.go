package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderv1 "github.com/nimbodex/microservices-factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/nimbodex/microservices-factory/shared/pkg/proto/payment/v1"
)

type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	UUID      string      `json:"uuid"`
	UserUUID  string      `json:"user_uuid"`
	PartUUIDs []string    `json:"part_uuids"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderService struct {
	mu              sync.RWMutex
	orders          map[string]*Order
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	inventoryConn   *grpc.ClientConn
	paymentConn     *grpc.ClientConn
}

// NewOrderService creates a new instance of OrderService and connects to external services.
func NewOrderService() *OrderService {
	service := &OrderService{
		orders: make(map[string]*Order),
	}

	service.connectToServices()

	return service
}

func (s *OrderService) connectToServices() {
	inventoryConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to inventory service: %v", err)
	} else {
		s.inventoryConn = inventoryConn
		s.inventoryClient = inventoryv1.NewInventoryServiceClient(inventoryConn)
		log.Println("Connected to inventory service")
	}

	paymentConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to payment service: %v", err)
	} else {
		s.paymentConn = paymentConn
		s.paymentClient = paymentv1.NewPaymentServiceClient(paymentConn)
		log.Println("Connected to payment service")
	}
}

// CreateOrder creates a new order with the specified parts for a user.
func (s *OrderService) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	log.Printf("Creating order for user %s with parts %v", req.UserUUID, req.PartUuids)

	if s.inventoryClient != nil {
		for _, partUUID := range req.PartUuids {
			_, err := s.inventoryClient.GetPart(ctx, &inventoryv1.GetPartRequest{
				Uuid: partUUID.String(),
			})
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

	partUUIDs := make([]string, len(req.PartUuids))
	for i, uuid := range req.PartUuids {
		partUUIDs[i] = uuid.String()
	}

	order := &Order{
		UUID:      orderUUID.String(),
		UserUUID:  req.UserUUID.String(),
		PartUUIDs: partUUIDs,
		Status:    StatusPendingPayment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.mu.Lock()
	s.orders[orderUUID.String()] = order
	s.mu.Unlock()

	log.Printf("Order %s created successfully", orderUUID)

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: 0.0,
	}, nil
}

// GetOrder retrieves an order by its UUID.
func (s *OrderService) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	log.Printf("Getting order %s", params.OrderUUID)

	s.mu.RLock()
	order, exists := s.orders[params.OrderUUID.String()]
	s.mu.RUnlock()

	if !exists {
		log.Printf("Order %s not found", params.OrderUUID)
		return &orderv1.NotFoundError{
			Error:   "order_not_found",
			Message: "order not found",
		}, nil
	}

	log.Printf("Order %s found with status %s", params.OrderUUID, order.Status)

	partUUIDs := make([]uuid.UUID, len(order.PartUUIDs))
	for i, uuidStr := range order.PartUUIDs {
		parsedUUID, err := uuid.Parse(uuidStr)
		if err != nil {
			log.Printf("Failed to parse part UUID %s: %v", uuidStr, err)
			return &orderv1.InternalServerError{
				Error:   "invalid_uuid",
				Message: fmt.Sprintf("invalid part UUID: %s", uuidStr),
			}, nil
		}
		partUUIDs[i] = parsedUUID
	}

	orderUUID, err := uuid.Parse(order.UUID)
	if err != nil {
		log.Printf("Failed to parse order UUID %s: %v", order.UUID, err)
		return &orderv1.InternalServerError{
			Error:   "invalid_uuid",
			Message: fmt.Sprintf("invalid order UUID: %s", order.UUID),
		}, nil
	}
	userUUID, err := uuid.Parse(order.UserUUID)
	if err != nil {
		log.Printf("Failed to parse user UUID %s: %v", order.UserUUID, err)
		return &orderv1.InternalServerError{
			Error:   "invalid_uuid",
			Message: fmt.Sprintf("invalid user UUID: %s", order.UserUUID),
		}, nil
	}

	return &orderv1.GetOrderResponse{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: 0.0,
		Status:     orderv1.OrderStatus(order.Status),
	}, nil
}

// PayOrder processes payment for an order using the specified payment method.
func (s *OrderService) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	log.Printf("Processing payment for order %s with method %s", params.OrderUUID, req.PaymentMethod)

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[params.OrderUUID.String()]
	if !exists {
		log.Printf("Order %s not found", params.OrderUUID)
		return &orderv1.NotFoundError{
			Error:   "order_not_found",
			Message: "order not found",
		}, nil
	}

	if order.Status != StatusPendingPayment {
		log.Printf("Order %s cannot be paid, current status: %s", params.OrderUUID, order.Status)
		return &orderv1.ConflictError{
			Error:   "invalid_status",
			Message: "order cannot be paid",
		}, nil
	}

	var transactionUUID uuid.UUID
	if s.paymentClient != nil {
		var paymentMethod paymentv1.PaymentMethod
		switch req.PaymentMethod {
		case orderv1.PaymentMethodCARD:
			paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
		case orderv1.PaymentMethodSBP:
			paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
		default:
			paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
		}

		paymentResp, err := s.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
			OrderUuid:     params.OrderUUID.String(),
			PaymentMethod: paymentMethod,
		})
		if err != nil {
			log.Printf("Payment failed for order %s: %v", params.OrderUUID, err)
			return &orderv1.InternalServerError{
				Error:   "payment_failed",
				Message: "payment processing failed",
			}, nil
		}
		parsedUUID, err := uuid.Parse(paymentResp.TransactionUuid)
		if err != nil {
			log.Printf("Failed to parse transaction UUID %s: %v", paymentResp.TransactionUuid, err)
			return &orderv1.InternalServerError{
				Error:   "invalid_uuid",
				Message: fmt.Sprintf("invalid transaction UUID: %s", paymentResp.TransactionUuid),
			}, nil
		}
		transactionUUID = parsedUUID
	} else {
		transactionUUID = uuid.New()
	}

	order.Status = StatusPaid
	order.UpdatedAt = time.Now()

	log.Printf("Payment successful for order %s, transaction: %s", params.OrderUUID, transactionUUID)

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// CancelOrder cancels an order if it is in pending payment status.
func (s *OrderService) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	log.Printf("Cancelling order %s", params.OrderUUID)

	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[params.OrderUUID.String()]
	if !exists {
		log.Printf("Order %s not found", params.OrderUUID)
		return &orderv1.NotFoundError{
			Error:   "order_not_found",
			Message: "order not found",
		}, nil
	}

	if order.Status != StatusPendingPayment {
		log.Printf("Order %s cannot be cancelled, current status: %s", params.OrderUUID, order.Status)
		return &orderv1.ConflictError{
			Error:   "invalid_status",
			Message: "order cannot be cancelled",
		}, nil
	}

	order.Status = StatusCancelled
	order.UpdatedAt = time.Now()

	log.Printf("Order %s cancelled successfully", params.OrderUUID)

	return &orderv1.CancelOrderNoContent{}, nil
}

// NewError creates a standardized internal server error response.
func (s *OrderService) NewError(ctx context.Context, err error) *orderv1.InternalServerErrorStatusCode {
	log.Printf("Internal error: %v", err)
	return &orderv1.InternalServerErrorStatusCode{
		StatusCode: 500,
		Response: orderv1.InternalServerError{
			Error:   "internal_error",
			Message: err.Error(),
		},
	}
}

// Close closes all external service connections.
func (s *OrderService) Close() {
	if s.inventoryConn != nil {
		if err := s.inventoryConn.Close(); err != nil {
			log.Printf("Failed to close inventory connection: %v", err)
		}
	}
	if s.paymentConn != nil {
		if err := s.paymentConn.Close(); err != nil {
			log.Printf("Failed to close payment connection: %v", err)
		}
	}
}
