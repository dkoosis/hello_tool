// file: internal/logging/slog.go
package logging

// slog.go provides an implementation of the Logger interface using Go's structured log/slog package.

import (
	"context"
	"log/slog"
	"os"
)

// SlogLogger wraps slog.Logger to implement our Logger interface.
type SlogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger creates a new structured logger based on slog.
// It takes a slog.Level to determine the minimum log level for messages to be recorded.
func NewSlogLogger(level slog.Level) *SlogLogger {
	// Create a new handler with the specified level
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	// Create a new logger with the handler
	logger := slog.New(handler)

	return &SlogLogger{
		logger: logger,
	}
}

// Debug logs a debug-level message with the underlying slog logger.
// Arguments are handled as key-value pairs.
func (l *SlogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs an info-level message with the underlying slog logger.
// Arguments are handled as key-value pairs.
func (l *SlogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning-level message with the underlying slog logger.
// Arguments are handled as key-value pairs.
func (l *SlogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs an error-level message with the underlying slog logger.
// Arguments are handled as key-value pairs.
func (l *SlogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// WithField returns a new SlogLogger instance with the specified key-value pair
// added to its structured logging context.
func (l *SlogLogger) WithField(key string, value any) Logger {
	// Create a new logger with the field
	newLogger := l.logger.With(key, value)

	return &SlogLogger{
		logger: newLogger,
	}
}

// SetupDefaultLogger initializes the default logger for the application using SlogLogger.
// It parses the string log level and configures a global logger instance.
// Call this early in main() to set up logging for the entire application.
func SetupDefaultLogger(level string) {
	// Convert string level to slog.Level
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		// Default to info level
		logLevel = slog.LevelInfo
	}

	// Create and set the default logger
	logger := NewSlogLogger(logLevel)
	SetDefaultLogger(logger)
}

// WithContext returns a new SlogLogger instance.
// For SlogLogger, this implementation primarily satisfies the Logger interface;
// slog itself does not directly use context for cancellation in its core logging methods
// in the same way some other logging libraries might. Contextual values can be added
// using WithField or by passing them as arguments to log methods.
// Here, "ctx:attached" is added as an example field to indicate context was passed.
func (l *SlogLogger) WithContext(_ context.Context) Logger { // Renamed ctx to _
	// slog.Logger doesn't have a dedicated WithContext method for propagation like some loggers.
	// We return a new logger with a marker field indicating context was considered.
	// True context propagation for slog would involve extracting values from ctx
	// and adding them as structured fields.
	return &SlogLogger{
		logger: l.logger.With("ctx", "attached"),
	}
}
