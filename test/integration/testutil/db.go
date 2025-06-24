package testutil

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB represents a test database connection
type TestDB struct {
	DB *gorm.DB
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnvOrDefault("TEST_DB_HOST", "localhost"),
		getEnvOrDefault("TEST_DB_PORT", "5432"),
		getEnvOrDefault("TEST_DB_USER", "postgres"),
		getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		getEnvOrDefault("TEST_DB_NAME", "ppob_test"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return &TestDB{DB: db}
}

// Cleanup performs cleanup after tests
func (tdb *TestDB) Cleanup() {
	// Get the underlying SQL DB
	sqlDB, err := tdb.DB.DB()
	if err != nil {
		log.Printf("Error getting underlying SQL DB: %v", err)
		return
	}
	sqlDB.Close()
}

// TruncateTables truncates all tables in the test database
func (tdb *TestDB) TruncateTables(tables ...string) error {
	for _, table := range tables {
		if err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GenerateTestID generates a unique test identifier
func GenerateTestID() string {
	return uuid.New().String()
}
