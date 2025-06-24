package usecase

import (
	"context"
	"encoding/json"
	"errors"
	orderRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/repository"
	paymentRepo "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/repository"
	"time"

	"github.com/google/uuid"

	cartClient "github.com/diki-haryadi/ecommerce-saga/internal/features/cart/delivery/grpc/client"
	orderClient "github.com/diki-haryadi/ecommerce-saga/internal/features/order/delivery/grpc/client"
	orderEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/order/domain/entity"
	paymentClient "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/delivery/grpc/client"
	paymentEntity "github.com/diki-haryadi/ecommerce-saga/internal/features/payment/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga"
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
	sagaRepo      repository.SagaRepository
	orderRepo     orderRepo.OrderRepository
	paymentRepo   paymentRepo.PaymentRepository
	orderClient   *orderClient.OrderClient
	paymentClient *paymentClient.PaymentClient
	cartClient    *cartClient.CartClient
}

func NewSagaUsecase(
	sagaRepo repository.SagaRepository,
	orderRepo orderRepo.OrderRepository,
	paymentRepo paymentRepo.PaymentRepository,
	orderClient *orderClient.OrderClient,
	paymentClient *paymentClient.PaymentClient,
	cartClient *cartClient.CartClient,
) *SagaUsecase {
	return &SagaUsecase{
		sagaRepo:      sagaRepo,
		orderRepo:     orderRepo,
		paymentRepo:   paymentRepo,
		orderClient:   orderClient,
		paymentClient: paymentClient,
		cartClient:    cartClient,
	}
}

// StartOrderSaga starts a new order saga transaction
func (u *SagaUsecase) StartOrderSaga(ctx context.Context, orderID, userID uuid.UUID, amount float64, paymentMethod string, metadata map[string]string) (*saga.SagaResponse, error) {
	// Create payload
	payload := OrderPaymentPayload{
		OrderID:  orderID,
		Amount:   amount,
		Currency: "USD",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create saga steps
	steps := []entity.SagaStep{
		{
			Name:    entity.StepCreateOrder,
			Payload: payloadBytes,
		},
		{
			Name:    entity.StepProcessPayment,
			Payload: payloadBytes,
		},
		{
			Name:    entity.StepUpdateInventory,
			Payload: payloadBytes,
		},
	}

	// Create saga
	sagaEntity := entity.NewSaga(entity.SagaTypeOrderPayment, steps)
	if err := u.sagaRepo.Create(ctx, sagaEntity); err != nil {
		return nil, err
	}

	// Start saga execution
	go u.executeSaga(context.Background(), sagaEntity)

	// Convert to response
	return &saga.SagaResponse{
		ID:        sagaEntity.ID,
		Type:      saga.Type(string(sagaEntity.Type)),
		Status:    saga.Status(string(sagaEntity.Status)),
		Steps:     convertStepsToResponse(sagaEntity.Steps),
		Metadata:  metadata,
		CreatedAt: sagaEntity.CreatedAt,
		UpdatedAt: sagaEntity.UpdatedAt,
	}, nil
}

// GetSagaStatus retrieves the status of a saga transaction
func (u *SagaUsecase) GetSagaStatus(ctx context.Context, sagaID uuid.UUID) (*saga.SagaResponse, error) {
	sagaEntity, err := u.sagaRepo.GetByID(ctx, sagaID)
	if err != nil {
		return nil, err
	}

	return &saga.SagaResponse{
		ID:        sagaEntity.ID,
		Type:      saga.Type(string(sagaEntity.Type)),
		Status:    saga.Status(string(sagaEntity.Status)),
		Steps:     convertStepsToResponse(sagaEntity.Steps),
		Metadata:  make(map[string]string),
		CreatedAt: sagaEntity.CreatedAt,
		UpdatedAt: sagaEntity.UpdatedAt,
	}, nil
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
		case entity.StepCreateOrder:
			err = u.executeCreateOrder(ctx, step)
		case entity.StepProcessPayment:
			err = u.executeProcessPayment(ctx, step)
		case entity.StepUpdateInventory:
			err = u.executeUpdateInventory(ctx, step)
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

// executeCreateOrder executes the CreateOrder step
func (u *SagaUsecase) executeCreateOrder(ctx context.Context, step *entity.SagaStep) error {
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

// executeUpdateInventory executes the UpdateInventory step
func (u *SagaUsecase) executeUpdateInventory(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	// Implementation of updating inventory
	return nil
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
		case entity.StepProcessPayment:
			err = u.compensateProcessPayment(ctx, step)
		case entity.StepUpdateInventory:
			err = u.compensateUpdateInventory(ctx, step)
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

// compensateUpdateInventory compensates the UpdateInventory step
func (u *SagaUsecase) compensateUpdateInventory(ctx context.Context, step *entity.SagaStep) error {
	var payload OrderPaymentPayload
	if err := json.Unmarshal(step.Payload, &payload); err != nil {
		return err
	}

	// Implementation of updating inventory
	return nil
}

// CompensateTransaction initiates compensation for a saga transaction
func (u *SagaUsecase) CompensateTransaction(ctx context.Context, sagaID, stepID uuid.UUID, reason string) (*saga.SagaResponse, error) {
	sagaEntity, err := u.sagaRepo.GetByID(ctx, sagaID)
	if err != nil {
		return nil, err
	}

	// Start compensation process
	go u.compensateSaga(context.Background(), sagaEntity)

	return &saga.SagaResponse{
		ID:        sagaEntity.ID,
		Type:      saga.Type(string(sagaEntity.Type)),
		Status:    saga.Status(string(sagaEntity.Status)),
		Steps:     convertStepsToResponse(sagaEntity.Steps),
		Metadata:  make(map[string]string),
		CreatedAt: sagaEntity.CreatedAt,
		UpdatedAt: sagaEntity.UpdatedAt,
	}, nil
}

// ListSagaTransactions retrieves a list of saga transactions
func (u *SagaUsecase) ListSagaTransactions(ctx context.Context, page, limit int32, status string, sagaType saga.Type) ([]*saga.SagaResponse, int64, error) {
	// TODO: Implement pagination and filtering
	return nil, 0, nil
}

func convertStepsToResponse(steps []entity.SagaStep) []saga.SagaStep {
	result := make([]saga.SagaStep, len(steps))
	for i, step := range steps {
		result[i] = saga.SagaStep{
			ID:                 step.ID,
			Name:               string(step.Name),
			Status:             saga.Status(string(step.Status)),
			Service:            "saga",
			Action:             string(step.Name),
			CompensationAction: "compensate_" + string(step.Name),
			Payload:            make(map[string]interface{}),
			ErrorMessage:       step.ErrorMessage,
			ExecutedAt:         step.CreatedAt,
		}
	}
	return result
}

// StartOrderPaymentSaga starts a new order-payment saga transaction
func (u *SagaUsecase) StartOrderPaymentSaga(ctx context.Context, orderID uuid.UUID) error {
	// Default values for demonstration
	amount := 0.0           // This should be fetched from the order
	userID := uuid.New()    // This should be fetched from the order
	paymentMethod := "card" // This should be fetched from the order
	metadata := make(map[string]string)

	_, err := u.StartOrderSaga(ctx, orderID, userID, amount, paymentMethod, metadata)
	return err
}
