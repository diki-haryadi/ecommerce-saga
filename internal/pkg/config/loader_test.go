package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigLoader(t *testing.T) {
	// Create temporary directory for test config files
	tempDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test config files
	devConfig := `
app:
  name: "test-app"
  environment: "development"
  host: "localhost"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "test_user"
  password: "test_pass"
  dbname: "test_db"

brokers:
  default: "kafka:default"
  kafka:
    default:
      enabled: true
      host: "localhost"
      port: "9092"
`

	prodConfig := `
app:
  name: "test-app"
  environment: "production"
  host: "0.0.0.0"
  port: 80

database:
  host: "${DB_HOST}"
  port: ${DB_PORT}
  user: "${DB_USER}"
  password: "${DB_PASSWORD}"
  dbname: "${DB_NAME}"

brokers:
  default: "kafka:default"
  kafka:
    default:
      enabled: true
      host: "${KAFKA_HOST}"
      port: "${KAFKA_PORT}"
`

	err = os.WriteFile(filepath.Join(tempDir, "config.development.yaml"), []byte(devConfig), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "config.production.yaml"), []byte(prodConfig), 0644)
	require.NoError(t, err)

	tests := []struct {
		name        string
		env         string
		envVars     map[string]string
		checkConfig func(*testing.T, *Config)
	}{
		{
			name: "development config",
			env:  "development",
			checkConfig: func(t *testing.T, cfg *Config) {
				require.Equal(t, "test-app", cfg.App.Name)
				require.Equal(t, "development", cfg.App.Environment)
				require.Equal(t, "localhost", cfg.App.Host)
				require.Equal(t, 8080, cfg.App.Port)

				require.Equal(t, "localhost", cfg.Database.Host)
				require.Equal(t, 5432, cfg.Database.Port)
				require.Equal(t, "test_user", cfg.Database.User)
				require.Equal(t, "test_pass", cfg.Database.Password)
				require.Equal(t, "test_db", cfg.Database.DBName)

				require.Equal(t, "kafka:default", cfg.Brokers.Default)
				require.True(t, cfg.Brokers.Kafka["default"].Enabled)
				require.Equal(t, "localhost", cfg.Brokers.Kafka["default"].Host)
				require.Equal(t, "9092", cfg.Brokers.Kafka["default"].Port)
			},
		},
		{
			name: "production config with env vars",
			env:  "production",
			envVars: map[string]string{
				"APP_DB_HOST":     "db.prod",
				"APP_DB_PORT":     "5432",
				"APP_DB_USER":     "prod_user",
				"APP_DB_PASSWORD": "prod_pass",
				"APP_DB_NAME":     "prod_db",
				"APP_KAFKA_HOST":  "kafka.prod",
				"APP_KAFKA_PORT":  "9093",
			},
			checkConfig: func(t *testing.T, cfg *Config) {
				require.Equal(t, "test-app", cfg.App.Name)
				require.Equal(t, "production", cfg.App.Environment)
				require.Equal(t, "0.0.0.0", cfg.App.Host)
				require.Equal(t, 80, cfg.App.Port)

				require.Equal(t, "db.prod", cfg.Database.Host)
				require.Equal(t, 5432, cfg.Database.Port)
				require.Equal(t, "prod_user", cfg.Database.User)
				require.Equal(t, "prod_pass", cfg.Database.Password)
				require.Equal(t, "prod_db", cfg.Database.DBName)

				require.Equal(t, "kafka:default", cfg.Brokers.Default)
				require.True(t, cfg.Brokers.Kafka["default"].Enabled)
				require.Equal(t, "kafka.prod", cfg.Brokers.Kafka["default"].Host)
				require.Equal(t, "9093", cfg.Brokers.Kafka["default"].Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("APP_ENV", tt.env)
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				os.Unsetenv("APP_ENV")
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			// Load configuration
			loader := NewConfigLoader()
			cfg, err := loader.LoadConfig(tempDir)
			require.NoError(t, err)
			require.NotNil(t, cfg)

			// Check configuration
			tt.checkConfig(t, cfg)
		})
	}
}

func TestConfigLoader_Defaults(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create minimal config file
	minConfig := `
app:
  name: "test-app"
`
	err = os.WriteFile(filepath.Join(tempDir, "config.yaml"), []byte(minConfig), 0644)
	require.NoError(t, err)

	// Load configuration
	loader := NewConfigLoader()
	cfg, err := loader.LoadConfig(tempDir)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Check default values
	require.Equal(t, "development", cfg.App.Environment)
	require.Equal(t, "localhost", cfg.App.Host)
	require.Equal(t, 8080, cfg.App.Port)

	require.Equal(t, "disable", cfg.Database.SSLMode)
	require.Equal(t, 5432, cfg.Database.Port)

	require.Equal(t, 6379, cfg.Redis.Port)
	require.Equal(t, 0, cfg.Redis.DB)

	require.Equal(t, 3600, cfg.JWT.AccessTokenTTL)
	require.Equal(t, 86400, cfg.JWT.RefreshTokenTTL)
	require.Equal(t, "HS256", cfg.JWT.SigningAlgorithm)

	require.Equal(t, "info", cfg.Logger.Level)
	require.Equal(t, "json", cfg.Logger.Format)
	require.Equal(t, "stdout", cfg.Logger.OutputPath)

	require.True(t, cfg.Monitoring.Enabled)
	require.Equal(t, "/metrics", cfg.Monitoring.MetricsPath)
	require.Equal(t, 9090, cfg.Monitoring.Port)
}
