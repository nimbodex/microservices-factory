package order

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/nimbodex/microservices-factory/order/internal/model"
)

// PostgresOrderRepository implements OrderRepository using PostgreSQL
type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a new PostgreSQL order repository
func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

// Create creates a new order in the repository
func (r *PostgresOrderRepository) Create(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (uuid, user_uuid, part_uuids, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	partUUIDStrings := make([]string, len(order.PartUUIDs))
	for i, partUUID := range order.PartUUIDs {
		partUUIDStrings[i] = partUUID.String()
	}

	_, err := r.db.ExecContext(ctx, query,
		order.UUID,
		order.UserUUID,
		pq.Array(partUUIDStrings),
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetByUUID retrieves an order by its UUID
func (r *PostgresOrderRepository) GetByUUID(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	query := `
		SELECT uuid, user_uuid, part_uuids, status, created_at, updated_at
		FROM orders
		WHERE uuid = $1
	`

	var order model.Order
	var partUUIDStrings []string

	err := r.db.QueryRowContext(ctx, query, orderUUID).Scan(
		&order.UUID,
		&order.UserUUID,
		pq.Array(&partUUIDStrings),
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order with UUID %s not found", orderUUID)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order.PartUUIDs = make([]uuid.UUID, len(partUUIDStrings))
	for i, partUUIDString := range partUUIDStrings {
		parsedUUID, parseErr := uuid.Parse(partUUIDString)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse part UUID %s: %w", partUUIDString, parseErr)
		}
		order.PartUUIDs[i] = parsedUUID
	}

	return &order, nil
}

// Update updates an existing order
func (r *PostgresOrderRepository) Update(ctx context.Context, order *model.Order) error {
	query := `
		UPDATE orders
		SET user_uuid = $2, part_uuids = $3, status = $4, updated_at = $5
		WHERE uuid = $1
	`

	partUUIDStrings := make([]string, len(order.PartUUIDs))
	for i, partUUID := range order.PartUUIDs {
		partUUIDStrings[i] = partUUID.String()
	}

	result, err := r.db.ExecContext(ctx, query,
		order.UUID,
		order.UserUUID,
		pq.Array(partUUIDStrings),
		order.Status,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order with UUID %s not found", order.UUID)
	}

	return nil
}

// Delete removes an order by its UUID
func (r *PostgresOrderRepository) Delete(ctx context.Context, orderUUID uuid.UUID) error {
	query := `DELETE FROM orders WHERE uuid = $1`

	result, err := r.db.ExecContext(ctx, query, orderUUID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order with UUID %s not found", orderUUID)
	}

	return nil
}

// List retrieves orders with pagination
func (r *PostgresOrderRepository) List(ctx context.Context, limit, offset int) ([]*model.Order, error) {
	query := `
		SELECT uuid, user_uuid, part_uuids, status, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("failed to close rows: %v", closeErr)
		}
	}()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		var partUUIDStrings []string

		scanErr := rows.Scan(
			&order.UUID,
			&order.UserUUID,
			pq.Array(&partUUIDStrings),
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("failed to scan order: %w", scanErr)
		}

		// Convert []string back to []uuid.UUID
		order.PartUUIDs = make([]uuid.UUID, len(partUUIDStrings))
		for i, partUUIDString := range partUUIDStrings {
			parsedUUID, parseErr := uuid.Parse(partUUIDString)
			if parseErr != nil {
				return nil, fmt.Errorf("failed to parse part UUID %s: %w", partUUIDString, parseErr)
			}
			order.PartUUIDs[i] = parsedUUID
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over orders: %w", err)
	}

	return orders, nil
}

// Close closes the database connection
func (r *PostgresOrderRepository) Close() error {
	return r.db.Close()
}
