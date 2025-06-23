package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, topic string, message []byte) error {
	// Ensure exchange exists
	err := r.channel.ExchangeDeclare(
		"saga",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Publish message
	err = r.channel.PublishWithContext(
		ctx,
		"saga", // exchange
		topic,  // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Subscribe(topic string, handler MessageHandler) error {
	// Ensure exchange exists
	err := r.channel.ExchangeDeclare(
		"saga",  // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Create queue
	q, err := r.channel.QueueDeclare(
		"",    // name (empty = auto-generated)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = r.channel.QueueBind(
		q.Name, // queue name
		topic,  // routing key
		"saga", // exchange
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Start consuming messages
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	// Handle messages in a goroutine
	go func() {
		for msg := range msgs {
			if err := handler(context.Background(), msg.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
