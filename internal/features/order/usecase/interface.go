package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/proto"
)

// Status represents the status of an order
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
	StatusCancelled  Status = "CANCELLED"
	StatusFailed     Status = "FAILED"
)

// OrderItem represents a single item in an order
type OrderItem struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
	Price     float64
	Subtotal  float64
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

// Usecase defines the interface for order business logic
type Usecase interface {
	// CreateOrder creates a new order from the user's cart
	CreateOrder(ctx context.Context, userID uuid.UUID, cartID uuid.UUID, paymentMethod, shippingAddress string) (*pb.Order, error)

	// GetOrder retrieves an order by ID
	GetOrder(ctx context.Context, userID, orderID uuid.UUID) (*pb.Order, error)

	// ListOrders retrieves a list of orders for a user
	ListOrders(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*pb.Order, int32, error)

	// CancelOrder cancels an order
	CancelOrder(ctx context.Context, userID, orderID uuid.UUID, reason string) error

	// UpdateOrderStatus updates the status of an order
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status Status) (*pb.Order, error)
}
