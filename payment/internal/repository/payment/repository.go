package payment

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/payment/internal/model"
)

// MemoryPaymentRepository implements PaymentRepository using in-memory storage
type MemoryPaymentRepository struct {
	mu       sync.RWMutex
	payments map[string]*model.Payment
}

// NewMemoryPaymentRepository creates a new in-memory payment repository
func NewMemoryPaymentRepository() *MemoryPaymentRepository {
	return &MemoryPaymentRepository{
		payments: make(map[string]*model.Payment),
	}
}

// Create creates a new payment
func (r *MemoryPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	paymentKey := payment.UUID.String()
	if _, exists := r.payments[paymentKey]; exists {
		return fmt.Errorf("payment with UUID %s already exists", payment.UUID)
	}

	// Create a copy to avoid external modifications
	paymentCopy := *payment
	r.payments[paymentKey] = &paymentCopy

	return nil
}

// GetByUUID retrieves a payment by its UUID
func (r *MemoryPaymentRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	payment, exists := r.payments[uuid.String()]
	if !exists {
		return nil, fmt.Errorf("payment with UUID %s not found", uuid)
	}

	// Return a copy to avoid external modifications
	paymentCopy := *payment
	return &paymentCopy, nil
}

// GetByOrderUUID retrieves a payment by order UUID
func (r *MemoryPaymentRepository) GetByOrderUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, payment := range r.payments {
		if payment.OrderUUID == orderUUID {
			// Return a copy to avoid external modifications
			paymentCopy := *payment
			return &paymentCopy, nil
		}
	}

	return nil, fmt.Errorf("payment for order %s not found", orderUUID)
}

// GetByTransactionUUID retrieves a payment by transaction UUID
func (r *MemoryPaymentRepository) GetByTransactionUUID(ctx context.Context, transactionUUID uuid.UUID) (*model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, payment := range r.payments {
		if payment.TransactionUUID == transactionUUID {
			// Return a copy to avoid external modifications
			paymentCopy := *payment
			return &paymentCopy, nil
		}
	}

	return nil, fmt.Errorf("payment with transaction UUID %s not found", transactionUUID)
}

// Update updates an existing payment
func (r *MemoryPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	paymentKey := payment.UUID.String()
	if _, exists := r.payments[paymentKey]; !exists {
		return fmt.Errorf("payment with UUID %s not found", payment.UUID)
	}

	// Create a copy to avoid external modifications
	paymentCopy := *payment
	r.payments[paymentKey] = &paymentCopy

	return nil
}

// Delete removes a payment by its UUID
func (r *MemoryPaymentRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	paymentKey := uuid.String()
	if _, exists := r.payments[paymentKey]; !exists {
		return fmt.Errorf("payment with UUID %s not found", uuid)
	}

	delete(r.payments, paymentKey)
	return nil
}
