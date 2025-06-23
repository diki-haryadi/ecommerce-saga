package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
	StatusRefunded  Status = "REFUNDED"
)

// PaymentDetails represents payment card details
type PaymentDetails struct {
	CardNumber  string
	ExpiryMonth string
	ExpiryYear  string
	CVV         string
	HolderName  string
}

// Usecase defines the payment business logic interface
type Usecase interface {
	CreatePayment(ctx context.Context, orderID uuid.UUID, amount float64, currency, paymentMethod string) (*PaymentResponse, error)
	GetPayment(ctx context.Context, paymentID uuid.UUID) (*PaymentResponse, error)
	ListPayments(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*PaymentResponse, int64, error)
	ProcessPayment(ctx context.Context, paymentID uuid.UUID, details *PaymentDetails) (*PaymentResponse, error)
	RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) (*PaymentResponse, string, error)
}

type PaymentResponse struct {
	ID                    uuid.UUID `json:"id"`
	OrderID               uuid.UUID `json:"order_id"`
	UserID                uuid.UUID `json:"user_id"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`
	Status                Status    `json:"status"`
	PaymentMethod         string    `json:"payment_method"`
	ProviderTransactionID string    `json:"provider_transaction_id,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Common errors
var (
	ErrNotFound            = NewError("payment not found")
	ErrOrderNotFound       = NewError("order not found")
	ErrInvalidProvider     = NewError("invalid payment provider")
	ErrCompleted           = NewError("payment already completed")
	ErrInvalidStatus       = NewError("invalid payment status")
	ErrProviderUnavailable = NewError("payment provider unavailable")
)

// Error represents a payment error
type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

// NewError creates a new payment error
func NewError(message string) *Error {
	return &Error{message: message}
}
