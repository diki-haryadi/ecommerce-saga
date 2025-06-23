package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ConfigLoader handles configuration loading from files and environment variables
type ConfigLoader struct {
	viper *viper.Viper
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		viper: viper.New(),
	}
}

// LoadConfig loads configuration from files and environment variables
func (l *ConfigLoader) LoadConfig(configPath string) (*Config, error) {
	// Set default configuration values
	l.setDefaults()

	// Load configuration from file
	if err := l.loadConfigFile(configPath); err != nil {
		return nil, err
	}

	// Load environment variables
	l.loadEnvVars()

	// Unmarshal configuration
	var config Config
	if err := l.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func (l *ConfigLoader) setDefaults() {
	l.viper.SetDefault("app.environment", "development")
	l.viper.SetDefault("app.host", "localhost")
	l.viper.SetDefault("app.port", 8080)

	l.viper.SetDefault("database.sslmode", "disable")
	l.viper.SetDefault("database.port", 5432)

	l.viper.SetDefault("redis.port", 6379)
	l.viper.SetDefault("redis.db", 0)

	l.viper.SetDefault("jwt.access_token_ttl", 3600)
	l.viper.SetDefault("jwt.refresh_token_ttl", 86400)
	l.viper.SetDefault("jwt.signing_algorithm", "HS256")

	l.viper.SetDefault("logger.level", "info")
	l.viper.SetDefault("logger.format", "json")
	l.viper.SetDefault("logger.output_path", "stdout")

	l.viper.SetDefault("monitoring.enabled", true)
	l.viper.SetDefault("monitoring.metrics_path", "/metrics")
	l.viper.SetDefault("monitoring.port", 9090)
}

// loadConfigFile loads configuration from file
func (l *ConfigLoader) loadConfigFile(configPath string) error {
	l.viper.SetConfigName("config")
	l.viper.SetConfigType("yaml")
	l.viper.AddConfigPath(configPath)

	// Get environment from env var or use default
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Try to load environment-specific config
	envConfigName := fmt.Sprintf("config.%s", env)
	l.viper.SetConfigName(envConfigName)

	configFile := filepath.Join(configPath, envConfigName+".yaml")
	if _, err := os.Stat(configFile); err == nil {
		if err := l.viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Fall back to default config
		l.viper.SetConfigName("config")
		if err := l.viper.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to read default config file: %w", err)
		}
	}

	return nil
}

// loadEnvVars loads configuration from environment variables
func (l *ConfigLoader) loadEnvVars() {
	// Set environment variables prefix
	l.viper.SetEnvPrefix("APP")

	// Replace dots with underscores in env vars
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable environment variables
	l.viper.AutomaticEnv()

	// Bind specific environment variables
	envVars := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"MONGODB_URI",
		"MONGODB_DATABASE",
		"MONGODB_USER",
		"MONGODB_PASSWORD",
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"JWT_SECRET",
		"KAFKA_HOST",
		"KAFKA_PORT",
		"KAFKA_USER",
		"KAFKA_PASSWORD",
		"NSQ_HOST",
		"NSQ_PORT",
		"NSQ_USER",
		"NSQ_PASSWORD",
		"NATS_HOST",
		"NATS_PORT",
		"NATS_USER",
		"NATS_PASSWORD",
	}

	for _, env := range envVars {
		if err := l.viper.BindEnv(strings.ToLower(env)); err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to bind env var %s: %v\n", env, err)
		}
	}
}
