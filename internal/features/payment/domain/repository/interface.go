package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
)

// PaymentRepository defines the interface for payment data persistence
type PaymentRepository interface {
	// Create saves a new payment to the database
	Create(ctx context.Context, payment *entity.Payment) error

	// GetByID retrieves a payment by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)

	// GetByOrderID retrieves a payment by order ID
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error)

	// Update updates an existing payment in the database
	Update(ctx context.Context, payment *entity.Payment) error

	// UpdateStatus updates the payment status
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error

	// Delete removes a payment from the database
	Delete(ctx context.Context, id uuid.UUID) error
}
