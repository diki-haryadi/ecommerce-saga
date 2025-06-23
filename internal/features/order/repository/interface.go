package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
)

// OrderRepository defines the interface for order data persistence
type OrderRepository interface {
	// Create saves a new order to the database
	Create(ctx context.Context, order *entity.Order) error

	// GetByID retrieves an order by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)

	// GetByUserID retrieves orders for a user
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Order, error)

	// Update updates an existing order in the database
	Update(ctx context.Context, order *entity.Order) error

	// UpdateStatus updates the order status
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.OrderStatus) error

	// Delete removes an order from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByUserID counts total orders for a user
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}
