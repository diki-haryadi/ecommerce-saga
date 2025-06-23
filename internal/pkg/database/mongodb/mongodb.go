package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config holds MongoDB configuration
type Config struct {
	URI      string
	Database string
	Username string
	Password string
	Options  map[string]interface{}
}

// Client represents a MongoDB client
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	config   *Config
}

// NewClient creates a new MongoDB client
func NewClient(config *Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client options
	clientOptions := options.Client().ApplyURI(config.URI)
	if config.Username != "" && config.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.Username,
			Password: config.Password,
		})
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{
		client:   client,
		database: client.Database(config.Database),
		config:   config,
	}, nil
}

// Close closes the MongoDB connection
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// Database returns the MongoDB database
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Collection returns a MongoDB collection
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}

// IsHealthy checks if the MongoDB connection is healthy
func (c *Client) IsHealthy(ctx context.Context) bool {
	err := c.client.Ping(ctx, readpref.Primary())
	return err == nil
}

// WithTransaction executes the given function within a transaction
func (c *Client) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := c.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start MongoDB session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// EnsureIndexes ensures that the given indexes exist for a collection
func (c *Client) EnsureIndexes(ctx context.Context, collection string, indexes []mongo.IndexModel) error {
	_, err := c.Collection(collection).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}
	return nil
}
