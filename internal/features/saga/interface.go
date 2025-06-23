package saga

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string
type Type string

const (
	StatusPending     Status = "PENDING"
	StatusInProgress  Status = "IN_PROGRESS"
	StatusCompleted   Status = "COMPLETED"
	StatusFailed      Status = "FAILED"
	StatusCompensated Status = "COMPENSATED"

	TypeOrder Type = "ORDER"
)

// Usecase defines the saga orchestrator interface
type Usecase interface {
	StartOrderSaga(ctx context.Context, orderID, userID uuid.UUID, amount float64, paymentMethod string, metadata map[string]string) (*SagaResponse, error)
	GetSagaStatus(ctx context.Context, sagaID uuid.UUID) (*SagaResponse, error)
	CompensateTransaction(ctx context.Context, sagaID, stepID uuid.UUID, reason string) (*SagaResponse, error)
	ListSagaTransactions(ctx context.Context, page, limit int32, status string, sagaType Type) ([]*SagaResponse, int64, error)
}

type SagaResponse struct {
	ID        uuid.UUID         `json:"id"`
	Type      Type              `json:"type"`
	Status    Status            `json:"status"`
	Steps     []SagaStep        `json:"steps"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type SagaStep struct {
	ID                 uuid.UUID              `json:"id"`
	Name               string                 `json:"name"`
	Status             Status                 `json:"status"`
	Service            string                 `json:"service"`
	Action             string                 `json:"action"`
	CompensationAction string                 `json:"compensation_action"`
	Payload            map[string]interface{} `json:"payload"`
	ErrorMessage       string                 `json:"error_message,omitempty"`
	ExecutedAt         time.Time              `json:"executed_at"`
}

// Common errors
var (
	ErrNotFound      = NewError("saga not found")
	ErrAlreadyExists = NewError("saga already exists")
	ErrInvalidStep   = NewError("invalid step order")
)

// Error represents a saga error
type Error struct {
	message string
}

func (e *Error) Error() string {
	return e.message
}

// NewError creates a new saga error
func NewError(message string) *Error {
	return &Error{message: message}
}
