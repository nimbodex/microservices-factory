package order

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/nimbodex/microservices-factory/order/internal/model"
)

// MemoryOrderRepository implements OrderRepository using in-memory storage
type MemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

// NewMemoryOrderRepository creates a new in-memory order repository
func NewMemoryOrderRepository() *MemoryOrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*model.Order),
	}
}

// Create creates a new order in the repository
func (r *MemoryOrderRepository) Create(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderKey := order.UUID.String()
	if _, exists := r.orders[orderKey]; exists {
		return fmt.Errorf("order with UUID %s already exists", order.UUID)
	}

	// Create a copy to avoid external modifications
	orderCopy := *order
	r.orders[orderKey] = &orderCopy

	return nil
}

// GetByUUID retrieves an order by its UUID
func (r *MemoryOrderRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid.String()]
	if !exists {
		return nil, fmt.Errorf("order with UUID %s not found", uuid)
	}

	// Return a copy to avoid external modifications
	orderCopy := *order
	return &orderCopy, nil
}

// Update updates an existing order
func (r *MemoryOrderRepository) Update(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderKey := order.UUID.String()
	if _, exists := r.orders[orderKey]; !exists {
		return fmt.Errorf("order with UUID %s not found", order.UUID)
	}

	// Create a copy to avoid external modifications
	orderCopy := *order
	r.orders[orderKey] = &orderCopy

	return nil
}

// Delete removes an order by its UUID
func (r *MemoryOrderRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderKey := uuid.String()
	if _, exists := r.orders[orderKey]; !exists {
		return fmt.Errorf("order with UUID %s not found", uuid)
	}

	delete(r.orders, orderKey)
	return nil
}

// List retrieves orders with pagination
func (r *MemoryOrderRepository) List(ctx context.Context, limit, offset int) ([]*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*model.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	// Simple pagination (in real implementation, you'd want to sort by created_at desc)
	start := offset
	end := offset + limit

	if start >= len(orders) {
		return []*model.Order{}, nil
	}

	if end > len(orders) {
		end = len(orders)
	}

	// Return copies to avoid external modifications
	result := make([]*model.Order, end-start)
	for i, order := range orders[start:end] {
		orderCopy := *order
		result[i] = &orderCopy
	}

	return result, nil
}
