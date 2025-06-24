package postgres

import (
	"context"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
)

// PaymentRepository implements the repository.PaymentRepository interface
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PostgreSQL payment postgres
func NewPaymentRepository(db *gorm.DB) usecase.Repository {
	return &PaymentRepository{
		db: db,
	}
}

// Create creates a new payment record
func (r *PaymentRepository) Create(ctx context.Context, payment *usecase.PaymentResponse) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*usecase.PaymentResponse, error) {
	var p usecase.PaymentResponse
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// List retrieves a list of payments with pagination
func (r *PaymentRepository) List(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*usecase.PaymentResponse, int64, error) {
	var payments []*usecase.PaymentResponse
	var total int64

	query := r.db.WithContext(ctx).Model(&usecase.PaymentResponse{})

	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(int(offset)).Limit(int(limit)).Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// Update updates a payment record
func (r *PaymentRepository) Update(ctx context.Context, payment *usecase.PaymentResponse) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// GetByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*usecase.PaymentResponse, error) {
	var p usecase.PaymentResponse
	if err := r.db.WithContext(ctx).First(&p, "order_id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// UpdateStatus updates the payment status
func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.PaymentStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Payment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete removes a payment from the database
func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Payment{}, "id = ?", id).Error
}

// Setup creates necessary indexes for the payment table
func (r *PaymentRepository) Setup(ctx context.Context) error {
	// Create indexes
	err := r.db.WithContext(ctx).Exec(`
		CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
		CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
		CREATE INDEX IF NOT EXISTS idx_payments_provider ON payments(provider);
	`).Error
	return err
}
