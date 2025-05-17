// File: internal/middleware/contextkeys.go
package middleware

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

// String makes it easier to debug.
func (c contextKey) String() string {
	return "middleware context key " + string(c)
}

var (
	// LoggerContextKey is the context key for the request-scoped logger.
	LoggerContextKey = contextKey("logger")
	// TraceIDContextKey is the context key for the trace ID.
	TraceIDContextKey = contextKey("traceID")
)
