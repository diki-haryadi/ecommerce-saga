package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/messaging"
)

type KafkaBroker struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	config   *messaging.BrokerConfig
	handlers map[string]messaging.MessageHandler
}

// NewKafkaBroker creates a new Kafka broker instance
func NewKafkaBroker(config *messaging.BrokerConfig) *KafkaBroker {
	return &KafkaBroker{
		config:   config,
		handlers: make(map[string]messaging.MessageHandler),
	}
}

// Connect establishes a connection to Kafka
func (k *KafkaBroker) Connect(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Create producer
	brokers := []string{fmt.Sprintf("%s:%s", k.config.Host, k.config.Port)}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	k.producer = producer

	// Create consumer
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	k.consumer = consumer

	return nil
}

// Close closes the Kafka connections
func (k *KafkaBroker) Close() error {
	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}
	if err := k.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka consumer: %w", err)
	}
	return nil
}

// Publish publishes a message to a Kafka topic
func (k *KafkaBroker) Publish(ctx context.Context, topic string, msg *messaging.Message) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.PublishedAt == 0 {
		msg.PublishedAt = time.Now().UnixNano()
	}

	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(msg.ID),
		Value: sarama.ByteEncoder(msg.Payload),
		Headers: func() []sarama.RecordHeader {
			headers := make([]sarama.RecordHeader, 0, len(msg.Headers))
			for k, v := range msg.Headers {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
			return headers
		}(),
	})

	if err != nil {
		return fmt.Errorf("failed to publish message to Kafka: %w", err)
	}

	return nil
}

// Subscribe subscribes to a Kafka topic
func (k *KafkaBroker) Subscribe(ctx context.Context, topic string, handler messaging.MessageHandler) error {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}

	k.handlers[topic] = handler

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				message := &messaging.Message{
					ID:          string(msg.Key),
					Topic:       msg.Topic,
					Payload:     msg.Value,
					PublishedAt: msg.Timestamp.UnixNano(),
					Headers:     make(map[string]string),
				}

				// Extract headers
				for _, header := range msg.Headers {
					message.Headers[string(header.Key)] = string(header.Value)
				}

				if err := handler(ctx, message); err != nil {
					// TODO: Implement error handling and retry logic
					fmt.Printf("Error processing message: %v\n", err)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Unsubscribe removes a subscription from a Kafka topic
func (k *KafkaBroker) Unsubscribe(ctx context.Context, topic string) error {
	delete(k.handlers, topic)
	return nil
}

// IsHealthy checks if the Kafka connection is healthy
func (k *KafkaBroker) IsHealthy(ctx context.Context) bool {
	// Check if producer and consumer are initialized
	if k.producer == nil || k.consumer == nil {
		return false
	}

	// Try to list topics as a health check
	_, err := k.consumer.Topics()
	return err == nil
}
