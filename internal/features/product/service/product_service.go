package service

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/usecase"
)

type productService struct {
	db *gorm.DB
}

// NewProductService creates a new product service
func NewProductService(db *gorm.DB) usecase.ProductService {
	return &productService{
		db: db,
	}
}

// GetProduct retrieves a product by ID
func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*usecase.Product, error) {
	var product struct {
		ID    uuid.UUID `gorm:"column:id"`
		Name  string    `gorm:"column:name"`
		Price float64   `gorm:"column:price"`
		Stock int       `gorm:"column:stock"`
	}

	err := s.db.WithContext(ctx).
		Table("products").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return &usecase.Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}, nil
}
