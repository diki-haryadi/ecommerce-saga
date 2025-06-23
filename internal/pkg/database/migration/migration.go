package migration

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/logger"
)

// Config holds migration configuration
type Config struct {
	DatabaseURL string
	SourceURL   string
}

// Manager handles database migrations
type Manager struct {
	migrate *migrate.Migrate
	config  *Config
}

// NewManager creates a new migration manager
func NewManager(config *Config) (*Manager, error) {
	m, err := migrate.New(config.SourceURL, config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return &Manager{
		migrate: m,
		config:  config,
	}, nil
}

// Up applies all available migrations
func (m *Manager) Up() error {
	logger.Info("starting database migration")
	startTime := time.Now()

	if err := m.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no migration needed")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Info("database migration completed",
		zap.Duration("duration", time.Since(startTime)),
	)
	return nil
}

// Down reverts all migrations
func (m *Manager) Down() error {
	logger.Info("reverting all migrations")
	startTime := time.Now()

	if err := m.migrate.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no migration to revert")
			return nil
		}
		return fmt.Errorf("failed to revert migrations: %w", err)
	}

	logger.Info("migrations reverted",
		zap.Duration("duration", time.Since(startTime)),
	)
	return nil
}

// Version returns the current migration version
func (m *Manager) Version() (uint, bool, error) {
	return m.migrate.Version()
}

// Force sets a specific migration version
func (m *Manager) Force(version int) error {
	return m.migrate.Force(version)
}

// Steps migrates up or down by a specific number of versions
func (m *Manager) Steps(n int) error {
	return m.migrate.Steps(n)
}

// Close closes the migration manager
func (m *Manager) Close() error {
	sourceErr, databaseErr := m.migrate.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed to close source: %w", sourceErr)
	}
	if databaseErr != nil {
		return fmt.Errorf("failed to close database: %w", databaseErr)
	}
	return nil
}
