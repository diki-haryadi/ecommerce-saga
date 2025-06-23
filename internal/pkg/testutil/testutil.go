package testutil

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB represents a test database connection
type TestDB struct {
	*gorm.DB
	Name string
}

// NewTestPostgres creates a new test PostgreSQL database
func NewTestPostgres(t *testing.T) *TestDB {
	t.Helper()

	// Get test database configuration
	dbName := "test_" + time.Now().Format("20060102150405")
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Create test database
	err = db.Exec("CREATE DATABASE " + dbName).Error
	require.NoError(t, err)

	// Connect to test database
	testDSN := dsn + " dbname=" + dbName
	testDB, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	require.NoError(t, err)

	return &TestDB{
		DB:   testDB,
		Name: dbName,
	}
}

// Cleanup removes the test database
func (db *TestDB) Cleanup(t *testing.T) {
	t.Helper()

	sqlDB, err := db.DB.DB()
	require.NoError(t, err)

	// Close connection
	err = sqlDB.Close()
	require.NoError(t, err)

	// Drop test database
	mainDB, err := gorm.Open(postgres.Open(os.Getenv("TEST_DATABASE_URL")), &gorm.Config{})
	require.NoError(t, err)

	err = mainDB.Exec("DROP DATABASE IF EXISTS " + db.Name).Error
	require.NoError(t, err)
}

// TestMongoDB represents a test MongoDB connection
type TestMongoDB struct {
	*mongo.Client
	*mongo.Database
	Name string
}

// NewTestMongoDB creates a new test MongoDB database
func NewTestMongoDB(t *testing.T) *TestMongoDB {
	t.Helper()

	// Get test database configuration
	dbName := "test_" + time.Now().Format("20060102150405")
	uri := os.Getenv("TEST_MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	// Ping database
	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	return &TestMongoDB{
		Client:   client,
		Database: client.Database(dbName),
		Name:     dbName,
	}
}

// Cleanup drops the test database and closes the connection
func (db *TestMongoDB) Cleanup(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	// Drop database
	err := db.Database.Drop(ctx)
	require.NoError(t, err)

	// Close connection
	err = db.Client.Disconnect(ctx)
	require.NoError(t, err)
}

// LoadFixture loads a test fixture file
func LoadFixture(t *testing.T, name string) []byte {
	t.Helper()

	path := filepath.Join("testdata", "fixtures", name)
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return data
}

// MockTime is a helper for testing time-based functionality
type MockTime struct {
	current time.Time
}

// NewMockTime creates a new mock time
func NewMockTime(t time.Time) *MockTime {
	return &MockTime{current: t}
}

// Now returns the current mock time
func (mt *MockTime) Now() time.Time {
	return mt.current
}

// Add adds a duration to the current mock time
func (mt *MockTime) Add(d time.Duration) {
	mt.current = mt.current.Add(d)
}

// Set sets the current mock time
func (mt *MockTime) Set(t time.Time) {
	mt.current = t
}

// Reset resets the mock time to the initial value
func (mt *MockTime) Reset(t time.Time) {
	mt.current = t
}
