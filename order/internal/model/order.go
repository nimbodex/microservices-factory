package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	StatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	StatusPaid           OrderStatus = "PAID"
	StatusAssembled      OrderStatus = "ASSEMBLED"
	StatusCancelled      OrderStatus = "CANCELLED"
)

// Order represents an order in the service layer
type Order struct {
	UUID      uuid.UUID   `json:"uuid"`
	UserUUID  uuid.UUID   `json:"user_uuid"`
	PartUUIDs []uuid.UUID `json:"part_uuids"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Part represents a part in the service layer
type Part struct {
	UUID  uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

// PaymentMethod represents payment method
type PaymentMethod string

const (
	PaymentMethodCard    PaymentMethod = "CARD"
	PaymentMethodSBP     PaymentMethod = "SBP"
	PaymentMethodUnknown PaymentMethod = "UNKNOWN"
)

// CreateOrderRequest represents request to create an order
type CreateOrderRequest struct {
	UserUUID  uuid.UUID   `json:"user_uuid"`
	PartUUIDs []uuid.UUID `json:"part_uuids"`
}

// PayOrderRequest represents request to pay an order
type PayOrderRequest struct {
	PaymentMethod PaymentMethod `json:"payment_method"`
}
