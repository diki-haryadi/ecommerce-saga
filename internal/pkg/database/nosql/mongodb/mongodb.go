package mongodb

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type mongoDB struct {
	config *database.Config
	client *mongo.Client
	mu     sync.RWMutex
}

// NewMongoDB creates a new MongoDB instance
func NewMongoDB(config *database.Config) database.Database {
	return &mongoDB{
		config: config,
	}
}

func (db *mongoDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		return nil
	}

	// Create MongoDB connection URI
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		db.config.Username,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		db.config.Database,
	)

	// Add additional options
	if len(db.config.AdditionalOpts) > 0 {
		uri += "?"
		for key, value := range db.config.AdditionalOpts {
			uri += fmt.Sprintf("%s=%s&", key, value)
		}
		uri = uri[:len(uri)-1] // Remove trailing &
	}

	// Configure client options
	opts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(uint64(db.config.MaxConnections)).
		SetMaxConnIdleTime(db.config.MaxIdleTime).
		SetConnectTimeout(db.config.ConnectTimeout)

	// Create client
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return fmt.Errorf("error connecting to mongodb: %w", err)
	}

	// Ping database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("error pinging mongodb: %w", err)
	}

	db.client = client
	return nil
}

func (db *mongoDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		if err := db.client.Disconnect(ctx); err != nil {
			return fmt.Errorf("error disconnecting from mongodb: %w", err)
		}
		db.client = nil
	}
	return nil
}

func (db *mongoDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return fmt.Errorf("database not connected")
	}
	return db.client.Ping(ctx, readpref.Primary())
}

func (db *mongoDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.client != nil
}

func (db *mongoDB) GetConfig() *database.Config {
	return db.config
}

func (db *mongoDB) Begin(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return fmt.Errorf("database not connected")
	}

	session, err := db.client.StartSession()
	if err != nil {
		return fmt.Errorf("error starting mongodb session: %w", err)
	}

	if err := session.StartTransaction(); err != nil {
		return fmt.Errorf("error starting mongodb transaction: %w", err)
	}

	// Store session in context
	ctx = context.WithValue(ctx, "session", session)
	return nil
}

func (db *mongoDB) Commit(ctx context.Context) error {
	session, ok := ctx.Value("session").(mongo.Session)
	if !ok {
		return fmt.Errorf("no session in context")
	}
	defer session.EndSession(ctx)
	return session.CommitTransaction(ctx)
}

func (db *mongoDB) Rollback(ctx context.Context) error {
	session, ok := ctx.Value("session").(mongo.Session)
	if !ok {
		return fmt.Errorf("no session in context")
	}
	defer session.EndSession(ctx)
	return session.AbortTransaction(ctx)
}

func (db *mongoDB) Execute(ctx context.Context, collection string, operations ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return nil, fmt.Errorf("database not connected")
	}

	coll := db.client.Database(db.config.Database).Collection(collection)

	// Check if we're in a transaction
	if session, ok := ctx.Value("session").(mongo.Session); ok {
		return session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			return executeOperation(sessCtx, coll, operations[0])
		})
	}

	return executeOperation(ctx, coll, operations[0])
}

func (db *mongoDB) Query(ctx context.Context, collection string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return nil, fmt.Errorf("database not connected")
	}

	coll := db.client.Database(db.config.Database).Collection(collection)
	filter := args[0]

	// Check if we're in a transaction
	if session, ok := ctx.Value("session").(mongo.Session); ok {
		return session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
			return coll.Find(sessCtx, filter)
		})
	}

	return coll.Find(ctx, filter)
}

func executeOperation(ctx context.Context, coll *mongo.Collection, operation interface{}) (interface{}, error) {
	switch op := operation.(type) {
	case mongo.InsertOneModel:
		return coll.InsertOne(ctx, op.Document)
	case mongo.UpdateOneModel:
		return coll.UpdateOne(ctx, op.Filter, op.Update)
	case mongo.DeleteOneModel:
		return coll.DeleteOne(ctx, op.Filter)
	default:
		return nil, fmt.Errorf("unsupported operation type")
	}
}
