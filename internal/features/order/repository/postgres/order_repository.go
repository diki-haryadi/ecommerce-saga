package postgres

import (
	"context"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
)

// OrderRepository implements the domain.OrderRepository interface
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new PostgreSQL order postgres
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Create saves a new order to the database
func (r *OrderRepository) Create(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves orders for a user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Update updates an existing order in the database
func (r *OrderRepository) Update(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// UpdateStatus updates the order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete removes an order from the database
func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Order{}, "id = ?", id).Error
}

// CountByUserID counts total orders for a user
func (r *OrderRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// Setup creates necessary indexes for the order tables
func (r *OrderRepository) Setup(ctx context.Context) error {
	// Create indexes
	err := r.db.WithContext(ctx).Exec(`
		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
		CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
	`).Error
	return err
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Order, error) {
	var orders []*entity.Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
