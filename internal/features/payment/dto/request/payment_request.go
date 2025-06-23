package request

import "github.com/google/uuid"

// ProcessPaymentRequest represents the request to process a payment
type ProcessPaymentRequest struct {
	OrderID  uuid.UUID `json:"order_id" validate:"required"`
	Provider string    `json:"provider" validate:"required,oneof=STRIPE PAYPAL"`
}

// UpdatePaymentStatusRequest represents the request to update a payment's status
type UpdatePaymentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PROCESSING SUCCESS FAILED"`
}

// PaymentWebhookRequest represents a payment webhook request from a provider
type PaymentWebhookRequest struct {
	ProviderTransactionID string `json:"provider_transaction_id" validate:"required"`
	Status                string `json:"status" validate:"required"`
	ErrorMessage          string `json:"error_message,omitempty"`
}
