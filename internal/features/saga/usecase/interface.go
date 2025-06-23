package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	pb "github.com/diki-haryadi/ecommerce-saga/internal/features/saga/delivery/grpc/proto"
)

// SagaStatus represents the status of a saga transaction
type SagaStatus string

const (
	SagaStatusPending      SagaStatus = "PENDING"
	SagaStatusProcessing   SagaStatus = "PROCESSING"
	SagaStatusCompleted    SagaStatus = "COMPLETED"
	SagaStatusFailed       SagaStatus = "FAILED"
	SagaStatusCompensating SagaStatus = "COMPENSATING"
	SagaStatusCompensated  SagaStatus = "COMPENSATED"
)

// SagaType represents the type of saga transaction
type SagaType string

const (
	SagaTypeOrder SagaType = "ORDER"
)

// StepStatus represents the status of a saga step
type StepStatus string

const (
	StepStatusPending      StepStatus = "PENDING"
	StepStatusCompleted    StepStatus = "COMPLETED"
	StepStatusFailed       StepStatus = "FAILED"
	StepStatusCompensating StepStatus = "COMPENSATING"
	StepStatusCompensated  StepStatus = "COMPENSATED"
)

// SagaStep represents a step in the saga transaction
type SagaStep struct {
	ID                 uuid.UUID
	Name               string
	Status             StepStatus
	Service            string
	Action             string
	CompensationAction string
	Payload            map[string]interface{}
	ErrorMessage       string
	ExecutedAt         time.Time
}

// SagaTransactionResponse represents a saga transaction
type SagaTransactionResponse struct {
	ID        uuid.UUID
	Type      SagaType
	Status    SagaStatus
	Steps     []*SagaStep
	Metadata  map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// StartOrderSagaRequest represents the request to start an order saga
type StartOrderSagaRequest struct {
	OrderID       uuid.UUID
	UserID        uuid.UUID
	Amount        float64
	PaymentMethod string
	Metadata      map[string]string
}

// ListSagaTransactionsRequest represents the request to list saga transactions
type ListSagaTransactionsRequest struct {
	Page   int
	Limit  int
	Status string
	Type   string
}

// Usecase defines the interface for saga orchestration
type Usecase interface {
	// StartOrderSaga starts a new order saga transaction
	StartOrderSaga(ctx context.Context, orderID, userID uuid.UUID, amount float64, paymentMethod string, metadata map[string]string) (*pb.SagaTransaction, error)

	// GetSagaStatus retrieves the status of a saga transaction
	GetSagaStatus(ctx context.Context, sagaID uuid.UUID) (*pb.SagaTransaction, error)

	// CompensateTransaction initiates compensation for a saga transaction
	CompensateTransaction(ctx context.Context, sagaID, stepID uuid.UUID, reason string) (*pb.SagaTransaction, error)

	// ListSagaTransactions retrieves a list of saga transactions
	ListSagaTransactions(ctx context.Context, page, limit int32, status, transactionType string) ([]*pb.SagaTransaction, int32, error)
}
