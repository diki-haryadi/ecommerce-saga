package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type MessageHandler func(ctx context.Context, message []byte) error

type Consumer struct {
	conn     *RabbitMQConnection
	handlers map[string]MessageHandler
}

func NewConsumer(conn *RabbitMQConnection) *Consumer {
	return &Consumer{
		conn:     conn,
		handlers: make(map[string]MessageHandler),
	}
}

func (c *Consumer) RegisterHandler(queue string, handler MessageHandler) {
	c.handlers[queue] = handler
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	for queue, handler := range c.handlers {
		// Create a new channel for each consumer
		ch, err := c.conn.conn.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		// Set QoS
		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		if err != nil {
			return err
		}

		// Start consuming
		msgs, err := ch.Consume(
			queue, // queue
			"",    // consumer
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			return err
		}

		go func(deliveries <-chan amqp.Delivery, handler MessageHandler) {
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-deliveries:
					if !ok {
						return
					}

					// Handle message
					err := handler(ctx, d.Body)
					if err != nil {
						log.Printf("Error handling message: %v", err)
						d.Nack(false, true) // Negative acknowledgement, requeue
						continue
					}

					d.Ack(false) // Acknowledge message
				}
			}
		}(msgs, handler)
	}

	return nil
}

// Helper function to unmarshal message
func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
