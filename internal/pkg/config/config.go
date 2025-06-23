package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	MongoDB    MongoDBConfig    `mapstructure:"mongodb"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Brokers    BrokersConfig    `mapstructure:"brokers"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	AccessTokenTTL   int    `mapstructure:"access_token_ttl"`
	RefreshTokenTTL  int    `mapstructure:"refresh_token_ttl"`
	SigningAlgorithm string `mapstructure:"signing_algorithm"`
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	MetricsPath string `mapstructure:"metrics_path"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
}

// BrokersConfig holds configuration for all message brokers
type BrokersConfig struct {
	Default string                          `mapstructure:"default"`
	Kafka   map[string]BrokerInstanceConfig `mapstructure:"kafka"`
	NSQ     map[string]BrokerInstanceConfig `mapstructure:"nsq"`
	NATS    map[string]BrokerInstanceConfig `mapstructure:"nats"`
}

// BrokerInstanceConfig holds configuration for a single broker instance
type BrokerInstanceConfig struct {
	Enabled  bool                   `mapstructure:"enabled"`
	Host     string                 `mapstructure:"host"`
	Port     string                 `mapstructure:"port"`
	Username string                 `mapstructure:"username"`
	Password string                 `mapstructure:"password"`
	Options  map[string]interface{} `mapstructure:"options"`
}

// LoadConfig loads configuration from files and environment variables
func LoadConfig(configPath string) (*Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	// Load environment-specific config
	env := getEnvironment()
	if env != "" {
		viper.SetConfigName(fmt.Sprintf("config.%s", env))
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Set environment variables prefix
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// getEnvironment returns the current environment
func getEnvironment() string {
	env := viper.GetString("environment")
	if env == "" {
		env = "development"
	}
	return env
}
