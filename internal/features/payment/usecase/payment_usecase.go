package usecase

import (
	"context"
	"errors"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/repository"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/usecase"
	"time"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/eventbus"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/payment/provider"
)

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrOrderNotFound    = errors.New("order not found")
	ErrStatusTransition = errors.New("invalid status transition")
	ErrPaymentCompleted = errors.New("payment is already completed")

	ErrNotFound            = errors.New("payment not found")
	ErrInvalidStatus       = errors.New("invalid payment status")
	ErrCompleted           = errors.New("payment is already completed")
	ErrProviderUnavailable = errors.New("payment provider is unavailable")
	ErrInvalidProvider     = errors.New("invalid payment provider")
)

// PaymentProvider defines the interface for payment providers
type PaymentProvider interface {
	ProcessPayment(ctx context.Context, payment *payment.Payment) error
	ValidateWebhook(payload []byte, signature string) error
}

type PaymentUsecase struct {
	paymentRepo     usecase.Repository
	orderRepo       orderRepo.OrderRepository
	paymentProvider provider.PaymentProvider
	eventBus        *eventbus.EventBus
}

// NewPaymentUsecase creates a new payment usecase
func NewPaymentUsecase(
	paymentRepo usecase.Repository,
	orderRepo orderRepo.OrderRepository,
	paymentProvider provider.PaymentProvider,
	eventBus *eventbus.EventBus,
) usecase.Usecase {
	return &PaymentUsecase{
		paymentRepo:     paymentRepo,
		orderRepo:       orderRepo,
		paymentProvider: paymentProvider,
		eventBus:        eventBus,
	}
}

// CreatePayment creates a new payment
func (u *PaymentUsecase) CreatePayment(ctx context.Context, orderID uuid.UUID, amount float64, currency, paymentMethod string) (*usecase.PaymentResponse, error) {
	// Validate order exists
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, usecase.ErrOrderNotFound
	}

	// Create payment record
	p := &usecase.PaymentResponse{
		ID:             uuid.New(),
		OrderID:        orderID,
		UserID:         order.UserID,
		Amount:         amount,
		Currency:       currency,
		usecase.Status: usecase.StatusPending,
		PaymentMethod:  paymentMethod,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := u.paymentRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

// GetPayment retrieves a payment by ID
func (u *PaymentUsecase) GetPayment(ctx context.Context, paymentID uuid.UUID) (*usecase.PaymentResponse, error) {
	p, err := u.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, usecase.ErrNotFound
	}
	return p, nil
}

// ListPayments retrieves a list of payments
func (u *PaymentUsecase) ListPayments(ctx context.Context, userID uuid.UUID, page, limit int32, status string) ([]*usecase.PaymentResponse, int64, error) {
	return u.paymentRepo.List(ctx, userID, page, limit, status)
}

// ProcessPayment processes a payment
func (u *PaymentUsecase) ProcessPayment(ctx context.Context, paymentID uuid.UUID, details *usecase.PaymentDetails) (*usecase.PaymentResponse, error) {
	p, err := u.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, usecase.ErrNotFound
	}

	if p.Status != usecase.StatusPending {
		return nil, usecase.ErrInvalidStatus
	}

	// Process payment with provider
	providerTxID, err := u.paymentProvider.ProcessPayment(ctx, details, p.Amount)
	if err != nil {
		p.Status = usecase.StatusFailed
		_ = u.paymentRepo.Update(ctx, p)
		return nil, err
	}

	p.Status = usecase.StatusCompleted
	p.ProviderTransactionID = providerTxID
	p.UpdatedAt = time.Now()

	if err := u.paymentRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	// Publish payment completed event
	u.eventBus.Publish("payment.completed", map[string]interface{}{
		"payment_id": p.ID,
		"order_id":   p.OrderID,
		"amount":     p.Amount,
		"status":     p.Status,
	})

	return p, nil
}

// RefundPayment processes a refund
func (u *PaymentUsecase) RefundPayment(ctx context.Context, paymentID uuid.UUID, amount float64, reason string) (*usecase.PaymentResponse, string, error) {
	p, err := u.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, "", err
	}
	if p == nil {
		return nil, "", usecase.ErrNotFound
	}

	if p.Status != usecase.StatusCompleted {
		return nil, "", usecase.ErrInvalidStatus
	}

	// Process refund with provider
	if err := u.paymentProvider.RefundPayment(ctx, p.ProviderTransactionID, amount); err != nil {
		return nil, "", err
	}

	p.Status = usecase.StatusRefunded
	p.UpdatedAt = time.Now()

	if err := u.paymentRepo.Update(ctx, p); err != nil {
		return nil, "", err
	}

	// Publish refund completed event
	u.eventBus.Publish("payment.refunded", map[string]interface{}{
		"payment_id": p.ID,
		"order_id":   p.OrderID,
		"amount":     amount,
		"reason":     reason,
	})

	return p, reason, nil
}

// GetPaymentByOrder retrieves a payment by order ID
func (u *PaymentUsecase) GetPaymentByOrder(ctx context.Context, orderID uuid.UUID) (*usecase.PaymentResponse, error) {
	payment, err := u.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

// UpdatePaymentStatus updates the status of a payment
func (u *PaymentUsecase) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, req *payment.UpdatePaymentStatusRequest) (*usecase.PaymentResponse, error) {
	// Get payment
	payment, err := u.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	// Validate status transition
	newStatus := payment.PaymentStatus(req.Status)
	if !payment.CanTransitionTo(newStatus) {
		if payment.IsCompleted() {
			return nil, ErrPaymentCompleted
		}
		return nil, ErrStatusTransition
	}

	// Update status
	if err := u.paymentRepo.UpdateStatus(ctx, id, newStatus); err != nil {
		return nil, err
	}

	// Get updated payment
	payment, err = u.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
