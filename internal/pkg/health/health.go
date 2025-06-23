package health

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/logger"
)

// Status represents the health status
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// Component represents a health check component
type Component struct {
	Name      string                 `json:"name"`
	Status    Status                 `json:"status"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     string                 `json:"error,omitempty"`
	LastCheck time.Time              `json:"last_check"`
}

// Checker defines the interface for health checks
type Checker interface {
	Check(ctx context.Context) (*Component, error)
}

// Health manages health checks for the application
type Health struct {
	checkers map[string]Checker
	mutex    sync.RWMutex
	cache    map[string]*Component
	ttl      time.Duration
}

// NewHealth creates a new health manager
func NewHealth(ttl time.Duration) *Health {
	return &Health{
		checkers: make(map[string]Checker),
		cache:    make(map[string]*Component),
		ttl:      ttl,
	}
}

// RegisterChecker registers a health checker
func (h *Health) RegisterChecker(name string, checker Checker) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.checkers[name] = checker
}

// UnregisterChecker removes a health checker
func (h *Health) UnregisterChecker(name string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.checkers, name)
	delete(h.cache, name)
}

// CheckHealth performs health checks on all registered components
func (h *Health) CheckHealth(ctx context.Context) map[string]*Component {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	results := make(map[string]*Component)
	now := time.Now()

	for name, checker := range h.checkers {
		// Check cache first
		if cached, ok := h.cache[name]; ok {
			if now.Sub(cached.LastCheck) < h.ttl {
				results[name] = cached
				continue
			}
		}

		// Perform health check
		component, err := checker.Check(ctx)
		if err != nil {
			logger.Error("health check failed",
				zap.String("component", name),
				zap.Error(err),
			)
			results[name] = &Component{
				Name:      name,
				Status:    StatusDown,
				Error:     err.Error(),
				LastCheck: now,
			}
		} else {
			component.LastCheck = now
			results[name] = component
		}

		// Update cache
		h.cache[name] = results[name]
	}

	return results
}

// IsHealthy returns true if all components are healthy
func (h *Health) IsHealthy(ctx context.Context) bool {
	results := h.CheckHealth(ctx)
	for _, component := range results {
		if component.Status == StatusDown {
			return false
		}
	}
	return true
}

// GetStatus returns the overall system status
type SystemStatus struct {
	Status     Status                `json:"status"`
	Components map[string]*Component `json:"components"`
	Timestamp  time.Time             `json:"timestamp"`
}

func (h *Health) GetStatus(ctx context.Context) *SystemStatus {
	components := h.CheckHealth(ctx)
	status := StatusUp

	for _, component := range components {
		if component.Status == StatusDown {
			status = StatusDown
			break
		}
	}

	return &SystemStatus{
		Status:     status,
		Components: components,
		Timestamp:  time.Now(),
	}
}
