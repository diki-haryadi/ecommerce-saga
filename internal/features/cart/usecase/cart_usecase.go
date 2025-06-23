package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/dto/response"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/repository"
)

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrItemNotFound    = errors.New("item not found in cart")
	ErrCartExpired     = errors.New("cart has expired")
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidQuantity = errors.New("invalid quantity")
)

type ProductService interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
}

type Product struct {
	ID    uuid.UUID
	Name  string
	Price float64
	Stock int
}

type CartUsecase struct {
	cartRepo       repository.CartRepository
	productService ProductService
	cartExpiry     time.Duration
}

func NewCartUsecase(cartRepo repository.CartRepository, productService ProductService, cartExpiry time.Duration) *CartUsecase {
	return &CartUsecase{
		cartRepo:       cartRepo,
		productService: productService,
		cartExpiry:     cartExpiry,
	}
}

func (u *CartUsecase) GetCart(ctx context.Context, userID uuid.UUID) (*response.CartResponse, error) {
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = entity.NewCart(userID, u.cartExpiry)
		if err := u.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	} else if cart.IsExpired() {
		return nil, ErrCartExpired
	}

	return response.NewCartResponse(cart), nil
}

func (u *CartUsecase) AddItem(ctx context.Context, userID uuid.UUID, req *request.AddItemRequest) (*response.CartResponse, error) {
	// Get or create cart
	cart, err := u.getOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if cart is expired
	if cart.IsExpired() {
		return nil, ErrCartExpired
	}

	// Get product details
	product, err := u.productService.GetProduct(ctx, req.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	// Add item to cart
	cart.AddItem(entity.CartItem{
		ProductID: product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Quantity:  req.Quantity,
	})

	// Save cart
	if err := u.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return response.NewCartResponse(cart), nil
}

func (u *CartUsecase) UpdateItem(ctx context.Context, userID uuid.UUID, req *request.UpdateItemRequest) (*response.CartResponse, error) {
	// Get cart
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.IsExpired() {
		return nil, ErrCartExpired
	}

	// Update item quantity
	if !cart.UpdateItemQuantity(req.ProductID, req.Quantity) {
		return nil, ErrItemNotFound
	}

	// Save cart
	if err := u.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return response.NewCartResponse(cart), nil
}

func (u *CartUsecase) RemoveItem(ctx context.Context, userID uuid.UUID, req *request.RemoveItemRequest) (*response.CartResponse, error) {
	// Get cart
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, ErrCartNotFound
	}

	// Check if cart is expired
	if cart.IsExpired() {
		return nil, ErrCartExpired
	}

	// Remove item
	cart.RemoveItem(req.ProductID)

	// Save cart
	if err := u.cartRepo.Update(ctx, cart); err != nil {
		return nil, err
	}

	return response.NewCartResponse(cart), nil
}

func (u *CartUsecase) ClearCart(ctx context.Context, userID uuid.UUID) error {
	// Get cart
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return nil
	}

	// Delete cart
	return u.cartRepo.Delete(ctx, cart.ID)
}

func (u *CartUsecase) getOrCreateCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, error) {
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = entity.NewCart(userID, u.cartExpiry)
		if err := u.cartRepo.Create(ctx, cart); err != nil {
			return nil, err
		}
	}

	return cart, nil
}
