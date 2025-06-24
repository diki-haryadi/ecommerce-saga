package aerospike

import (
	"context"
	"fmt"
	"sync"
	"time"

	as "github.com/aerospike/aerospike-client-go/v6"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

type aerospikeDB struct {
	config *database.Config
	client *as.Client
	mu     sync.RWMutex
}

// NewAerospike creates a new Aerospike database instance
func NewAerospike(config *database.Config) database.Database {
	return &aerospikeDB{
		config: config,
	}
}

func (db *aerospikeDB) Connect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		return nil
	}

	// Create client policy
	policy := as.NewClientPolicy()
	policy.Timeout = db.config.ConnectTimeout
	policy.ConnectionQueueSize = db.config.MaxConnections
	policy.User = db.config.Username
	policy.Password = db.config.Password

	// Add additional options
	for key, _ := range db.config.AdditionalOpts {
		switch key {
		case "max_retries":
			if maxRetries := db.config.MaxRetries; maxRetries > 0 {
				policy.MaxRetries = maxRetries
				policy.RetryOnTimeout = true
			}
		case "idle_timeout":
			if timeout := db.config.MaxIdleTime; timeout > 0 {
				policy.IdleTimeout = timeout
			}
		}
	}

	// Create client
	client, err := as.NewClientWithPolicy(policy, db.config.Host, db.config.Port)
	if err != nil {
		return fmt.Errorf("error connecting to aerospike: %w", err)
	}

	db.client = client
	return nil
}

func (db *aerospikeDB) Disconnect(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.client != nil {
		db.client.Close()
		db.client = nil
	}
	return nil
}

func (db *aerospikeDB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return fmt.Errorf("database not connected")
	}

	if !db.client.IsConnected() {
		return fmt.Errorf("database connection lost")
	}
	return nil
}

func (db *aerospikeDB) IsConnected() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.client != nil && db.client.IsConnected()
}

func (db *aerospikeDB) GetConfig() *database.Config {
	return db.config
}

func (db *aerospikeDB) Begin(ctx context.Context) error {
	// Aerospike supports atomic operations but not traditional transactions
	return nil
}

func (db *aerospikeDB) Commit(ctx context.Context) error {
	// Aerospike supports atomic operations but not traditional transactions
	return nil
}

func (db *aerospikeDB) Rollback(ctx context.Context) error {
	// Aerospike supports atomic operations but not traditional transactions
	return nil
}

func (db *aerospikeDB) Execute(ctx context.Context, operation string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return nil, fmt.Errorf("database not connected")
	}

	if len(args) < 3 {
		return nil, fmt.Errorf("insufficient arguments for operation")
	}

	// Extract common parameters
	namespace, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("namespace must be a string")
	}

	setName, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("set name must be a string")
	}

	key, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	asKey, err := as.NewKey(namespace, setName, key)
	if err != nil {
		return nil, fmt.Errorf("error creating aerospike key: %w", err)
	}

	switch operation {
	case "Put":
		if len(args) < 4 {
			return nil, fmt.Errorf("missing bins for Put operation")
		}
		bins, ok := args[3].(as.BinMap)
		if !ok {
			return nil, fmt.Errorf("bins must be an aerospike.BinMap")
		}
		return nil, db.client.Put(nil, asKey, bins)

	case "Get":
		record, err := db.client.Get(nil, asKey)
		if err != nil {
			return nil, err
		}
		return record, nil

	case "Delete":
		existed, err := db.client.Delete(nil, asKey)
		if err != nil {
			return nil, err
		}
		return existed, nil

	case "Touch":
		ttl, ok := args[3].(time.Duration)
		if !ok {
			return nil, fmt.Errorf("TTL must be a time.Duration")
		}
		writePolicy := as.NewWritePolicy(0, uint32(ttl.Seconds()))
		return nil, db.client.Touch(writePolicy, asKey)

	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

func (db *aerospikeDB) Query(ctx context.Context, namespace string, args ...interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.client == nil {
		return nil, fmt.Errorf("database not connected")
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("insufficient arguments for query")
	}

	setName, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("set name must be a string")
	}

	stmt := as.NewStatement(namespace, setName)

	// Add filters if provided
	if len(args) > 1 {
		if filters, ok := args[1].([]as.Filter); ok {
			for _, filter := range filters {
				stmt.SetFilter(&filter)
			}
		}
	}

	// Execute query
	recordset, err := db.client.Query(nil, stmt)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return recordset, nil
}
