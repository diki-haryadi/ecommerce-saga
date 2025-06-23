package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	NSQ      NSQConfig      `mapstructure:"nsq"`
	NATS     NATSConfig     `mapstructure:"nats"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AuthConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	TokenExpiry   time.Duration `mapstructure:"token_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	GroupID string   `mapstructure:"group_id"`
}

type NSQConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type NATSConfig struct {
	URL     string `mapstructure:"url"`
	Cluster string `mapstructure:"cluster"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetPostgresDSN returns PostgreSQL connection string
func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Postgres.Host,
		c.Database.Postgres.Port,
		c.Database.Postgres.User,
		c.Database.Postgres.Password,
		c.Database.Postgres.DBName,
		c.Database.Postgres.SSLMode,
	)
}

// GetRabbitMQURL returns RabbitMQ connection URL
func (c *Config) GetRabbitMQURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.VHost,
	)
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
