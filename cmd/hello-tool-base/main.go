// cmd/hello-tool-base/main.go
package main

import (
	"context" // Added for graceful shutdown
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal" // Added for signal handling
	"syscall"   // Added for signal handling

	// Added for server timeouts
	"github.com/cockroachdb/errors"
	// Ensure these import paths are correct for your project structure
	"github.com/dkoosis/hello-tool-base/internal/apperrors"
	"github.com/dkoosis/hello-tool-base/internal/buildinfo"
	"github.com/dkoosis/hello-tool-base/internal/config"
	"github.com/dkoosis/hello-tool-base/internal/logging"
)

var (
	log logging.Logger
)

// GreetingResponse defines the structure for a successful greeting.
type GreetingResponse struct {
	Message string `json:"message"`
}

// ClientErrorResponse defines the structure for a JSON error message to the client.
type ClientErrorResponse struct {
	Error   string `json:"error"`             // User-facing error message
	Details string `json:"details,omitempty"` // Optional user-facing details
}

// Helper function to respond with JSON.
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		wrappedErr := errors.Wrap(err, "respondWithJSON: failed to encode JSON response payload")
		log.Error("Failed to encode JSON response", "error", fmt.Sprintf("%+v", wrappedErr))
	}
}

// Helper function to respond with a JSON error to the client.
func respondWithError(w http.ResponseWriter, statusCode int, clientMessage string, clientDetails string, internalErr error) {
	log.Error("Server-side error occurred",
		"internal_error", fmt.Sprintf("%+v", internalErr),
		"client_message", clientMessage,
		"client_details", clientDetails,
		"http_status_code", statusCode,
	)
	respondWithJSON(w, statusCode, ClientErrorResponse{Error: clientMessage, Details: clientDetails})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	reqLogger := log.WithField("handler", "helloHandler").WithField("remote_addr", r.RemoteAddr).WithField("path", r.URL.Path)
	reqLogger.Info("Received request")

	name := r.URL.Query().Get("name")
	if name == "" {
		internalErr := apperrors.NewInvalidParamsError(
			"helloHandler: missing 'name' query parameter",
			nil, // No underlying cause
			map[string]interface{}{
				"parameter_name": "name",
				"query_path":     r.URL.String(),
				"method":         r.Method,
			},
		)
		respondWithError(w, http.StatusBadRequest,
			"Invalid Request Parameter",
			"The 'name' query parameter is required.",
			internalErr)
		return
	}

	message := fmt.Sprintf("Hello, %s, from your Go Cloud Run service!", name)
	response := GreetingResponse{Message: message}
	respondWithJSON(w, http.StatusOK, response)
	reqLogger.Info("Successfully processed /hello request", "name", name)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	reqLogger := log.WithField("handler", "rootHandler").WithField("remote_addr", r.RemoteAddr).WithField("path", r.URL.Path)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		reqLogger.Warn("Path not found")
		return
	}
	reqLogger.Info("Received request for root path")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	writeAndLog := func(format string, args ...interface{}) {
		if _, err := fmt.Fprintf(w, format, args...); err != nil {
			wrappedErr := errors.Wrapf(err, "rootHandler: failed to write response for path %s", r.URL.Path)
			reqLogger.Error("Failed to write to response stream", "error", fmt.Sprintf("%+v", wrappedErr))
		}
	}

	writeAndLog("Hello World Root! This is the base Go service. Try /hello?name=YourName\n")
	writeAndLog("---\n")
	writeAndLog("Service: %s\n", os.Getenv("SERVICE_NAME")) // Using SERVICE_NAME from env for consistency
	writeAndLog("Version: %s\n", buildinfo.Version)
	writeAndLog("Commit:  %s\n", buildinfo.CommitHash)
	writeAndLog("Built:   %s\n", buildinfo.BuildDate)

	reqLogger.Info("Successfully processed / request")
}

func main() {
	// Setup logger first. The level might come from config later.
	// For now, hardcoding to debug for development.
	logging.SetupDefaultLogger("debug")
	log = logging.GetLogger("hello-tool-base-main")

	// Load configuration
	// A common pattern is to allow config path via flag or env var.
	// For this template, we'll assume "config.yaml" in the current dir,
	// or rely on defaults + env vars if it's not found.
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.yaml" // Default config path
	}

	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		log.Error("Failed to load configuration. Shutting down.", "path", cfgPath, "error", fmt.Sprintf("%+v", err))
		os.Exit(1) // Exit if config loading fails
	}

	// Re-setup logger if log level is in config (optional, depends on your config structure)
	// Example: if cfg.Logging.Level != "" { logging.SetupDefaultLogger(cfg.Logging.Level) }
	// log = logging.GetLogger("hello-tool-base-main") // Re-get logger in case level changed

	log.Info("Service starting...",
		"name", cfg.Server.Name,
		"version", buildinfo.Version,
		"commit", buildinfo.CommitHash,
		"buildDate", buildinfo.BuildDate,
		"port", cfg.Server.Port,
		"readTimeout", cfg.Server.ReadTimeout,
		"writeTimeout", cfg.Server.WriteTimeout,
		"idleTimeout", cfg.Server.IdleTimeout,
		"gracefulTimeout", cfg.Server.GracefulTimeout,
	)

	// Setup HTTP request multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", rootHandler)
	// Placeholder for health handler (from Suggestion #2)
	// mux.HandleFunc("/health", healthHandler)

	// Configure the HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start the server in a goroutine so that it doesn't block.
	go func() {
		log.Info("Server listening", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed to start or encountered an unexpected error. Shutting down.", "error", fmt.Sprintf("%+v", err))
			// Consider if a more specific apperror is needed here or if os.Exit is appropriate.
			// For catastrophic startup failure, os.Exit(1) is common.
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// Listen for syscall.SIGINT (Ctrl+C) and syscall.SIGTERM (sent by Docker, Kubernetes, Cloud Run)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	receivedSignal := <-quit
	log.Info("Received signal, initiating shutdown...", "signal", receivedSignal.String())

	// Create a context with a timeout for the graceful shutdown.
	// This is the maximum time we allow for existing requests to finish.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel() // Ensure cancel is called to release resources associated with the context

	// Attempt to gracefully shut down the server.
	// Shutdown() will stop accepting new connections and wait for active connections to finish.
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed (requests may have been interrupted).", "error", fmt.Sprintf("%+v", err))
		// os.Exit(1) // Or handle more gracefully, depending on requirements
	} else {
		log.Info("Server exited gracefully.")
	}
}
