package entity

import (
	"time"

	"github.com/google/uuid"
)

// SagaStatus represents the status of a saga
type SagaStatus string

const (
	SagaStatusPending      SagaStatus = "PENDING"
	SagaStatusProcessing   SagaStatus = "PROCESSING"
	SagaStatusCompleted    SagaStatus = "COMPLETED"
	SagaStatusFailed       SagaStatus = "FAILED"
	SagaStatusCompensating SagaStatus = "COMPENSATING"
)

// StepStatus represents the status of a saga step
type StepStatus string

const (
	StepStatusPending     StepStatus = "PENDING"
	StepStatusSuccess     StepStatus = "SUCCESS"
	StepStatusFailed      StepStatus = "FAILED"
	StepStatusCancelled   StepStatus = "CANCELLED"
	StepStatusCompleted   StepStatus = "COMPLETED"
	StepStatusCompensated StepStatus = "COMPENSATED"
)

// StepType represents the type of step in the saga
type StepType string

const (
	StepCreateOrder     StepType = "CREATE_ORDER"
	StepProcessPayment  StepType = "PROCESS_PAYMENT"
	StepUpdateInventory StepType = "UPDATE_INVENTORY"
)

// SagaType represents the type of saga
type SagaType string

const (
	SagaTypeOrderPayment SagaType = "ORDER_PAYMENT"
)

// SagaStep represents a step in the saga
type SagaStep struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SagaID       uuid.UUID  `json:"saga_id" gorm:"type:uuid;not null"`
	Name         StepType   `json:"name" gorm:"type:varchar(255);not null"`
	Status       StepStatus `json:"status" gorm:"type:varchar(50);not null"`
	Order        int        `json:"order" gorm:"not null"`
	Payload      []byte     `json:"payload" gorm:"type:jsonb"`
	ErrorMessage string     `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Saga represents a distributed transaction
type Saga struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Type      SagaType   `json:"type" gorm:"type:varchar(50);not null"`
	Status    SagaStatus `json:"status" gorm:"type:varchar(50);not null"`
	Steps     []SagaStep `json:"steps" gorm:"foreignKey:SagaID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NewSaga creates a new saga
func NewSaga(sagaType SagaType, steps []SagaStep) *Saga {
	now := time.Now()
	sagaID := uuid.New()

	// Initialize steps
	for i := range steps {
		steps[i].ID = uuid.New()
		steps[i].SagaID = sagaID
		steps[i].Status = StepStatusPending
		steps[i].Order = i
		steps[i].CreatedAt = now
		steps[i].UpdatedAt = now
	}

	return &Saga{
		ID:        sagaID,
		Type:      sagaType,
		Status:    SagaStatusPending,
		Steps:     steps,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateStepStatus updates the status of a step
func (s *Saga) UpdateStepStatus(stepID uuid.UUID, status StepStatus, errorMessage string) {
	for i := range s.Steps {
		if s.Steps[i].ID == stepID {
			s.Steps[i].Status = status
			s.Steps[i].ErrorMessage = errorMessage
			s.Steps[i].UpdatedAt = time.Now()
			break
		}
	}
	s.UpdateStatus()
}

// UpdateStatus updates the saga status based on step statuses
func (s *Saga) UpdateStatus() {
	allCompleted := true
	anyFailed := false
	anyCompensated := false

	for _, step := range s.Steps {
		if step.Status == StepStatusFailed {
			anyFailed = true
			break
		}
		if step.Status == StepStatusCompensated {
			anyCompensated = true
		}
		if step.Status != StepStatusCompleted {
			allCompleted = false
		}
	}

	if anyFailed {
		s.Status = SagaStatusFailed
	} else if anyCompensated {
		s.Status = SagaStatusCompensating
	} else if allCompleted {
		s.Status = SagaStatusCompleted
	}

	s.UpdatedAt = time.Now()
}

// GetNextStep gets the next pending step
func (s *Saga) GetNextStep() *SagaStep {
	for _, step := range s.Steps {
		if step.Status == StepStatusPending {
			return &step
		}
	}
	return nil
}

// GetStepByID gets a step by its ID
func (s *Saga) GetStepByID(stepID uuid.UUID) *SagaStep {
	for _, step := range s.Steps {
		if step.ID == stepID {
			return &step
		}
	}
	return nil
}

// IsCompleted checks if the saga is completed
func (s *Saga) IsCompleted() bool {
	return s.Status == SagaStatusCompleted
}

// IsFailed checks if the saga has failed
func (s *Saga) IsFailed() bool {
	return s.Status == SagaStatusFailed
}

// IsCompensating checks if the saga is in compensating state
func (s *Saga) IsCompensating() bool {
	return s.Status == SagaStatusCompensating
}

type SagaStepResult struct {
	Step      SagaStep    `json:"step"`
	Status    StepStatus  `json:"status"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Retries   int         `json:"retries"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

type SagaTransaction struct {
	ID                uuid.UUID        `json:"id"`
	OrderID           uuid.UUID        `json:"order_id"`
	Status            SagaStatus       `json:"status"`
	Steps             []SagaStepResult `json:"steps"`
	CompensationSteps []SagaStep       `json:"compensation_steps"`
	CurrentStep       SagaStep         `json:"current_step"`
	Error             string           `json:"error,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	Timeout           time.Duration    `json:"timeout"`
	MaxRetries        int              `json:"max_retries"`
}

func NewSagaTransaction(orderID uuid.UUID) *SagaTransaction {
	return &SagaTransaction{
		ID:                uuid.New(),
		OrderID:           orderID,
		Status:            SagaStatusPending,
		Steps:             make([]SagaStepResult, 0),
		CompensationSteps: make([]SagaStep, 0),
		CurrentStep:       SagaStep{Name: StepCreateOrder, Status: StepStatusPending},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Timeout:           5 * time.Minute,
		MaxRetries:        3,
	}
}

func (s *SagaTransaction) AddStepResult(step SagaStep, status StepStatus, err error) {
	result := SagaStepResult{
		Step:      step,
		Status:    status,
		Timestamp: time.Now(),
	}
	if err != nil {
		result.Error = err.Error()
	}
	s.Steps = append(s.Steps, result)
	s.UpdatedAt = time.Now()
}

func (s *SagaTransaction) AddCompensationStep(step SagaStep) {
	s.CompensationSteps = append(s.CompensationSteps, step)
}

func (s *SagaTransaction) UpdateStatus(status SagaStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

func (s *SagaTransaction) SetCurrentStep(step SagaStep) {
	s.CurrentStep = step
	s.UpdatedAt = time.Now()
}

func (s *SagaTransaction) IsCompleted() bool {
	return s.Status == SagaStatusCompleted
}

func (s *SagaTransaction) IsFailed() bool {
	return s.Status == SagaStatusFailed
}

func (s *SagaTransaction) IsTimeout() bool {
	return time.Since(s.UpdatedAt) > s.Timeout
}
