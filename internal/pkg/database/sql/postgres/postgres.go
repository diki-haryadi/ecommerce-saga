package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type postgresDB struct {
	config *database.Config
	pool   *pgxpool.Pool
	mu     sync.RWMutex
}

// NewPostgresDB creates a new PostgreSQL database instance
func NewPostgresDB(config *database.Config) database.Database {
	return &postgresDB{
		config: config,
	}
}

func (db *postgresDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.pool != nil {
		return nil
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		db.config.Username,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		db.config.Database,
		db.config.SSLMode,
	)

	// Add additional options
	for key, value := range db.config.AdditionalOpts {
		dsn += fmt.Sprintf("&%s=%s", key, value)
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("error parsing postgres config: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = int32(db.config.MaxConnections)
	poolConfig.MaxConnLifetime = db.config.MaxIdleTime
	poolConfig.ConnConfig.ConnectTimeout = db.config.ConnectTimeout

	// Create connection pool
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %w", err)
	}

	db.pool = pool
	return nil
}

func (db *postgresDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.pool != nil {
		db.pool.Close()
		db.pool = nil
	}
	return nil
}

func (db *postgresDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.pool == nil {
		return fmt.Errorf("database not connected")
	}
	return db.pool.Ping(ctx)
}

func (db *postgresDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.pool != nil
}

func (db *postgresDB) GetConfig() *database.Config {
	return db.config
}

func (db *postgresDB) Begin(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.pool == nil {
		return fmt.Errorf("database not connected")
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	// Store transaction in context
	ctx = context.WithValue(ctx, "tx", tx)
	return nil
}

func (db *postgresDB) Commit(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(pgxpool.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Commit(ctx)
}

func (db *postgresDB) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(pgxpool.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Rollback(ctx)
}

func (db *postgresDB) Execute(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Check if we're in a transaction
	if tx, ok := ctx.Value("tx").(pgxpool.Tx); ok {
		return tx.Exec(ctx, query, args...)
	}

	return db.pool.Exec(ctx, query, args...)
}

func (db *postgresDB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.pool == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Check if we're in a transaction
	if tx, ok := ctx.Value("tx").(pgxpool.Tx); ok {
		return tx.Query(ctx, query, args...)
	}

	return db.pool.Query(ctx, query, args...)
}
