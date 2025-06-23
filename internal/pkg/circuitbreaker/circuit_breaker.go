package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/logger"
)

// State represents the state of the circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Config holds circuit breaker configuration
type Config struct {
	Threshold      int           // Number of failures before opening
	Timeout        time.Duration // How long to wait before attempting recovery
	HalfOpenCalls  int           // Number of calls to allow in half-open state
	RetryInterval  time.Duration // Time between retries in half-open state
	FailureHandler func(error)   // Optional handler for failures
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config         *Config
	state          State
	failures       int
	lastFailure    time.Time
	halfOpenCalls  int
	mutex          sync.RWMutex
	failureHandler func(error)
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config *Config) *CircuitBreaker {
	if config.FailureHandler == nil {
		config.FailureHandler = func(err error) {
			logger.Error("circuit breaker failure", zap.Error(err))
		}
	}

	return &CircuitBreaker{
		config:         config,
		state:          StateClosed,
		failureHandler: config.FailureHandler,
	}
}

// Execute executes the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mutex.RLock()
	state := cb.state
	cb.mutex.RUnlock()

	switch state {
	case StateOpen:
		if time.Since(cb.lastFailure) > cb.config.Timeout {
			cb.mutex.Lock()
			cb.state = StateHalfOpen
			cb.halfOpenCalls = 0
			cb.mutex.Unlock()
		} else {
			return ErrCircuitOpen
		}

	case StateHalfOpen:
		cb.mutex.Lock()
		if cb.halfOpenCalls >= cb.config.HalfOpenCalls {
			cb.mutex.Unlock()
			return ErrCircuitOpen
		}
		cb.halfOpenCalls++
		cb.mutex.Unlock()
	}

	err := fn()
	if err != nil {
		cb.recordFailure(err)
		return err
	}

	cb.recordSuccess()
	return nil
}

// recordFailure records a failure and potentially opens the circuit
func (cb *CircuitBreaker) recordFailure(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.config.Threshold {
		cb.state = StateOpen
		cb.failures = 0
	}

	if cb.failureHandler != nil {
		cb.failureHandler(err)
	}
}

// recordSuccess records a success and potentially closes the circuit
func (cb *CircuitBreaker) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch cb.state {
	case StateHalfOpen:
		cb.state = StateClosed
		cb.failures = 0
		cb.halfOpenCalls = 0
	case StateClosed:
		cb.failures = 0
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to its initial state
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenCalls = 0
}
