package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
)

type PostgresSagaRepository struct {
	db *gorm.DB
}

func NewPostgresSagaRepository(db *gorm.DB) *PostgresSagaRepository {
	return &PostgresSagaRepository{
		db: db,
	}
}

// DB returns the underlying database connection
func (r *PostgresSagaRepository) DB() *gorm.DB {
	return r.db
}

func (r *PostgresSagaRepository) Create(ctx context.Context, saga *entity.SagaTransaction) error {
	return r.db.WithContext(ctx).Create(saga).Error
}

func (r *PostgresSagaRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SagaTransaction, error) {
	var saga entity.SagaTransaction
	if err := r.db.WithContext(ctx).First(&saga, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get saga: %w", err)
	}
	return &saga, nil
}

func (r *PostgresSagaRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.SagaTransaction, error) {
	var saga entity.SagaTransaction
	if err := r.db.WithContext(ctx).First(&saga, "order_id = ?", orderID).Error; err != nil {
		return nil, fmt.Errorf("failed to get saga: %w", err)
	}
	return &saga, nil
}

func (r *PostgresSagaRepository) Update(ctx context.Context, saga *entity.SagaTransaction) error {
	return r.db.WithContext(ctx).Save(saga).Error
}

func (r *PostgresSagaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SagaTransaction{}, "id = ?", id).Error
}

func (r *PostgresSagaRepository) GetPendingSagas(ctx context.Context) ([]*entity.SagaTransaction, error) {
	var sagas []*entity.SagaTransaction
	err := r.db.WithContext(ctx).
		Where("status IN ?", []entity.SagaStatus{entity.SagaStatusPending, entity.SagaStatusProcessing}).
		Find(&sagas).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending sagas: %w", err)
	}
	return sagas, nil
}

// Migrations returns the database migrations for the saga tables
func (r *PostgresSagaRepository) Migrations() []string {
	return []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE TYPE saga_status AS ENUM ('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED', 'COMPENSATING');`,
		`CREATE TYPE step_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED', 'CANCELLED', 'COMPLETED', 'COMPENSATED');`,
		`CREATE TYPE step_type AS ENUM ('CREATE_ORDER', 'PROCESS_PAYMENT', 'UPDATE_INVENTORY');`,
		`
		CREATE TABLE IF NOT EXISTS saga_transactions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			order_id UUID NOT NULL,
			status saga_status NOT NULL,
			current_step JSONB NOT NULL,
			steps JSONB NOT NULL,
			compensation_steps JSONB NOT NULL,
			error TEXT,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
			timeout INTERVAL NOT NULL,
			max_retries INTEGER NOT NULL
		);
		`,
		`CREATE INDEX IF NOT EXISTS idx_saga_transactions_order_id ON saga_transactions(order_id);`,
		`CREATE INDEX IF NOT EXISTS idx_saga_transactions_status ON saga_transactions(status);`,
	}
}
