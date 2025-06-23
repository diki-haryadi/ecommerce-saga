package response

import (
	"time"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/cart/domain/entity"
)

// CartItemResponse represents a cart item in responses
type CartItemResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
}

// CartResponse represents a cart in responses
type CartResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	Total     float64            `json:"total"`
	ExpiresAt time.Time          `json:"expires_at"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// NewCartResponse creates a new cart response from a cart entity
func NewCartResponse(cart *entity.Cart) *CartResponse {
	items := make([]CartItemResponse, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = CartItemResponse{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Price * float64(item.Quantity),
		}
	}

	return &CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     items,
		Total:     cart.Total,
		ExpiresAt: cart.ExpiresAt,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"description,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
