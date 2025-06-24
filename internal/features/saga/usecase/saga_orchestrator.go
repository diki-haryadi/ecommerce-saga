package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/features/saga/domain/entity"
	"github.com/diki-haryadi/ecommerce-saga/internal/shared/messaging"
)

var (
	ErrInvalidStep      = errors.New("invalid saga step")
	ErrStepTimeout      = errors.New("step timeout")
	ErrSagaAlreadyExist = errors.New("saga already exists")
	ErrNotFound         = errors.New("saga not found")
	ErrAlreadyExist     = errors.New("saga already exists")
	ErrInvalidStatus    = errors.New("invalid status")
)

type SagaRepository interface {
	Create(ctx context.Context, saga *entity.SagaTransaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SagaTransaction, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entity.SagaTransaction, error)
	Update(ctx context.Context, saga *entity.SagaTransaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPendingSagas(ctx context.Context) ([]*entity.SagaTransaction, error)
}

type SagaOrchestrator struct {
	sagaRepo      SagaRepository
	messageBroker messaging.MessageBroker
	stepTimeout   time.Duration
	maxRetries    int
}

func NewSagaOrchestrator(
	sagaRepo SagaRepository,
	messageBroker messaging.MessageBroker,
	stepTimeout time.Duration,
	maxRetries int,
) *SagaOrchestrator {
	return &SagaOrchestrator{
		sagaRepo:      sagaRepo,
		messageBroker: messageBroker,
		stepTimeout:   stepTimeout,
		maxRetries:    maxRetries,
	}
}

func (o *SagaOrchestrator) StartOrderPaymentSaga(ctx context.Context, orderID uuid.UUID) error {
	// Check if saga already exists for this order
	if _, err := o.sagaRepo.GetByOrderID(ctx, orderID); err == nil {
		return ErrSagaAlreadyExist
	}

	// Create new saga transaction
	saga := entity.NewSagaTransaction(orderID)
	saga.Timeout = o.stepTimeout
	saga.MaxRetries = o.maxRetries

	// Save saga to postgres
	if err := o.sagaRepo.Create(ctx, saga); err != nil {
		return fmt.Errorf("failed to create saga: %w", err)
	}

	// Start first step
	return o.processStep(ctx, saga)
}

func (o *SagaOrchestrator) ProcessStepResult(ctx context.Context, sagaID uuid.UUID, step entity.SagaStep, status entity.StepStatus, err error) error {
	saga, err := o.sagaRepo.GetByID(ctx, sagaID)
	if err != nil {
		return fmt.Errorf("saga not found: %w", err)
	}

	// Add step result
	saga.AddStepResult(step, status, err)

	// Handle step result
	switch status {
	case entity.StepStatusSuccess:
		return o.handleSuccessfulStep(ctx, saga)
	case entity.StepStatusFailed:
		return o.handleFailedStep(ctx, saga)
	default:
		return o.sagaRepo.Update(ctx, saga)
	}
}

func (o *SagaOrchestrator) ProcessPendingSagas(ctx context.Context) error {
	sagas, err := o.sagaRepo.GetPendingSagas(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending sagas: %w", err)
	}

	for _, saga := range sagas {
		if err := o.processSaga(ctx, saga); err != nil {
			// Log error but continue processing other sagas
			fmt.Printf("Error processing saga %s: %v\n", saga.ID, err)
		}
	}

	return nil
}

func (o *SagaOrchestrator) processSaga(ctx context.Context, saga *entity.SagaTransaction) error {
	// Check for timeout
	if saga.IsTimeout() {
		saga.UpdateStatus(entity.SagaStatusFailed)
		if err := o.sagaRepo.Update(ctx, saga); err != nil {
			return err
		}
		return o.startCompensation(ctx, saga)
	}

	// Process current step
	return o.processStep(ctx, saga)
}

func (o *SagaOrchestrator) handleSuccessfulStep(ctx context.Context, saga *entity.SagaTransaction) error {
	nextStep := o.getNextStep(saga.CurrentStep)
	if nextStep == "" {
		// No more steps, saga completed
		saga.UpdateStatus(entity.SagaStatusCompleted)
		return o.sagaRepo.Update(ctx, saga)
	}

	// Update saga with next step
	saga.SetCurrentStep(entity.SagaStep{Name: nextStep})
	if err := o.sagaRepo.Update(ctx, saga); err != nil {
		return err
	}

	// Process next step
	return o.processStep(ctx, saga)
}

func (o *SagaOrchestrator) handleFailedStep(ctx context.Context, saga *entity.SagaTransaction) error {
	// Mark saga as failed
	saga.UpdateStatus(entity.SagaStatusFailed)
	if err := o.sagaRepo.Update(ctx, saga); err != nil {
		return err
	}

	// Start compensation
	return o.startCompensation(ctx, saga)
}

func (o *SagaOrchestrator) processStep(ctx context.Context, saga *entity.SagaTransaction) error {
	// Prepare step message
	msg := map[string]interface{}{
		"saga_id":  saga.ID,
		"order_id": saga.OrderID,
		"step":     saga.CurrentStep,
	}

	// Get topic for current step
	topic := o.getTopicForStep(saga.CurrentStep.Name)

	// Publish step message
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal step message: %w", err)
	}

	return o.messageBroker.Publish(ctx, topic, msgBytes)
}

func (o *SagaOrchestrator) startCompensation(ctx context.Context, saga *entity.SagaTransaction) error {
	// Process compensation steps in reverse order
	for i := len(saga.CompensationSteps) - 1; i >= 0; i-- {
		step := saga.CompensationSteps[i]
		msg := map[string]interface{}{
			"saga_id":  saga.ID,
			"order_id": saga.OrderID,
			"step":     step,
		}

		// Get compensation topic
		topic := o.getCompensationTopicForStep(step.Name)

		// Publish compensation message
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal compensation message: %w", err)
		}

		if err := o.messageBroker.Publish(ctx, topic, msgBytes); err != nil {
			return err
		}
	}

	return nil
}

func (o *SagaOrchestrator) getNextStep(currentStep entity.SagaStep) entity.StepType {
	switch currentStep.Name {
	case entity.StepCreateOrder:
		return entity.StepProcessPayment
	case entity.StepProcessPayment:
		return entity.StepUpdateInventory
	case entity.StepUpdateInventory:
		return "" // No more steps
	default:
		return ""
	}
}

func (o *SagaOrchestrator) getTopicForStep(step entity.StepType) string {
	return fmt.Sprintf("saga.%s", step)
}

func (o *SagaOrchestrator) getCompensationTopicForStep(step entity.StepType) string {
	return fmt.Sprintf("saga.compensation.%s", step)
}
