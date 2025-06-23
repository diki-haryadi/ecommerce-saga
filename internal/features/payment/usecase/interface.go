package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/grpc/proto"
)

// Status represents the status of a payment
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusRefunded   Status = "REFUNDED"
)

// PaymentDetails represents payment card details
type PaymentDetails struct {
	CardNumber  string
	ExpiryMonth string
	ExpiryYear  string
	CVV         string
	HolderName  string
}

// PaymentResponse represents a payment
type PaymentResponse struct {
	ID                    uuid.UUID
	OrderID               uuid.UUID
	UserID                uuid.UUID
	Amount                float64
	Currency              string
	Status                Status
	PaymentMethod         string
	ProviderTransactionID string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	OrderID       uuid.UUID
	Amount        float64
	Currency      string
	PaymentMethod string
}

// ListPaymentsRequest represents the request to list payments
type ListPaymentsRequest struct {
	Page   int
	Limit  int
	Status string
}

// Usecase defines the interface for payment business logic
type Usecase interface {
	// CreatePayment creates a new payment for an order
	CreatePayment(ctx context.Context, orderID uuid.UUID, amount float64, currency, paymentMethod string) (*pb.Payment, error)

	// GetPayment retrieves a payment by ID
	GetPayment(ctx context.Context, paymentID uuid.UUID) (*pb.Payment, error)

	// ListPayments retrieves a list of payments for a user
	ListPayments(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*pb.Payment, int32, error)

	// ProcessPayment processes a payment with the provided details
	ProcessPayment(ctx context.Context, paymentID uuid.UUID, details *pb.PaymentDetails) (*pb.Payment, error)

	// RefundPayment processes a refund for a payment
	RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) (*pb.Payment, string, error)
}
