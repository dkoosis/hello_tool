// ErrorEnhanced: 2025-05-17
// Package main executes the hello-tool-base service.
// file: cmd/hello-tool-base/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// Added for server timeouts
	"github.com/cockroachdb/errors"
	"github.com/dkoosis/hello-tool-base/internal/apperrors"
	"github.com/dkoosis/hello-tool-base/internal/buildinfo"
	"github.com/dkoosis/hello-tool-base/internal/config"
	"github.com/dkoosis/hello-tool-base/internal/logging"
	"github.com/dkoosis/hello-tool-base/internal/middleware"
)

var (
	// appLog is the global logger, primarily for startup and shutdown messages.
	// Request-specific logging will use the logger from context.
	appLog logging.Logger
	// cfg holds the application configuration, moved to package level for access in main.
	cfg *config.Config
)

// GreetingResponse defines the structure for a successful greeting.
type GreetingResponse struct {
	Message string `json:"message"`
}

// ClientErrorResponse defines the structure for a JSON error message to the client.
type ClientErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// respondWithJSON is a helper function to respond with JSON.
// It now takes a logger for consistent error logging.
func respondWithJSON(l logging.Logger, w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Use the passed-in logger which should have trace_id
		wrappedErr := errors.Wrap(err, "respondWithJSON: failed to encode JSON response payload")
		l.Error("Failed to encode JSON response", "error", fmt.Sprintf("%+v", wrappedErr))
	}
}

// respondWithError is a helper function to respond with a JSON error to the client.
// It now takes a logger.
func respondWithError(l logging.Logger, w http.ResponseWriter, statusCode int, clientMessage string, clientDetails string, internalErr error) {
	// Use the passed-in logger which should have trace_id
	l.Error("Server-side error occurred",
		"internal_error", fmt.Sprintf("%+v", internalErr), // Ensure internalErr is wrapped with stack if not already
		"client_message", clientMessage,
		"client_details", clientDetails,
		"http_status_code", statusCode,
	)
	respondWithJSON(l, w, statusCode, ClientErrorResponse{Error: clientMessage, Details: clientDetails})
}

// helloHandler handles requests to the /hello endpoint.
// It expects a 'name' query parameter and returns a personalized greeting.
// If the 'name' parameter is missing, it returns a 400 Bad Request.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the request-scoped logger from context
	reqLogger := middleware.GetLoggerFromContext(r.Context()).WithField("handler", "helloHandler")
	reqLogger.Info("Received request", "remote_addr", r.RemoteAddr, "path", r.URL.Path)

	name := r.URL.Query().Get("name")
	if name == "" {
		internalErr := apperrors.NewInvalidParamsError(
			"helloHandler: missing 'name' query parameter",
			nil, // No underlying cause for this specific validation error
			map[string]interface{}{
				"parameter_name": "name",
				"query_path":     r.URL.String(),
				"method":         r.Method,
			},
		)
		respondWithError(reqLogger, w, http.StatusBadRequest,
			"Invalid Request Parameter",
			"The 'name' query parameter is required.",
			internalErr)
		return
	}

	message := fmt.Sprintf("Hello, %s, from your Go Cloud Run service!", name)
	response := GreetingResponse{Message: message}
	respondWithJSON(reqLogger, w, http.StatusOK, response)
	reqLogger.Info("Successfully processed /hello request", "name", name)
}

// rootHandler handles requests to the / (root) endpoint.
// It returns general service information, including build details and a welcome message.
// For any other path, it returns a 404 Not Found.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	reqLogger := middleware.GetLoggerFromContext(r.Context()).WithField("handler", "rootHandler")

	if r.URL.Path != "/" {
		// Log before calling http.NotFound, as it writes the response
		reqLogger.Warn("Path not found", "remote_addr", r.RemoteAddr, "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	reqLogger.Info("Received request for root path", "remote_addr", r.RemoteAddr)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	writeAndLog := func(format string, args ...interface{}) {
		if _, err := fmt.Fprintf(w, format, args...); err != nil {
			// Using cockroachdb/errors.Wrapf to add context and preserve stack trace if any from fmt.Fprintf
			wrappedErr := errors.Wrapf(err, "rootHandler: failed to write response for path %s", r.URL.Path)
			// Use reqLogger here as it has trace_id and handler context
			reqLogger.Error("Failed to write to response stream", "error", fmt.Sprintf("%+v", wrappedErr))
		}
	}

	writeAndLog("Hello World Root! This is the base Go service. Try /hello?name=YourName\n")
	writeAndLog("---\n")
	serviceName := cfg.Server.Name // Use configured server name
	if serviceName == "" {
		serviceName = "hello-tool-base" // Fallback
	}
	writeAndLog("Service: %s\n", serviceName)
	writeAndLog("Version: %s\n", buildinfo.Version)
	writeAndLog("Commit:  %s\n", buildinfo.CommitHash)
	writeAndLog("Built:   %s\n", buildinfo.BuildDate)

	// Fetch and write the TraceID
	traceID := middleware.GetTraceIDFromContext(r.Context())
	if traceID != "" {
		writeAndLog("TraceID: %s\n", traceID)
	}

	reqLogger.Info("Successfully processed / request")
}

// healthHandler provides a health check endpoint for the service.
// It returns the service's current operational status, version, commit, and build date.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	reqLogger := middleware.GetLoggerFromContext(r.Context()).WithField("handler", "healthHandler")
	// For frequent health checks, detailed logging per request might be too verbose.
	// reqLogger.Debug("Received health check request", "path", r.URL.Path) // Use Debug if preferred

	healthStatus := struct {
		Status    string `json:"status"`
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildDate string `json:"buildDate"`
		TraceID   string `json:"traceId,omitempty"` // Optionally include traceID in health response
	}{
		Status:    "OK",
		Version:   buildinfo.Version,
		Commit:    buildinfo.CommitHash,
		BuildDate: buildinfo.BuildDate,
		TraceID:   middleware.GetTraceIDFromContext(r.Context()), // Get trace ID from context
	}
	respondWithJSON(reqLogger, w, http.StatusOK, healthStatus) // Pass logger
	// reqLogger.Debug("Successfully processed /health request")
}

// main is the entry point for the application.
// It initializes configuration, logging, sets up HTTP routes, and starts the server.
// It also handles graceful shutdown on SIGINT or SIGTERM signals.
func main() {
	// Setup default logger for application-level logs (startup, shutdown)
	logging.SetupDefaultLogger("debug") // Or your desired default level
	appLog = logging.GetLogger("hello-tool")

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}

	var err error
	cfg, err = config.LoadFromFile(cfgPath)
	if err != nil {
		// err from LoadFromFile should already be well-wrapped by cockroachdb/errors.
		// Logging with %+v will include the stack trace.
		appLog.Error("Failed to load configuration. Shutting down.", "path", cfgPath, "error", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	// Use appLog for startup messages
	appLog.Info("Service starting...",
		"name", cfg.Server.Name,
		"version", buildinfo.Version,
		"commit", buildinfo.CommitHash,
		"buildDate", buildinfo.BuildDate,
		"port", cfg.Server.Port,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", rootHandler)

	// Apply middleware
	// The Tracing middleware is instantiated with the appLog, which will serve as the
	// base for request-scoped loggers created by the middleware.
	tracingMiddleware := middleware.Tracing(appLog)
	handlerWithTracing := tracingMiddleware(mux)
	// Add other middleware here if needed, e.g.:
	// handlerWithAuth := authMiddleware(handlerWithTracing)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      handlerWithTracing, // Use the middleware-wrapped handler
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		appLog.Info("Server listening", "address", srv.Addr)
		// errors.Is correctly checks for http.ErrServerClosed
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Error from ListenAndServe might not be wrapped by cockroachdb yet, so wrap it for consistent stack trace.
			// However, standard library errors like this are often clear enough.
			// For consistency with prompt, let's ensure it's potentially wrapped if we want to add our own context/type.
			// But since we just log it, fmt.Sprintf("%+v", err) is often sufficient if it's already rich.
			// Let's keep it as is, as %+v on stdlib errors can also be informative.
			// If we were to pass it to apperrors, then we'd wrap:
			// wrappedErr := errors.Wrap(err, "main: server ListenAndServe failed")
			appLog.Error("Server failed to start or encountered an error. Shutting down.", "error", fmt.Sprintf("%+v", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	receivedSignal := <-quit
	appLog.Info("Received signal, initiating shutdown...", "signal", receivedSignal.String())

	// Create a context with a timeout for graceful shutdown
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancelShutdown()

	// Error from srv.Shutdown might also be a standard library error.
	if err := srv.Shutdown(shutdownCtx); err != nil {
		// wrappedErr := errors.Wrap(err, "main: server Shutdown failed")
		appLog.Error("Server shutdown failed", "error", fmt.Sprintf("%+v", err))
	} else {
		appLog.Info("Server exited gracefully.")
	}
}
