package usecase

import (
	"context"
	//pb "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/proto"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPaid       Status = "PAID"
	StatusShipped    Status = "SHIPPED"
	StatusDelivered  Status = "DELIVERED"
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
	StatusCancelled  Status = "CANCELLED"
	StatusFailed     Status = "FAILED"
)

// Usecase defines the order business logic interface
type Usecase interface {
	CreateOrder(ctx context.Context, userID, cartID uuid.UUID, paymentMethod, shippingAddress string) (*OrderResponse, error)
	GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*OrderResponse, error)
	ListOrders(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*OrderResponse, int64, error)
	CancelOrder(ctx context.Context, userID, orderID uuid.UUID, reason string) error
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status Status) (*OrderResponse, error)
	// CreateOrder creates a new order from the user's cart
	//CreateOrder(ctx context.Context, userID uuid.UUID, cartID uuid.UUID, paymentMethod, shippingAddress string) (*pb.Order, error)
	//
	//// GetOrder retrieves an order by ID
	//GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*pb.Order, error)
	//
	//// ListOrders retrieves a list of orders for a user
	//ListOrders(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*pb.Order, int32, error)
	//
	//// CancelOrder cancels an order
	//CancelOrder(ctx context.Context, userID, orderID uuid.UUID, reason string) error
	//
	//// UpdateOrderStatus updates the status of an order
	//UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status Status) (*pb.Order, error)
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

// Order represents an order
type Order struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Items           []*OrderItem
	Total           float64
	Status          Status
	PaymentMethod   string
	PaymentID       *uuid.UUID
	ShippingAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	CartID          uuid.UUID
	PaymentMethod   string
	ShippingAddress string
}

// ListOrdersRequest represents the request to list orders
type ListOrdersRequest struct {
	Page   int
	Limit  int
	Status string
}
