package nsq

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/messaging"
)

type NSQBroker struct {
	producer *nsq.Producer
	consumer *nsq.Consumer
	config   *messaging.BrokerConfig
	handlers map[string]messaging.MessageHandler
}

// NewNSQBroker creates a new NSQ broker instance
func NewNSQBroker(config *messaging.BrokerConfig) *NSQBroker {
	return &NSQBroker{
		config:   config,
		handlers: make(map[string]messaging.MessageHandler),
	}
}

// Connect establishes a connection to NSQ
func (n *NSQBroker) Connect(ctx context.Context) error {
	// Configure NSQ
	config := nsq.NewConfig()
	config.MaxInFlight = 10

	// Create producer
	producer, err := nsq.NewProducer(fmt.Sprintf("%s:%s", n.config.Host, n.config.Port), config)
	if err != nil {
		return fmt.Errorf("failed to create NSQ producer: %w", err)
	}
	n.producer = producer

	return nil
}

// Close closes the NSQ connections
func (n *NSQBroker) Close() error {
	if n.producer != nil {
		n.producer.Stop()
	}
	if n.consumer != nil {
		n.consumer.Stop()
	}
	return nil
}

// Publish publishes a message to an NSQ topic
func (n *NSQBroker) Publish(ctx context.Context, topic string, msg *messaging.Message) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.PublishedAt == 0 {
		msg.PublishedAt = time.Now().UnixNano()
	}

	err := n.producer.Publish(topic, msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to publish message to NSQ: %w", err)
	}

	return nil
}

// Subscribe subscribes to an NSQ topic
func (n *NSQBroker) Subscribe(ctx context.Context, topic string, handler messaging.MessageHandler) error {
	config := nsq.NewConfig()
	config.MaxInFlight = 10

	consumer, err := nsq.NewConsumer(topic, "default", config)
	if err != nil {
		return fmt.Errorf("failed to create NSQ consumer: %w", err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		message := &messaging.Message{
			ID:          uuid.New().String(),
			Topic:       topic,
			Payload:     msg.Body,
			PublishedAt: msg.Timestamp.UnixNano(),
			Headers:     make(map[string]string),
		}

		if err := handler(ctx, message); err != nil {
			// TODO: Implement error handling and retry logic
			fmt.Printf("Error processing message: %v\n", err)
			return err
		}

		return nil
	}))

	err = consumer.ConnectToNSQD(fmt.Sprintf("%s:%s", n.config.Host, n.config.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to NSQ: %w", err)
	}

	n.consumer = consumer
	n.handlers[topic] = handler

	return nil
}

// Unsubscribe removes a subscription from an NSQ topic
func (n *NSQBroker) Unsubscribe(ctx context.Context, topic string) error {
	if n.consumer != nil {
		n.consumer.Stop()
	}
	delete(n.handlers, topic)
	return nil
}

// IsHealthy checks if the NSQ connection is healthy
func (n *NSQBroker) IsHealthy(ctx context.Context) bool {
	if n.producer == nil {
		return false
	}

	// Ping the NSQ daemon
	err := n.producer.Ping()
	return err == nil
}
