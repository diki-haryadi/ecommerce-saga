package provider

import (
	"context"
	"fmt"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"
	"time"
)

// Config holds payment provider configuration
type Config struct {
	APIKey          string
	APISecret       string
	TimeoutDuration time.Duration
	RetryAttempts   int
	WebhookEndpoint string
}

// PaymentProvider defines the interface for payment providers
type PaymentProvider interface {
	ProcessPayment(ctx context.Context, details *usecase.PaymentDetails, amount float64) (string, error)
	RefundPayment(ctx context.Context, transactionID string, amount float64) error
}

// provider implements PaymentProvider interface
type provider struct {
	config Config
}

// NewPaymentProvider creates a new payment provider based on the type
func NewPaymentProvider(providerType string, config Config) (PaymentProvider, error) {
	switch providerType {
	case "stripe":
		return newStripeProvider(config), nil
	case "midtrans":
		return newMidtransProvider(config), nil
	default:
		return nil, fmt.Errorf("unsupported payment provider: %s", providerType)
	}
}

// Mock implementations for now - these would be replaced with actual implementations
type stripeProvider struct {
	provider
}

func newStripeProvider(config Config) PaymentProvider {
	return &stripeProvider{
		provider: provider{config: config},
	}
}

func (p *stripeProvider) ProcessPayment(ctx context.Context, details *usecase.PaymentDetails, amount float64) (string, error) {
	// TODO: Implement actual Stripe payment processing
	return "stripe_mock_transaction_id", nil
}

func (p *stripeProvider) RefundPayment(ctx context.Context, transactionID string, amount float64) error {
	// TODO: Implement actual Stripe refund
	return nil
}

type midtransProvider struct {
	provider
}

func newMidtransProvider(config Config) PaymentProvider {
	return &midtransProvider{
		provider: provider{config: config},
	}
}

func (p *midtransProvider) ProcessPayment(ctx context.Context, details *usecase.PaymentDetails, amount float64) (string, error) {
	// TODO: Implement actual Midtrans payment processing
	return "midtrans_mock_transaction_id", nil
}

func (p *midtransProvider) RefundPayment(ctx context.Context, transactionID string, amount float64) error {
	// TODO: Implement actual Midtrans refund
	return nil
}
