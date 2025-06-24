package database

import (
	"time"

	"github.com/diki-haryadi/ecommerce-saga/internal/pkg/database"
)

// ConnectionBuilder helps build database configurations
type ConnectionBuilder struct {
	config *database.Config
}

// NewConnectionBuilder creates a new connection builder
func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{
		config: &database.Config{
			MaxConnections: 10,
			ConnectTimeout: 10 * time.Second,
			MaxIdleTime:    5 * time.Minute,
			MaxRetries:     3,
			RetryInterval:  time.Second,
			SSLMode:        "disable",
			AdditionalOpts: make(map[string]string),
		},
	}
}

// WithHost sets the host
func (b *ConnectionBuilder) WithHost(host string) *ConnectionBuilder {
	b.config.Host = host
	return b
}

// WithPort sets the port
func (b *ConnectionBuilder) WithPort(port int) *ConnectionBuilder {
	b.config.Port = port
	return b
}

// WithCredentials sets the credentials
func (b *ConnectionBuilder) WithCredentials(username, password string) *ConnectionBuilder {
	b.config.Username = username
	b.config.Password = password
	return b
}

// WithDatabase sets the database name
func (b *ConnectionBuilder) WithDatabase(database string) *ConnectionBuilder {
	b.config.Database = database
	return b
}

// WithMaxConnections sets the maximum number of connections
func (b *ConnectionBuilder) WithMaxConnections(max int) *ConnectionBuilder {
	b.config.MaxConnections = max
	return b
}

// WithConnectTimeout sets the connection timeout
func (b *ConnectionBuilder) WithConnectTimeout(timeout time.Duration) *ConnectionBuilder {
	b.config.ConnectTimeout = timeout
	return b
}

// WithMaxIdleTime sets the maximum idle time
func (b *ConnectionBuilder) WithMaxIdleTime(duration time.Duration) *ConnectionBuilder {
	b.config.MaxIdleTime = duration
	return b
}

// WithRetryPolicy sets the retry policy
func (b *ConnectionBuilder) WithRetryPolicy(maxRetries int, interval time.Duration) *ConnectionBuilder {
	b.config.MaxRetries = maxRetries
	b.config.RetryInterval = interval
	return b
}

// WithSSLMode sets the SSL mode
func (b *ConnectionBuilder) WithSSLMode(mode string) *ConnectionBuilder {
	b.config.SSLMode = mode
	return b
}

// WithOption adds an additional option
func (b *ConnectionBuilder) WithOption(key, value string) *ConnectionBuilder {
	b.config.AdditionalOpts[key] = value
	return b
}

// Build returns the final configuration
func (b *ConnectionBuilder) Build() *database.Config {
	return b.config
}

// Example usage:
// config := NewConnectionBuilder().
//     WithHost("localhost").
//     WithPort(5432).
//     WithCredentials("user", "pass").
//     WithDatabase("mydb").
//     WithMaxConnections(20).
//     WithSSLMode("require").
//     Build()
