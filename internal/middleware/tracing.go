// Package middleware provides HTTP server middleware functionalities.
// Middleware are functions that wrap HTTP handlers to add cross-cutting concerns
// such as request tracing, logging context enrichment, authentication, and more,
// before or after the main handler logic is executed.
// file: internal/middleware/tracing.go
package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dkoosis/hello-tool-base/internal/logging" // Adjust to your actual module path
	"github.com/google/uuid"
)

const (
	// HeaderCloudTraceContext is the HTTP header field used by Google Cloud
	// for propagating trace context information.
	HeaderCloudTraceContext = "X-Cloud-Trace-Context"
	// HeaderTraceID is a common HTTP header field used to return the trace ID
	// back to the client or for other tracing systems.
	HeaderTraceID = "X-Trace-ID"
)

// Tracing is a middleware that adds a trace ID to each incoming HTTP request.
// It retrieves an existing trace ID from the X-Cloud-Trace-Context header if present,
// otherwise, it generates a new UUID-based trace ID. This trace ID is then set
// in the X-Trace-ID response header.
//
// Furthermore, it creates a request-scoped logger enriched with this trace ID
// and adds both the trace ID and the logger to the request's context.
// The baseLogger provided is used as the foundation for these request-scoped loggers.
func Tracing(baseLogger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get(HeaderCloudTraceContext)
			if traceID == "" {
				traceID = fmt.Sprintf("generated-%s", uuid.New().String())
			}
			w.Header().Set(HeaderTraceID, traceID)

			ctx := context.WithValue(r.Context(), TraceIDContextKey, traceID)

			// Create request-scoped logger WITH ONLY trace_id added by middleware.
			// The baseLogger already has its component (e.g., "app" or "hello-tool-base-main").
			// Handlers will add their specific context (like "handler":"helloHandler").
			requestLogger := baseLogger.WithField("trace_id", traceID)
			ctx = context.WithValue(ctx, LoggerContextKey, requestLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetLoggerFromContext retrieves the request-scoped logger from the provided context.Context.
// It expects the logger to have been set using the LoggerContextKey by the Tracing middleware (or similar).
// If no logger is found in the context (which ideally should not happen if middleware is correctly applied),
// it returns a fallback logger and logs a warning about the missing logger.
func GetLoggerFromContext(ctx context.Context) logging.Logger {
	logger, ok := ctx.Value(LoggerContextKey).(logging.Logger)
	if !ok || logger == nil {
		// This should ideally not happen if middleware is correctly applied.
		// Fallback to a default logger or a new one with a warning.
		fallbackLogger := logging.GetLogger("context_logger_fallback")
		fallbackLogger.Warn("Logger not found in context, using fallback.")
		return fallbackLogger
	}
	return logger
}

// GetTraceIDFromContext retrieves the trace ID from the provided context.Context.
// It expects the trace ID to have been set using the TraceIDContextKey by the Tracing middleware.
// If no trace ID is found in the context, it returns an empty string.
func GetTraceIDFromContext(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDContextKey).(string)
	if !ok {
		return ""
	}
	return traceID
}
