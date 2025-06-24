package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository"
)

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new PostgreSQL cart postgres
func NewCartRepository(db *gorm.DB) repository.CartRepository {
	return &cartRepository{
		db: db,
	}
}

// Create saves a new cart to the database
func (r *cartRepository) Create(ctx context.Context, cart *entity.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

// GetByID retrieves a cart by its ID
func (r *cartRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&cart, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetByUserID retrieves a cart by user ID
func (r *cartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&cart, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// AddItem adds a new item to the cart
func (r *cartRepository) AddItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// Clear removes all items from the cart
func (r *cartRepository) Clear(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.CartItem{}).Error
}

// Update updates an existing cart in the database
func (r *cartRepository) Update(ctx context.Context, cart *entity.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}

// Delete removes a cart from the database
func (r *cartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete cart items first
		if err := tx.Delete(&entity.CartItem{}, "cart_id = ?", id).Error; err != nil {
			return err
		}
		// Then delete the cart
		return tx.Delete(&entity.Cart{}, "id = ?", id).Error
	})
}

// DeleteExpired removes all expired carts from the database
func (r *cartRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get expired cart IDs
		var expiredCartIDs []uuid.UUID
		if err := tx.Model(&entity.Cart{}).
			Where("expires_at < NOW()").
			Pluck("id", &expiredCartIDs).Error; err != nil {
			return err
		}

		if len(expiredCartIDs) == 0 {
			return nil
		}

		// Delete cart items first
		if err := tx.Delete(&entity.CartItem{}, "cart_id IN ?", expiredCartIDs).Error; err != nil {
			return err
		}

		// Then delete the carts
		return tx.Delete(&entity.Cart{}, "id IN ?", expiredCartIDs).Error
	})
}
