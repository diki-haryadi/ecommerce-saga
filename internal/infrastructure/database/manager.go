package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/nosql/aerospike"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/nosql/cassandra"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/nosql/dynamodb"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/nosql/mongodb"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/sql/mysql"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/sql/oracle"
	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database/sql/postgres"
)

// DatabaseType represents the type of database
type DatabaseType string

const (
	// SQL Databases
	TypePostgres DatabaseType = "postgres"
	TypeMySQL    DatabaseType = "mysql"
	TypeOracle   DatabaseType = "oracle"

	// NoSQL Databases
	TypeMongoDB   DatabaseType = "mongodb"
	TypeDynamoDB  DatabaseType = "dynamodb"
	TypeCassandra DatabaseType = "cassandra"
	TypeAerospike DatabaseType = "aerospike"
)

// Manager manages database connections using the Singleton pattern
type Manager struct {
	connections map[string]database.Database
	mu          sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetInstance returns the singleton instance of Manager
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{
			connections: make(map[string]database.Database),
		}
	})
	return instance
}

// DatabaseFactory creates database instances based on type
type DatabaseFactory interface {
	CreateDatabase(config *database.Config) (database.Database, error)
}

// PostgresFactory implements DatabaseFactory for PostgreSQL
type PostgresFactory struct{}

func (f *PostgresFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return postgres.NewPostgresDB(config), nil
}

// MongoDBFactory implements DatabaseFactory for MongoDB
type MongoDBFactory struct{}

func (f *MongoDBFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return mongodb.NewMongoDB(config), nil
}

// GetFactory returns the appropriate database factory
func GetFactory(dbType DatabaseType) (DatabaseFactory, error) {
	switch dbType {
	// SQL Databases
	case TypePostgres:
		return &PostgresFactory{}, nil
	case TypeMySQL:
		return &MySQLFactory{}, nil
	case TypeOracle:
		return &OracleFactory{}, nil

	// NoSQL Databases
	case TypeMongoDB:
		return &MongoDBFactory{}, nil
	case TypeDynamoDB:
		return &DynamoDBFactory{}, nil
	case TypeCassandra:
		return &CassandraFactory{}, nil
	case TypeAerospike:
		return &AerosplikeFactory{}, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// AddConnection adds a new database connection
func (m *Manager) AddConnection(name string, dbType DatabaseType, config *database.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.connections[name]; exists {
		return fmt.Errorf("connection %s already exists", name)
	}

	factory, err := GetFactory(dbType)
	if err != nil {
		return err
	}

	db, err := factory.CreateDatabase(config)
	if err != nil {
		return err
	}

	m.connections[name] = db
	return nil
}

// GetConnection returns an existing database connection
func (m *Manager) GetConnection(name string) (database.Database, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("connection %s not found", name)
	}

	return db, nil
}

// RemoveConnection removes a database connection
func (m *Manager) RemoveConnection(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, exists := m.connections[name]
	if !exists {
		return fmt.Errorf("connection %s not found", name)
	}

	if err := db.Disconnect(context.Background()); err != nil {
		return err
	}

	delete(m.connections, name)
	return nil
}

// CloseAll closes all database connections
func (m *Manager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, db := range m.connections {
		if err := db.Disconnect(context.Background()); err != nil {
			errs = append(errs, fmt.Errorf("error closing %s: %w", name, err))
		}
	}

	m.connections = make(map[string]database.Database)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}

// Add factory implementations for each database type
type MySQLFactory struct{}

func (f *MySQLFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return mysql.NewMySQLDB(config), nil
}

type OracleFactory struct{}

func (f *OracleFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return oracle.NewOracleDB(config), nil
}

type DynamoDBFactory struct{}

func (f *DynamoDBFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return dynamodb.NewDynamoDB(config), nil
}

type CassandraFactory struct{}

func (f *CassandraFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return cassandra.NewCassandraDB(config), nil
}

type AerosplikeFactory struct{}

func (f *AerosplikeFactory) CreateDatabase(config *database.Config) (database.Database, error) {
	return aerospike.NewAerospike(config), nil
}
