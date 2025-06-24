package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) Create(ctx context.Context, cart *entity.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *CartRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Items").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Items").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) AddItem(ctx context.Context, item *entity.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CartRepository) Clear(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.CartItem{}).Error
}

func (r *CartRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete cart items first
		if err := tx.Where("cart_id = ?", id).Delete(&entity.CartItem{}).Error; err != nil {
			return err
		}
		// Delete cart
		return tx.Delete(&entity.Cart{}, "id = ?", id).Error
	})
}

func (r *CartRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get expired cart IDs
		var expiredCarts []entity.Cart
		if err := tx.Where("expires_at < ?", time.Now()).Find(&expiredCarts).Error; err != nil {
			return err
		}

		if len(expiredCarts) == 0 {
			return nil
		}

		// Extract cart IDs
		cartIDs := make([]uuid.UUID, len(expiredCarts))
		for i, cart := range expiredCarts {
			cartIDs[i] = cart.ID
		}

		// Delete cart items first
		if err := tx.Where("cart_id IN ?", cartIDs).Delete(&entity.CartItem{}).Error; err != nil {
			return err
		}

		// Delete carts
		return tx.Where("id IN ?", cartIDs).Delete(&entity.Cart{}).Error
	})
}

func (r *CartRepository) Update(ctx context.Context, cart *entity.Cart) error {
	return r.db.WithContext(ctx).Save(cart).Error
}
