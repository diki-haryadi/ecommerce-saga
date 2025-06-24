package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

// PoolConfig represents pool configuration
type PoolConfig struct {
	MinConnections int
	MaxConnections int
	IdleTimeout    time.Duration
	RetryInterval  time.Duration
	MaxRetries     int
}

// Pool manages a pool of database connections
type Pool struct {
	config      *PoolConfig
	factory     DatabaseFactory
	dbConfig    *database.Config
	connections chan database.Database
	mu          sync.RWMutex
	closed      bool
}

// NewPool creates a new connection pool
func NewPool(factory DatabaseFactory, dbConfig *database.Config, poolConfig *PoolConfig) (*Pool, error) {
	if poolConfig.MinConnections > poolConfig.MaxConnections {
		return nil, fmt.Errorf("min connections cannot be greater than max connections")
	}

	pool := &Pool{
		config:      poolConfig,
		factory:     factory,
		dbConfig:    dbConfig,
		connections: make(chan database.Database, poolConfig.MaxConnections),
	}

	// Initialize minimum connections
	for i := 0; i < poolConfig.MinConnections; i++ {
		conn, err := pool.createConnection()
		if err != nil {
			return nil, fmt.Errorf("error initializing pool: %w", err)
		}
		pool.connections <- conn
	}

	// Start connection monitor
	go pool.monitor()

	return pool, nil
}

// Get retrieves a connection from the pool
func (p *Pool) Get(ctx context.Context) (database.Database, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, fmt.Errorf("pool is closed")
	}
	p.mu.RUnlock()

	select {
	case conn := <-p.connections:
		if !conn.IsConnected() {
			// Connection is stale, create a new one
			if err := conn.Connect(ctx); err != nil {
				return nil, fmt.Errorf("error reconnecting: %w", err)
			}
		}
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// No connections available, try to create a new one
		if len(p.connections) < p.config.MaxConnections {
			conn, err := p.createConnection()
			if err != nil {
				return nil, fmt.Errorf("error creating new connection: %w", err)
			}
			return conn, nil
		}
		// Wait for an available connection
		select {
		case conn := <-p.connections:
			return conn, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Put returns a connection to the pool
func (p *Pool) Put(conn database.Database) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return fmt.Errorf("pool is closed")
	}
	p.mu.RUnlock()

	select {
	case p.connections <- conn:
		return nil
	default:
		// Pool is full, close the connection
		return conn.Disconnect(context.Background())
	}
}

// Close closes all connections in the pool
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	close(p.connections)

	var errs []error
	for conn := range p.connections {
		if err := conn.Disconnect(context.Background()); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}

// createConnection creates a new database connection
func (p *Pool) createConnection() (database.Database, error) {
	conn, err := p.factory.CreateDatabase(p.dbConfig)
	if err != nil {
		return nil, err
	}

	if err := conn.Connect(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

// monitor periodically checks connections and removes stale ones
func (p *Pool) monitor() {
	ticker := time.NewTicker(p.config.IdleTimeout)
	defer ticker.Stop()

	for {
		<-ticker.C

		p.mu.RLock()
		if p.closed {
			p.mu.RUnlock()
			return
		}
		p.mu.RUnlock()

		currentConns := len(p.connections)
		for i := 0; i < currentConns; i++ {
			select {
			case conn := <-p.connections:
				if !conn.IsConnected() {
					// Connection is stale, create a new one
					if newConn, err := p.createConnection(); err == nil {
						p.connections <- newConn
					}
				} else {
					p.connections <- conn
				}
			default:
				// No more connections to check
				break
			}
		}
	}
}
