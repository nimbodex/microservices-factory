package model

import (
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents payment method
type PaymentMethod string

const (
	PaymentMethodUnknown PaymentMethod = "UNKNOWN"
	PaymentMethodCard    PaymentMethod = "CARD"
	PaymentMethodSBP     PaymentMethod = "SBP"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

// Payment represents a payment in the service layer
type Payment struct {
	UUID            uuid.UUID     `json:"uuid"`
	OrderUUID       uuid.UUID     `json:"order_uuid"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	Amount          float64       `json:"amount"`
	Status          PaymentStatus `json:"status"`
	TransactionUUID uuid.UUID     `json:"transaction_uuid"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// PayOrderRequest represents request to pay an order
type PayOrderRequest struct {
	OrderUUID     uuid.UUID     `json:"order_uuid"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	Amount        float64       `json:"amount"`
}
