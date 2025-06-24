package database

import (
	"context"
	"time"
)

// Config represents common database configuration
type Config struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Database       string
	MaxConnections int
	ConnectTimeout time.Duration
	MaxIdleTime    time.Duration
	MaxRetries     int
	RetryInterval  time.Duration
	SSLMode        string
	AdditionalOpts map[string]string
}

// Connection defines the common interface for all database connections
type Connection interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context) error
	IsConnected() bool
	GetConfig() *Config
}

// Transaction defines the common interface for database transactions
type Transaction interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// QueryExecutor defines the common interface for executing queries
type QueryExecutor interface {
	Execute(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	Query(ctx context.Context, query string, args ...interface{}) (interface{}, error)
}

// Database combines all database operations
type Database interface {
	Connection
	Transaction
	QueryExecutor
}

// Factory defines the interface for creating database connections
type Factory interface {
	CreateConnection(config *Config) (Database, error)
}
