package types

import "context"

// Message represents a message to be sent/received
type Message struct {
	ID          string                 `json:"id"`
	Topic       string                 `json:"topic"`
	Payload     []byte                 `json:"payload"`
	Headers     map[string]string      `json:"headers,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	PublishedAt int64                  `json:"published_at"`
}

// MessageHandler is a function that processes a message
type MessageHandler func(context.Context, *Message) error

// MessageBroker defines the interface for message brokers
type MessageBroker interface {
	// Connect establishes a connection to the message broker
	Connect(ctx context.Context) error
	// Close closes the connection
	Close() error
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, msg *Message) error
	// Subscribe subscribes to a topic
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	// Unsubscribe removes a subscription
	Unsubscribe(ctx context.Context, topic string) error
	// IsHealthy checks if the broker connection is healthy
	IsHealthy(ctx context.Context) bool
}

// BrokerConfig holds the configuration for message brokers
type BrokerConfig struct {
	Type     string                 `json:"type"`     // rabbitmq, kafka, nsq, nats
	Host     string                 `json:"host"`     // Broker host
	Port     string                 `json:"port"`     // Broker port
	Username string                 `json:"username"` // Authentication username
	Password string                 `json:"password"` // Authentication password
	Options  map[string]interface{} `json:"options"`  // Additional broker-specific options
}
