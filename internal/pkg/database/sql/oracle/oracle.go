package oracle

import (
	"context"
	"fmt"
	"sync"

	"github.com/godror/godror"
	"github.com/jmoiron/sqlx"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type oracleDB struct {
	config *database.Config
	db     *sqlx.DB
	mu     sync.RWMutex
}

// NewOracleDB creates a new Oracle database instance
func NewOracleDB(config *database.Config) database.Database {
	return &oracleDB{
		config: config,
	}
}

func (db *oracleDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.db != nil {
		return nil
	}

	// Create Oracle connection parameters
	params := godror.ConnectionParams{
		StandaloneConnection: true,
		ConnectString:        fmt.Sprintf("%s:%d/%s", db.config.Host, db.config.Port, db.config.Database),
		Username:             db.config.Username,
		Password:             godror.NewPassword(db.config.Password),
	}

	// Add additional options
	for key, value := range db.config.AdditionalOpts {
		params.SetSessionParamOnInit(key, value)
	}

	// Connect to database
	sqlxDB, err := sqlx.ConnectContext(ctx, "godror", params.StringWithPassword())
	if err != nil {
		return fmt.Errorf("error connecting to oracle: %w", err)
	}

	// Configure connection pool
	sqlxDB.SetMaxOpenConns(db.config.MaxConnections)
	sqlxDB.SetConnMaxIdleTime(db.config.MaxIdleTime)

	db.db = sqlxDB
	return nil
}

func (db *oracleDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.db != nil {
		if err := db.db.Close(); err != nil {
			return fmt.Errorf("error disconnecting from oracle: %w", err)
		}
		db.db = nil
	}
	return nil
}

func (db *oracleDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.db == nil {
		return fmt.Errorf("database not connected")
	}
	return db.db.PingContext(ctx)
}

func (db *oracleDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.db != nil
}

func (db *oracleDB) GetConfig() *database.Config {
	return db.config
}

func (db *oracleDB) Begin(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.db == nil {
		return fmt.Errorf("database not connected")
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	// Store transaction in context
	ctx = context.WithValue(ctx, "tx", tx)
	return nil
}

func (db *oracleDB) Commit(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Commit()
}

func (db *oracleDB) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Rollback()
}

func (db *oracleDB) Execute(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Check if we're in a transaction
	if tx, ok := ctx.Value("tx").(*sqlx.Tx); ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return db.db.ExecContext(ctx, query, args...)
}

func (db *oracleDB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Check if we're in a transaction
	if tx, ok := ctx.Value("tx").(*sqlx.Tx); ok {
		return tx.QueryxContext(ctx, query, args...)
	}

	return db.db.QueryxContext(ctx, query, args...)
}
