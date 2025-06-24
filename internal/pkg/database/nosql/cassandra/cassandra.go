package cassandra

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gocql/gocql"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type cassandraDB struct {
	config  *database.Config
	session *gocql.Session
	mu      sync.RWMutex
}

// NewCassandraDB creates a new Cassandra database instance
func NewCassandraDB(config *database.Config) database.Database {
	return &cassandraDB{
		config: config,
	}
}

func (db *cassandraDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.session != nil {
		return nil
	}

	// Create cluster config
	cluster := gocql.NewCluster(db.config.Host)
	cluster.Port = db.config.Port
	cluster.Keyspace = db.config.Database
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = db.config.ConnectTimeout
	cluster.ConnectTimeout = db.config.ConnectTimeout
	cluster.NumConns = db.config.MaxConnections

	if db.config.Username != "" && db.config.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: db.config.Username,
			Password: db.config.Password,
		}
	}

	// Add additional options
	for key, value := range db.config.AdditionalOpts {
		switch key {
		case "consistency":
			if cons := gocql.ParseConsistency(value); cons != gocql.Any {
				cluster.Consistency = cons
			}
		case "retry_policy":
			if value == "exponential" {
				cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
					Min:        time.Second,
					Max:        time.Minute,
					NumRetries: db.config.MaxRetries,
				}
			}
		}
	}

	// Create session
	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("error creating cassandra session: %w", err)
	}

	db.session = session
	return nil
}

func (db *cassandraDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.session != nil {
		db.session.Close()
		db.session = nil
	}
	return nil
}

func (db *cassandraDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.session == nil {
		return fmt.Errorf("database not connected")
	}

	return db.session.Query("SELECT release_version FROM system.local").Exec()
}

func (db *cassandraDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.session != nil && !db.session.Closed()
}

func (db *cassandraDB) GetConfig() *database.Config {
	return db.config
}

func (db *cassandraDB) Begin(ctx context.Context) error {
	// Cassandra doesn't support traditional transactions
	// but we can use lightweight transactions (LWT) in queries
	return nil
}

func (db *cassandraDB) Commit(ctx context.Context) error {
	// Cassandra doesn't support traditional transactions
	return nil
}

func (db *cassandraDB) Rollback(ctx context.Context) error {
	// Cassandra doesn't support traditional transactions
	return nil
}

func (db *cassandraDB) Execute(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.session == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Create and execute query
	q := db.session.Query(query, args...)
	return nil, q.Exec()
}

func (db *cassandraDB) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.session == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Create and execute query
	q := db.session.Query(query, args...)
	return q.Iter(), nil
}
