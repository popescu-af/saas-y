package log

import (
	"log"

	"go.uber.org/zap"
)

var logger = newLogger()

func newLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger - %v", err)
	}
	return logger
}

// Sync flushes the log cache.
func Sync() {
	if err := logger.Sync(); err != nil {
		log.Printf("failed to flush log cache - %v", err)
	}
}

// Debug logs a message with DEBUG level.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Error logs a message with ERROR level.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal logs a message with FATAL level.
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Info logs a message with INFO level.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a message with WARN level.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
