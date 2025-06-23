package nats

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/messaging/types"
)

type NATSBroker struct {
	conn     *nats.Conn
	config   *types.BrokerConfig
	handlers map[string]types.MessageHandler
	subs     map[string]*nats.Subscription
}

// NewNATSBroker creates a new NATS broker instance
func NewNATSBroker(config *types.BrokerConfig) *NATSBroker {
	return &NATSBroker{
		config:   config,
		handlers: make(map[string]types.MessageHandler),
		subs:     make(map[string]*nats.Subscription),
	}
}

// Connect establishes a connection to NATS
func (n *NATSBroker) Connect(ctx context.Context) error {
	// Configure NATS options
	opts := []nats.Option{
		nats.Name("ecommerce-saga"),
		nats.ReconnectWait(time.Second * 10),
		nats.MaxReconnects(-1),
	}

	// Add credentials if provided
	if n.config.Username != "" && n.config.Password != "" {
		opts = append(opts, nats.UserInfo(n.config.Username, n.config.Password))
	}

	// Connect to NATS
	conn, err := nats.Connect(fmt.Sprintf("nats://%s:%s", n.config.Host, n.config.Port), opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	n.conn = conn

	return nil
}

// Close closes the NATS connection
func (n *NATSBroker) Close() error {
	for _, sub := range n.subs {
		sub.Unsubscribe()
	}
	if n.conn != nil {
		n.conn.Close()
	}
	return nil
}

// Publish publishes a message to a NATS subject
func (n *NATSBroker) Publish(ctx context.Context, subject string, msg *types.Message) error {
	if msg.ID == "" {
		msg.ID = uuid.New().String()
	}
	if msg.PublishedAt == 0 {
		msg.PublishedAt = time.Now().UnixNano()
	}

	// Create NATS message
	natsMsg := &nats.Msg{
		Subject: subject,
		Data:    msg.Payload,
		Header:  nats.Header{},
	}

	// Add headers
	for k, v := range msg.Headers {
		natsMsg.Header.Add(k, v)
	}

	// Add metadata
	natsMsg.Header.Add("message_id", msg.ID)
	natsMsg.Header.Add("published_at", fmt.Sprintf("%d", msg.PublishedAt))

	err := n.conn.PublishMsg(natsMsg)
	if err != nil {
		return fmt.Errorf("failed to publish message to NATS: %w", err)
	}

	return nil
}

// Subscribe subscribes to a NATS subject
func (n *NATSBroker) Subscribe(ctx context.Context, subject string, handler types.MessageHandler) error {
	sub, err := n.conn.Subscribe(subject, func(msg *nats.Msg) {
		message := &types.Message{
			ID:      msg.Header.Get("message_id"),
			Topic:   msg.Subject,
			Payload: msg.Data,
			Headers: make(map[string]string),
		}

		// Extract headers
		for k := range msg.Header {
			message.Headers[k] = msg.Header.Get(k)
		}

		// Extract published timestamp
		if ts := msg.Header.Get("published_at"); ts != "" {
			if publishedAt, err := strconv.ParseInt(ts, 10, 64); err == nil {
				message.PublishedAt = publishedAt
			} else {
				message.PublishedAt = time.Now().UnixNano()
			}
		}

		if err := handler(ctx, message); err != nil {
			// TODO: Implement error handling and retry logic
			fmt.Printf("Error processing message: %v\n", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to NATS subject: %w", err)
	}

	n.subs[subject] = sub
	n.handlers[subject] = handler

	return nil
}

// Unsubscribe removes a subscription from a NATS subject
func (n *NATSBroker) Unsubscribe(ctx context.Context, subject string) error {
	if sub, ok := n.subs[subject]; ok {
		if err := sub.Unsubscribe(); err != nil {
			return fmt.Errorf("failed to unsubscribe from NATS subject: %w", err)
		}
		delete(n.subs, subject)
	}
	delete(n.handlers, subject)
	return nil
}

// IsHealthy checks if the NATS connection is healthy
func (n *NATSBroker) IsHealthy(ctx context.Context) bool {
	return n.conn != nil && n.conn.IsConnected()
}
