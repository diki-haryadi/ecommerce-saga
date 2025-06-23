package entity

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusSuccess    PaymentStatus = "SUCCESS"
	PaymentStatusFailed     PaymentStatus = "FAILED"
)

// PaymentProvider represents the payment provider
type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "STRIPE"
	PaymentProviderPayPal PaymentProvider = "PAYPAL"
)

// Payment represents a payment in the system
type Payment struct {
	ID                    uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	OrderID               uuid.UUID       `json:"order_id" gorm:"type:uuid;not null"`
	Amount                float64         `json:"amount" gorm:"not null"`
	Currency              string          `json:"currency" gorm:"type:varchar(3);not null"`
	Status                PaymentStatus   `json:"status" gorm:"type:varchar(50);not null"`
	Provider              PaymentProvider `json:"provider" gorm:"type:varchar(50);not null"`
	ProviderTransactionID string          `json:"provider_transaction_id" gorm:"type:varchar(255)"`
	ErrorMessage          string          `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
}

// NewPayment creates a new payment
func NewPayment(orderID uuid.UUID, amount float64, currency string, provider PaymentProvider) *Payment {
	return &Payment{
		ID:        uuid.New(),
		OrderID:   orderID,
		Amount:    amount,
		Currency:  currency,
		Status:    PaymentStatusPending,
		Provider:  provider,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateStatus updates the payment status
func (p *Payment) UpdateStatus(status PaymentStatus) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// SetProviderTransactionID sets the provider's transaction ID
func (p *Payment) SetProviderTransactionID(txID string) {
	p.ProviderTransactionID = txID
	p.UpdatedAt = time.Now()
}

// SetError sets the error message and updates status to failed
func (p *Payment) SetError(message string) {
	p.Status = PaymentStatusFailed
	p.ErrorMessage = message
	p.UpdatedAt = time.Now()
}

// IsCompleted checks if the payment is in a final state
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusSuccess || p.Status == PaymentStatusFailed
}

// CanTransitionTo checks if the payment can transition to the given status
func (p *Payment) CanTransitionTo(status PaymentStatus) bool {
	if p.IsCompleted() {
		return false
	}

	switch p.Status {
	case PaymentStatusPending:
		return status == PaymentStatusProcessing || status == PaymentStatusFailed
	case PaymentStatusProcessing:
		return status == PaymentStatusSuccess || status == PaymentStatusFailed
	default:
		return false
	}
}
