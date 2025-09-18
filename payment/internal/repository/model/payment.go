package model

import (
	"time"
)

// Payment represents a payment in the repository layer
type Payment struct {
	UUID            string    `json:"uuid"`
	OrderUUID       string    `json:"order_uuid"`
	PaymentMethod   string    `json:"payment_method"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	TransactionUUID string    `json:"transaction_uuid"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
