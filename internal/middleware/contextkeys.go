// contextkeys.go defines unexported keys used for storing and retrieving
// values from a context.Context, specifically within the middleware package.
// Using unexported types for context keys avoids collisions with keys defined in other packages.
// file: internal/middleware/contextkeys.go
package middleware

// contextKey is an unexported type for context keys to avoid collisions.
// This prevents other packages from using the same key names to store values in the context,
// ensuring that middleware-specific values are not accidentally overwritten or accessed.
type contextKey string

// String makes it easier to debug context keys if they are logged or inspected.
// It returns a string representation of the contextKey.
func (c contextKey) String() string {
	return "middleware context key " + string(c)
}

var (
	// LoggerContextKey is the context key used to store and retrieve the
	// request-scoped logger instance within a context.Context.
	LoggerContextKey = contextKey("logger")
	// TraceIDContextKey is the context key used to store and retrieve the
	// trace ID associated with a request within a context.Context.
	TraceIDContextKey = contextKey("traceID")
)
