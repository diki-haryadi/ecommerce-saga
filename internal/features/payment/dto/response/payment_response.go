package response

import (
	"time"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
)

// PaymentResponse represents a payment in responses
type PaymentResponse struct {
	ID                    uuid.UUID `json:"id"`
	OrderID               uuid.UUID `json:"order_id"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`
	Status                string    `json:"status"`
	Provider              string    `json:"provider"`
	ProviderTransactionID string    `json:"provider_transaction_id,omitempty"`
	ErrorMessage          string    `json:"error_message,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// NewPaymentResponse creates a new payment response from a payment entity
func NewPaymentResponse(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:                    payment.ID,
		OrderID:               payment.OrderID,
		Amount:                payment.Amount,
		Currency:              payment.Currency,
		Status:                string(payment.Status),
		Provider:              string(payment.Provider),
		ProviderTransactionID: payment.ProviderTransactionID,
		ErrorMessage:          payment.ErrorMessage,
		CreatedAt:             payment.CreatedAt,
		UpdatedAt:             payment.UpdatedAt,
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
