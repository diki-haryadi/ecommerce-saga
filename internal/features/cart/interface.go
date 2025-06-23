package cart

import (
	"context"

	"github.com/google/uuid"
)

// Usecase defines the cart business logic interface
type Usecase interface {
	AddItem(ctx context.Context, userID, productID uuid.UUID, quantity int32) (*Cart, error)
	RemoveItem(ctx context.Context, userID, cartItemID uuid.UUID) error
	UpdateItem(ctx context.Context, userID, cartItemID uuid.UUID, quantity int32) (*Cart, error)
	GetCart(ctx context.Context, userID uuid.UUID) (*Cart, error)
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

// Cart represents a shopping cart
type Cart struct {
	ID     uuid.UUID  `json:"id"`
	UserID uuid.UUID  `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	Price     float64   `json:"price"`
}

// Common errors
var (
	ErrCartNotFound     = NewError("cart not found")
	ErrCartItemNotFound = NewError("cart item not found")
	ErrProductNotFound  = NewError("product not found")
	ErrInvalidQuantity  = NewError("invalid quantity")
	ErrOutOfStock       = NewError("product out of stock")
)

// Error represents a cart error
type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

// NewError creates a new cart error
func NewError(message string) *Error {
	return &Error{message: message}
}
