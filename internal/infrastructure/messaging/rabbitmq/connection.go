package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQConnection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQConnection{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQConnection) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQConnection) Close() {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}
}
