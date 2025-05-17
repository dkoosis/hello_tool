// File: internal/middleware/tracing.go
package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dkoosis/hello-tool-base/internal/logging" // Adjust to your actual module path
	"github.com/google/uuid"
)

const (
	// HeaderCloudTraceContext is the header used by Google Cloud for trace context.
	HeaderCloudTraceContext = "X-Cloud-Trace-Context"
	// HeaderTraceID is a common header to return the trace ID.
	HeaderTraceID = "X-Trace-ID"
)

// Tracing adds a trace ID to each request and a request-scoped logger to the context.
// File: internal/middleware/tracing.go
// ...
// Tracing adds a trace ID to each request and a request-scoped logger to the context.
func Tracing(baseLogger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ... (traceID generation and response header setting is the same) ...
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

// GetLoggerFromContext retrieves the request-scoped logger from the context.
// If no logger is found, it returns a new logger with a warning.
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

// GetTraceIDFromContext retrieves the trace ID from the context.
// Returns an empty string if not found.
func GetTraceIDFromContext(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDContextKey).(string)
	if !ok {
		return ""
	}
	return traceID
}
