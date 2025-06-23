package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/auth/domain/entity"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	// Create saves a new user to the database
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates an existing user in the database
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateRefreshToken updates the refresh token for a user
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error

	// GetByRefreshToken retrieves a user by their refresh token
	GetByRefreshToken(ctx context.Context, token string) (*entity.User, error)
}
