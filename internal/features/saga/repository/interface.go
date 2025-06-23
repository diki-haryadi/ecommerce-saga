package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
)

// SagaRepository defines the interface for saga data persistence
type SagaRepository interface {
	// Create saves a new saga to the database
	Create(ctx context.Context, saga *entity.Saga) error

	// GetByID retrieves a saga by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Saga, error)

	// Update updates an existing saga in the database
	Update(ctx context.Context, saga *entity.Saga) error

	// UpdateStepStatus updates the status of a saga step
	UpdateStepStatus(ctx context.Context, sagaID, stepID uuid.UUID, status entity.StepStatus, errorMessage string) error

	// GetPendingSagas retrieves all pending sagas
	GetPendingSagas(ctx context.Context) ([]*entity.Saga, error)

	// GetFailedSagas retrieves all failed sagas
	GetFailedSagas(ctx context.Context) ([]*entity.Saga, error)

	// GetCompensatingSagas retrieves all compensating sagas
	GetCompensatingSagas(ctx context.Context) ([]*entity.Saga, error)
}
