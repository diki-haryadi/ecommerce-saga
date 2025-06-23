package messaging

import (
	"context"
	"fmt"
	"sync"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/config"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/messaging/factory"
)

// BrokerManager manages multiple message broker instances
type BrokerManager struct {
	config    *config.BrokersConfig
	brokers   map[string]MessageBroker
	mutex     sync.RWMutex
	factories map[factory.BrokerType]factory.BrokerFactory
}

// NewBrokerManager creates a new broker manager
func NewBrokerManager(config *config.BrokersConfig) *BrokerManager {
	manager := &BrokerManager{
		config:    config,
		brokers:   make(map[string]MessageBroker),
		factories: make(map[factory.BrokerType]factory.BrokerFactory),
	}

	// Initialize factories for each broker type
	manager.factories[factory.KafkaBroker] = factory.NewBrokerFactory(factory.KafkaBroker, config)
	manager.factories[factory.NSQBroker] = factory.NewBrokerFactory(factory.NSQBroker, config)
	manager.factories[factory.NATSBroker] = factory.NewBrokerFactory(factory.NATSBroker, config)

	return manager
}

// GetBroker returns a broker instance by name and type
func (m *BrokerManager) GetBroker(brokerType factory.BrokerType, instanceName string) (MessageBroker, error) {
	m.mutex.RLock()
	broker, exists := m.brokers[getBrokerKey(brokerType, instanceName)]
	m.mutex.RUnlock()

	if exists {
		return broker, nil
	}

	return m.createBroker(brokerType, instanceName)
}

// GetDefaultBroker returns the default broker instance
func (m *BrokerManager) GetDefaultBroker() (MessageBroker, error) {
	if m.config.Default == "" {
		return nil, fmt.Errorf("no default broker configured")
	}

	// Parse default broker string (format: "type:instance")
	brokerType, instanceName, err := parseBrokerString(m.config.Default)
	if err != nil {
		return nil, err
	}

	return m.GetBroker(brokerType, instanceName)
}

// createBroker creates a new broker instance
func (m *BrokerManager) createBroker(brokerType factory.BrokerType, instanceName string) (MessageBroker, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if broker was created while waiting for lock
	if broker, exists := m.brokers[getBrokerKey(brokerType, instanceName)]; exists {
		return broker, nil
	}

	// Get factory for broker type
	brokerFactory, exists := m.factories[brokerType]
	if !exists {
		return nil, fmt.Errorf("unsupported broker type: %s", brokerType)
	}

	// Create broker instance
	broker, err := brokerFactory.CreateBroker(instanceName)
	if err != nil {
		return nil, err
	}

	// Connect to broker
	if err := broker.Connect(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}

	// Store broker instance
	m.brokers[getBrokerKey(brokerType, instanceName)] = broker

	return broker, nil
}

// Close closes all broker connections
func (m *BrokerManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var errs []error
	for _, broker := range m.brokers {
		if err := broker.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close some broker connections: %v", errs)
	}

	return nil
}

// getBrokerKey returns a unique key for a broker instance
func getBrokerKey(brokerType factory.BrokerType, instanceName string) string {
	return fmt.Sprintf("%s:%s", brokerType, instanceName)
}

// parseBrokerString parses a broker string in the format "type:instance"
func parseBrokerString(s string) (factory.BrokerType, string, error) {
	var brokerType, instanceName string
	_, err := fmt.Sscanf(s, "%s:%s", &brokerType, &instanceName)
	if err != nil {
		return "", "", fmt.Errorf("invalid broker string format (expected 'type:instance'): %w", err)
	}

	return factory.BrokerType(brokerType), instanceName, nil
}
