package factory

import (
	"fmt"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/broker/kafka"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/broker/nats"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/broker/nsq"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/config"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/messaging/types"
)

// BrokerType represents the type of message broker
type BrokerType string

const (
	KafkaBroker BrokerType = "kafka"
	NSQBroker   BrokerType = "nsq"
	NATSBroker  BrokerType = "nats"
)

// BrokerFactory creates message brokers
type BrokerFactory interface {
	CreateBroker(instanceName string) (types.MessageBroker, error)
}

// KafkaBrokerFactory creates Kafka brokers
type KafkaBrokerFactory struct {
	config *config.BrokersConfig
}

// NSQBrokerFactory creates NSQ brokers
type NSQBrokerFactory struct {
	config *config.BrokersConfig
}

// NATSBrokerFactory creates NATS brokers
type NATSBrokerFactory struct {
	config *config.BrokersConfig
}

// NewBrokerFactory creates a new broker factory based on the broker type
func NewBrokerFactory(brokerType BrokerType, config *config.BrokersConfig) BrokerFactory {
	switch brokerType {
	case KafkaBroker:
		return &KafkaBrokerFactory{config: config}
	case NSQBroker:
		return &NSQBrokerFactory{config: config}
	case NATSBroker:
		return &NATSBrokerFactory{config: config}
	default:
		return nil
	}
}

// CreateBroker creates a new Kafka broker instance
func (f *KafkaBrokerFactory) CreateBroker(instanceName string) (types.MessageBroker, error) {
	brokerConfig, ok := f.config.Kafka[instanceName]
	if !ok {
		return nil, fmt.Errorf("kafka broker instance '%s' not found in config", instanceName)
	}

	if !brokerConfig.Enabled {
		return nil, fmt.Errorf("kafka broker instance '%s' is disabled", instanceName)
	}

	return kafka.NewKafkaBroker(&types.BrokerConfig{
		Host:     brokerConfig.Host,
		Port:     brokerConfig.Port,
		Username: brokerConfig.Username,
		Password: brokerConfig.Password,
		Options:  brokerConfig.Options,
	}), nil
}

// CreateBroker creates a new NSQ broker instance
func (f *NSQBrokerFactory) CreateBroker(instanceName string) (types.MessageBroker, error) {
	brokerConfig, ok := f.config.NSQ[instanceName]
	if !ok {
		return nil, fmt.Errorf("nsq broker instance '%s' not found in config", instanceName)
	}

	if !brokerConfig.Enabled {
		return nil, fmt.Errorf("nsq broker instance '%s' is disabled", instanceName)
	}

	return nsq.NewNSQBroker(&types.BrokerConfig{
		Host:     brokerConfig.Host,
		Port:     brokerConfig.Port,
		Username: brokerConfig.Username,
		Password: brokerConfig.Password,
		Options:  brokerConfig.Options,
	}), nil
}

// CreateBroker creates a new NATS broker instance
func (f *NATSBrokerFactory) CreateBroker(instanceName string) (types.MessageBroker, error) {
	brokerConfig, ok := f.config.NATS[instanceName]
	if !ok {
		return nil, fmt.Errorf("nats broker instance '%s' not found in config", instanceName)
	}

	if !brokerConfig.Enabled {
		return nil, fmt.Errorf("nats broker instance '%s' is disabled", instanceName)
	}

	return nats.NewNATSBroker(&types.BrokerConfig{
		Host:     brokerConfig.Host,
		Port:     brokerConfig.Port,
		Username: brokerConfig.Username,
		Password: brokerConfig.Password,
		Options:  brokerConfig.Options,
	}), nil
}
