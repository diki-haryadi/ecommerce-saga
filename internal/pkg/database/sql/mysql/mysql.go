package mysql

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type mysqlDB struct {
	config *database.Config
	db     *sqlx.DB
	mu     sync.RWMutex
}

// NewMySQLDB creates a new MySQL database instance
func NewMySQLDB(config *database.Config) database.Database {
	return &mysqlDB{
		config: config,
	}
}

func (db *mysqlDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.db != nil {
		return nil
	}

	// Add additional options
	mysqlConfig := mysql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", db.config.Host, db.config.Port),
		DBName:               db.config.Database,
		User:                 db.config.Username,
		Passwd:               db.config.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
		Timeout:              db.config.ConnectTimeout,
		MaxAllowedPacket:     4 << 20, // 4 MB
	}

	// Add additional options from config
	for key, value := range db.config.AdditionalOpts {
		mysqlConfig.Params[key] = value
	}

	// Connect to database
	sqlxDB, err := sqlx.ConnectContext(ctx, "mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return fmt.Errorf("error connecting to mysql: %w", err)
	}

	// Configure connection pool
	sqlxDB.SetMaxOpenConns(db.config.MaxConnections)
	sqlxDB.SetConnMaxIdleTime(db.config.MaxIdleTime)

	db.db = sqlxDB
	return nil
}

func (db *mysqlDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.db != nil {
		if err := db.db.Close(); err != nil {
			return fmt.Errorf("error disconnecting from mysql: %w", err)
		}
		db.db = nil
	}
	return nil
}

func (db *mysqlDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.db == nil {
		return fmt.Errorf("database not connected")
	}
	return db.db.PingContext(ctx)
}

func (db *mysqlDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.db != nil
}

func (db *mysqlDB) GetConfig() *database.Config {
	return db.config
}

func (db *mysqlDB) Begin(ctx context.Context) error {
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

func (db *mysqlDB) Commit(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Commit()
}

func (db *mysqlDB) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}
	return tx.Rollback()
}

func (db *mysqlDB) Execute(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
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

func (db *mysqlDB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
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
