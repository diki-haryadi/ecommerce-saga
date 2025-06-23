package order

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusPaid      Status = "PAID"
	StatusShipped   Status = "SHIPPED"
	StatusDelivered Status = "DELIVERED"
	StatusCancelled Status = "CANCELLED"
)

// Usecase defines the order business logic interface
type Usecase interface {
	CreateOrder(ctx context.Context, userID, cartID uuid.UUID, paymentMethod, shippingAddress string) (*OrderResponse, error)
	GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*OrderResponse, error)
	ListOrders(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*OrderResponse, int64, error)
	CancelOrder(ctx context.Context, userID, orderID uuid.UUID, reason string) error
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status Status) (*OrderResponse, error)
}

type OrderResponse struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Status      Status      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
}

// Common errors
var (
	ErrNotFound          = NewError("order not found")
	ErrCartNotFound      = NewError("cart not found")
	ErrCartEmpty         = NewError("cart is empty")
	ErrCancelled         = NewError("order is already cancelled")
	ErrCompleted         = NewError("order is already completed")
	ErrInvalidStatus     = NewError("invalid order status")
	ErrStatusTransition  = NewError("invalid status transition")
	ErrOrderAlreadyFinal = NewError("order is in final state")
)

// Error represents an order error
type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

// NewError creates a new order error
func NewError(message string) *Error {
	return &Error{message: message}
}
