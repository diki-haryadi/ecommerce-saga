package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
)

// CartRepository defines the interface for cart data persistence
type CartRepository interface {
	// Create saves a new cart to the database
	Create(ctx context.Context, cart *entity.Cart) error

	// GetByID retrieves a cart by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Cart, error)

	// GetByUserID retrieves a cart by user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error)

	// Update updates an existing cart in the database
	Update(ctx context.Context, cart *entity.Cart) error

	// Delete removes a cart from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteExpired removes all expired carts from the database
	DeleteExpired(ctx context.Context) error
}
