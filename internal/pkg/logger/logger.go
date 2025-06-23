package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string // json or console
	OutputPath string
}

// Initialize initializes the logger with the given configuration
func Initialize(config *Config) error {
	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create encoder based on format
	var encoder zapcore.Encoder
	switch config.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create core
	file, err := os.OpenFile(config.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(file),
		level,
	)

	// Create logger
	log = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return nil
}

// Debug logs a debug message
func Debug(msg string, fields ...zapcore.Field) {
	log.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func With(fields ...zapcore.Field) *zap.Logger {
	return log.With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	return log.Sync()
}

// Fields creates a list of Zap fields from a map
func Fields(fields map[string]interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// NewContext creates a new context with the given fields
func NewContext(fields map[string]interface{}) *zap.Logger {
	return With(Fields(fields)...)
}
