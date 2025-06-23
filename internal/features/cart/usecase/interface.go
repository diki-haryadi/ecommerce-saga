package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/proto"
)

// CartItem represents a single item in the cart
type CartItem struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
	Price     float64
	Subtotal  float64
}

// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Items     []*CartItem
	Total     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AddItemRequest represents the request to add an item to cart
type AddItemRequest struct {
	ProductID uuid.UUID
	Quantity  int32
}

// UpdateItemRequest represents the request to update a cart item
type UpdateItemRequest struct {
	CartItemID uuid.UUID
	Quantity   int32
}

// Usecase defines the interface for cart business logic
type Usecase interface {
	// AddItem adds a new item to the cart
	AddItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int32) (*pb.Cart, error)

	// RemoveItem removes an item from the cart
	RemoveItem(ctx context.Context, userID, cartItemID uuid.UUID) error

	// UpdateItem updates the quantity of an item in the cart
	UpdateItem(ctx context.Context, userID, cartItemID uuid.UUID, quantity int32) (*pb.Cart, error)

	// GetCart retrieves the user's cart
	GetCart(ctx context.Context, userID uuid.UUID) (*pb.Cart, error)

	// ClearCart removes all items from the cart
	ClearCart(ctx context.Context, userID uuid.UUID) error
}
