package postgres

import (
	"context"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentRepository implements the repository.PaymentRepository interface
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new PostgreSQL payment repository
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

// Create saves a new payment to the database
func (r *PaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetByID retrieves a payment by its ID
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.WithContext(ctx).First(&payment, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update updates an existing payment in the database
func (r *PaymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
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
