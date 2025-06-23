package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

type Field = zapcore.Field

type zapLogger struct {
	logger *zap.Logger
}

func NewLogger(level string) (Logger, error) {
	config := zap.NewProductionConfig()

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Customize output format
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "" // Disable stacktrace

	logger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: logger,
	}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
	os.Exit(1)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	// Extract trace ID or other context values here
	return l
}

// Field constructors
func String(key, value string) Field {
	return zap.String(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

func Float64(key string, value float64) Field {
	return zap.Float64(key, value)
}

func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

func Error(err error) Field {
	return zap.Error(err)
}

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}
