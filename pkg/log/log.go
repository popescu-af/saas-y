package log

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeFormat = "2020-01-01 20:20:20.999"

// Context is the type for log contexts.
type Context map[string]interface{}

type loggerWrapper struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

var defaultLogger = newLogger()

func newLogger() *loggerWrapper {
	logLevel := zap.NewAtomicLevel()
	enc := newEncoder()
	core := zapcore.NewCore(enc, os.Stdout, logLevel)
	return &loggerWrapper{
		logger: zap.New(core).Named("root"),
		level:  logLevel,
	}
}

func newEncoder() zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.CallerKey = ""     // disable caller
	encCfg.StacktraceKey = "" // disable stack trace
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encCfg.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(timeFormat))
	})
	return zapcore.NewConsoleEncoder(encCfg)
}

// SetLevel sets the logging level.
func SetLevel(logLevel string) {
	switch logLevel {
	case "DEBUG":
		defaultLogger.level.SetLevel(zapcore.DebugLevel)
	case "INFO":
		defaultLogger.level.SetLevel(zapcore.InfoLevel)
	case "WARN":
		defaultLogger.level.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		defaultLogger.level.SetLevel(zapcore.ErrorLevel)
	case "CRITICAL":
		defaultLogger.level.SetLevel(zapcore.FatalLevel)
	}
}

// Sync flushes the log cache.
func Sync() {
	if err := defaultLogger.logger.Sync(); err != nil {
		log.Printf("failed to flush log cache - %v", err)
	}
}

func logCtx(msg string, ctx Context, logFn func(string, ...zap.Field)) {
	fields := make([]zap.Field, 0, len(ctx))

	for k, v := range ctx {
		fields = append(fields, zap.Any(k, v))
	}

	logFn(msg, fields...)
}

// Debug logs a message with DEBUG level.
func Debug(msg string) {
	DebugCtx(msg, nil)
}

// DebugCtx logs a message with DEBUG level.
func DebugCtx(msg string, ctx Context) {
	logCtx(msg, ctx, defaultLogger.logger.Debug)
}

// Error logs a message with ERROR level.
func Error(msg string) {
	ErrorCtx(msg, nil)
}

// ErrorCtx logs a message with ERROR level.
func ErrorCtx(msg string, ctx Context) {
	logCtx(msg, ctx, defaultLogger.logger.Error)
}

// Fatal logs a message with FATAL level.
func Fatal(msg string) {
	FatalCtx(msg, nil)
}

// FatalCtx logs a message with FATAL level.
func FatalCtx(msg string, ctx Context) {
	logCtx(msg, ctx, defaultLogger.logger.Fatal)
}

// Info logs a message with INFO level.
func Info(msg string) {
	InfoCtx(msg, nil)
}

// InfoCtx logs a message with INFO level.
func InfoCtx(msg string, ctx Context) {
	logCtx(msg, ctx, defaultLogger.logger.Info)
}

// Warn logs a message with WARN level.
func Warn(msg string) {
	WarnCtx(msg, nil)
}

// WarnCtx logs a message with WARN level.
func WarnCtx(msg string, ctx Context) {
	logCtx(msg, ctx, defaultLogger.logger.Warn)
}
