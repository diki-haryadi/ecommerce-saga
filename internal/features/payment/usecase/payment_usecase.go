package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"

	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/dto/request"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/dto/response"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/payment/repository"
)

var (
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrInvalidStatus       = errors.New("invalid payment status")
	ErrStatusTransition    = errors.New("invalid status transition")
	ErrPaymentCompleted    = errors.New("payment is already completed")
	ErrInvalidProvider     = errors.New("invalid payment provider")
	ErrProviderUnavailable = errors.New("payment provider is unavailable")
)

// PaymentProvider defines the interface for payment providers
type PaymentProvider interface {
	ProcessPayment(ctx context.Context, payment *entity.Payment) error
	ValidateWebhook(payload []byte, signature string) error
}

type PaymentUsecase struct {
	paymentRepo      repository.PaymentRepository
	orderRepo        orderRepo.OrderRepository
	paymentProviders map[entity.PaymentProvider]PaymentProvider
}

func NewPaymentUsecase(
	paymentRepo repository.PaymentRepository,
	orderRepo orderRepo.OrderRepository,
	providers map[entity.PaymentProvider]PaymentProvider,
) *PaymentUsecase {
	return &PaymentUsecase{
		paymentRepo:      paymentRepo,
		orderRepo:        orderRepo,
		paymentProviders: providers,
	}
}

// ProcessPayment processes a new payment for an order
func (u *PaymentUsecase) ProcessPayment(ctx context.Context, req *request.ProcessPaymentRequest) (*response.PaymentResponse, error) {
	// Get order
	order, err := u.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Check if payment already exists
	existingPayment, err := u.paymentRepo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if existingPayment != nil {
		return response.NewPaymentResponse(existingPayment), nil
	}

	// Create payment
	provider := entity.PaymentProvider(req.Provider)
	payment := entity.NewPayment(req.OrderID, order.TotalAmount, "USD", provider)

	// Save payment
	if err := u.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// Get payment provider
	paymentProvider, ok := u.paymentProviders[provider]
	if !ok {
		return nil, ErrInvalidProvider
	}

	// Process payment asynchronously
	go func() {
		ctx := context.Background()
		if err := paymentProvider.ProcessPayment(ctx, payment); err != nil {
			payment.SetError(err.Error())
			u.paymentRepo.Update(ctx, payment)
		}
	}()

	return response.NewPaymentResponse(payment), nil
}

// GetPayment retrieves a payment by ID
func (u *PaymentUsecase) GetPayment(ctx context.Context, id uuid.UUID) (*response.PaymentResponse, error) {
	payment, err := u.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return response.NewPaymentResponse(payment), nil
}

// GetPaymentByOrder retrieves a payment by order ID
func (u *PaymentUsecase) GetPaymentByOrder(ctx context.Context, orderID uuid.UUID) (*response.PaymentResponse, error) {
	payment, err := u.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return response.NewPaymentResponse(payment), nil
}

// UpdatePaymentStatus updates the status of a payment
func (u *PaymentUsecase) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, req *request.UpdatePaymentStatusRequest) (*response.PaymentResponse, error) {
	// Get payment
	payment, err := u.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	// Validate status transition
	newStatus := entity.PaymentStatus(req.Status)
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

	return response.NewPaymentResponse(payment), nil
}

// HandleWebhook handles payment provider webhook
func (u *PaymentUsecase) HandleWebhook(ctx context.Context, provider entity.PaymentProvider, payload []byte, signature string) error {
	// Get payment provider
	paymentProvider, ok := u.paymentProviders[provider]
	if !ok {
		return ErrInvalidProvider
	}

	// Validate webhook
	if err := paymentProvider.ValidateWebhook(payload, signature); err != nil {
		return err
	}

	// Process webhook asynchronously
	go func() {
		// Implementation depends on provider-specific webhook format
		// Update payment status based on webhook data
	}()

	return nil
}
