package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	orderEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/repository"
	paymentEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	paymentRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/repository"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/repository"
)

var (
	ErrSagaNotFound     = errors.New("saga not found")
	ErrInvalidSagaType  = errors.New("invalid saga type")
	ErrInvalidStepOrder = errors.New("invalid step order")
)

// OrderPaymentPayload represents the payload for order-payment saga
type OrderPaymentPayload struct {
	OrderID  uuid.UUID `json:"order_id"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
}

type SagaUsecase struct {
	sagaRepo    repository.SagaRepository
	orderRepo   orderRepo.OrderRepository
	paymentRepo paymentRepo.PaymentRepository
}

func NewSagaUsecase(
	sagaRepo repository.SagaRepository,
	orderRepo orderRepo.OrderRepository,
	paymentRepo paymentRepo.PaymentRepository,
) *SagaUsecase {
	return &SagaUsecase{
		sagaRepo:    sagaRepo,
		orderRepo:   orderRepo,
		paymentRepo: paymentRepo,
	}
}

// StartOrderPaymentSaga starts a new order-payment saga
func (u *SagaUsecase) StartOrderPaymentSaga(ctx context.Context, orderID uuid.UUID) error {
	// Get order
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Create payload
	payload := OrderPaymentPayload{
		OrderID:  orderID,
		Amount:   order.TotalAmount,
		Currency: "USD",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create saga steps
	steps := []entity.SagaStep{
		{
			Name:    "ValidateOrder",
			Payload: payloadBytes,
		},
		{
			Name:    "ProcessPayment",
			Payload: payloadBytes,
		},
		{
			Name:    "UpdateOrderStatus",
			Payload: payloadBytes,
		},
	}

	// Create saga
	saga := entity.NewSaga(entity.SagaTypeOrderPayment, steps)
	if err := u.sagaRepo.Create(ctx, saga); err != nil {
		return err
	}

	// Start saga execution
	go u.executeSaga(context.Background(), saga)

	return nil
}

// executeSaga executes a saga
func (u *SagaUsecase) executeSaga(ctx context.Context, saga *entity.Saga) {
	for {
		step := saga.GetNextStep()
		if step == nil {
			break
		}

		var err error
		switch step.Name {
		case "ValidateOrder":
			err = u.executeValidateOrder(ctx, step)
		case "ProcessPayment":
			err = u.executeProcessPayment(ctx, step)
		case "UpdateOrderStatus":
			err = u.executeUpdateOrderStatus(ctx, step)
		}

		if err != nil {
			u.handleStepFailure(ctx, saga, step, err)
			return
		}

		step.Status = entity.StepStatusCompleted
		if err := u.sagaRepo.UpdateStepStatus(ctx, saga.ID, step.ID, entity.StepStatusCompleted, ""); err != nil {
			// Log error but continue
			continue
		}
	}
}

// executeValidateOrder executes the ValidateOrder step
func (u *SagaUsecase) executeValidateOrder(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	order, err := u.orderRepo.GetByID(ctx, payload.OrderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	if order.Status != orderEntity.OrderStatusPending {
		return errors.New("invalid order status")
	}

	return nil
}

// executeProcessPayment executes the ProcessPayment step
func (u *SagaUsecase) executeProcessPayment(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	payment := paymentEntity.NewPayment(
		payload.OrderID,
		payload.Amount,
		payload.Currency,
		paymentEntity.PaymentProviderStripe,
	)

	if err := u.paymentRepo.Create(ctx, payment); err != nil {
		return err
	}

	// In a real application, we would integrate with the payment provider here
	// For now, we'll simulate a successful payment
	time.Sleep(2 * time.Second)
	payment.Status = paymentEntity.PaymentStatusSuccess
	payment.UpdatedAt = time.Now()

	return u.paymentRepo.Update(ctx, payment)
}

// executeUpdateOrderStatus executes the UpdateOrderStatus step
func (u *SagaUsecase) executeUpdateOrderStatus(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	payment, err := u.paymentRepo.GetByOrderID(ctx, payload.OrderID)
	if err != nil {
		return err
	}

	var orderStatus orderEntity.OrderStatus
	switch payment.Status {
	case paymentEntity.PaymentStatusSuccess:
		orderStatus = orderEntity.OrderStatusConfirmed
	case paymentEntity.PaymentStatusFailed:
		orderStatus = orderEntity.OrderStatusFailed
	default:
		return errors.New("invalid payment status")
	}

	return u.orderRepo.UpdateStatus(ctx, payload.OrderID, orderStatus)
}

// handleStepFailure handles a step failure
func (u *SagaUsecase) handleStepFailure(ctx context.Context, saga *entity.Saga, step *entity.SagaStep, err error) {
	// Update step status
	u.sagaRepo.UpdateStepStatus(ctx, saga.ID, step.ID, entity.StepStatusFailed, err.Error())

	// Start compensation
	go u.compensateSaga(context.Background(), saga)
}

// compensateSaga compensates a failed saga
func (u *SagaUsecase) compensateSaga(ctx context.Context, saga *entity.Saga) {
	// Reverse through completed steps
	for i := len(saga.Steps) - 1; i >= 0; i-- {
		step := &saga.Steps[i]
		if step.Status != entity.StepStatusCompleted {
			continue
		}

		var err error
		switch step.Name {
		case "ProcessPayment":
			err = u.compensateProcessPayment(ctx, step)
		case "UpdateOrderStatus":
			err = u.compensateUpdateOrderStatus(ctx, step)
		}

		if err != nil {
			// Log error but continue compensation
			continue
		}

		step.Status = entity.StepStatusCompensated
		u.sagaRepo.UpdateStepStatus(ctx, saga.ID, step.ID, entity.StepStatusCompensated, "")
	}
}

// compensateProcessPayment compensates the ProcessPayment step
func (u *SagaUsecase) compensateProcessPayment(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	payment, err := u.paymentRepo.GetByOrderID(ctx, payload.OrderID)
	if err != nil {
		return err
	}

	payment.Status = paymentEntity.PaymentStatusFailed
	payment.ErrorMessage = "Payment compensated due to saga failure"
	payment.UpdatedAt = time.Now()

	return u.paymentRepo.Update(ctx, payment)
}

// compensateUpdateOrderStatus compensates the UpdateOrderStatus step
func (u *SagaUsecase) compensateUpdateOrderStatus(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	return u.orderRepo.UpdateStatus(ctx, payload.OrderID, orderEntity.OrderStatusFailed)
}
