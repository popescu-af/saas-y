package log

import (
	"log"

	"go.uber.org/zap"
)

// Context is the type for log contexts.
type Context map[string]interface{}

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

func logCtx(msg string, ctx Context, logFn func(string, ...zap.Field)) {
	fields := make([]zap.Field, 0, len(ctx))

	for k, v := range ctx {
		fields := append(fields, zap.Any(k, v))
	}

	logFn(msg, fields...)
}

// Debug logs a message with DEBUG level.
func Debug(msg string) {
	DebugCtx(msg, nil)
}

// DebugCtx logs a message with DEBUG level.
func DebugCtx(msg string, ctx Context) {
	logCtx(msg, ctx, logger.Debug)
}

// Error logs a message with ERROR level.
func Error(msg string) {
	ErrorCtx(msg, nil)
}

// ErrorCtx logs a message with ERROR level.
func ErrorCtx(msg string, ctx Context) {
	logCtx(msg, ctx, logger.Error)
}

// Fatal logs a message with FATAL level.
func Fatal(msg string) {
	FatalCtx(msg, nil)
}

// FatalCtx logs a message with FATAL level.
func FatalCtx(msg string, ctx Context) {
	logCtx(msg, ctx, logger.Fatal)
}

// Info logs a message with INFO level.
func Info(msg string) {
	InfoCtx(msg, nil)
}

// InfoCtx logs a message with INFO level.
func InfoCtx(msg string, ctx Context) {
	logCtx(msg, ctx, logger.Info)
}

// Warn logs a message with WARN level.
func Warn(msg string) {
	WarnCtx(msg, nil)
}

// WarnCtx logs a message with WARN level.
func WarnCtx(msg string, ctx Context) {
	logCtx(msg, ctx, logger.Warn)
}
