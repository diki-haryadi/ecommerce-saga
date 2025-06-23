package messaging

import (
	"context"
)

// MessageBroker defines the interface for message broker operations
type MessageBroker interface {
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe subscribes to a topic with a message handler
	Subscribe(topic string, handler MessageHandler) error

	// Close closes the message broker connection
	Close() error
}

// MessageHandler is a function that handles received messages
type MessageHandler func(ctx context.Context, message []byte) error
