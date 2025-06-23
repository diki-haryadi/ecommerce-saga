package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn *RabbitMQConnection
}

func NewPublisher(conn *RabbitMQConnection) *Publisher {
	return &Publisher{
		conn: conn,
	}
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Publish message
	return p.conn.Channel().Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Publisher) DeclareExchange(exchange string, kind string) error {
	return p.conn.Channel().ExchangeDeclare(
		exchange, // name
		kind,     // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}

func (p *Publisher) DeclareQueue(name string) error {
	_, err := p.conn.Channel().QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (p *Publisher) BindQueue(queue, exchange, routingKey string) error {
	return p.conn.Channel().QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	)
}
