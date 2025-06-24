package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
)

// SagaRepository implements the repository.SagaRepository interface
type SagaRepository struct {
	db *gorm.DB
}

// NewSagaRepository creates a new PostgreSQL saga repository
func NewSagaRepository(db *gorm.DB) *SagaRepository {
	return &SagaRepository{
		db: db,
	}
}

// Create saves a new saga to the database
func (r *SagaRepository) Create(ctx context.Context, saga *entity.Saga) error {
	return r.db.WithContext(ctx).Create(saga).Error
}

// GetByID retrieves a saga by its ID
func (r *SagaRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Saga, error) {
	var saga entity.Saga
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		First(&saga, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &saga, nil
}

// Update updates an existing saga in the database
func (r *SagaRepository) Update(ctx context.Context, saga *entity.Saga) error {
	return r.db.WithContext(ctx).Save(saga).Error
}

// UpdateStepStatus updates the status of a saga step
func (r *SagaRepository) UpdateStepStatus(ctx context.Context, sagaID, stepID uuid.UUID, status entity.StepStatus, errorMessage string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update step status
		err := tx.Model(&entity.SagaStep{}).
			Where("id = ? AND saga_id = ?", stepID, sagaID).
			Updates(map[string]interface{}{
				"status":        status,
				"error_message": errorMessage,
			}).Error
		if err != nil {
			return err
		}

		// Get saga to update its status
		saga, err := r.GetByID(ctx, sagaID)
		if err != nil {
			return err
		}

		// Update saga status
		saga.UpdateStatus()
		return tx.Save(saga).Error
	})
}

// GetPendingSagas retrieves all pending sagas
func (r *SagaRepository) GetPendingSagas(ctx context.Context) ([]*entity.Saga, error) {
	var sagas []*entity.Saga
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Where("status = ?", entity.SagaStatusPending).
		Find(&sagas).Error
	return sagas, err
}

// GetFailedSagas retrieves all failed sagas
func (r *SagaRepository) GetFailedSagas(ctx context.Context) ([]*entity.Saga, error) {
	var sagas []*entity.Saga
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Where("status = ?", entity.SagaStatusFailed).
		Find(&sagas).Error
	return sagas, err
}

// GetCompensatingSagas retrieves all compensating sagas
func (r *SagaRepository) GetCompensatingSagas(ctx context.Context) ([]*entity.Saga, error) {
	var sagas []*entity.Saga
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\" ASC")
		}).
		Where("status = ?", entity.SagaStatusCompensating).
		Find(&sagas).Error
	return sagas, err
}

// Setup creates necessary indexes for the saga tables
func (r *SagaRepository) Setup(ctx context.Context) error {
	// Create indexes
	err := r.db.WithContext(ctx).Exec(`
		CREATE INDEX IF NOT EXISTS idx_sagas_status ON sagas(status);
		CREATE INDEX IF NOT EXISTS idx_saga_steps_saga_id ON saga_steps(saga_id);
		CREATE INDEX IF NOT EXISTS idx_saga_steps_status ON saga_steps(status);
	`).Error
	return err
}

func (r *SagaRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Saga, error) {
	var saga entity.Saga
	err := r.db.WithContext(ctx).
		Preload("Steps").
		Where("order_id = ?", orderID).
		First(&saga).Error
	if err != nil {
		return nil, err
	}
	return &saga, nil
}
