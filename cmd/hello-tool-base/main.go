// cmd/hello-tool-base/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"
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
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
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
			nil,
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
	serviceName := cfg.Server.Name // Use configured server name
	if serviceName == "" {
		serviceName = "hello-tool-base" // Fallback if not in config for some reason
	}
	writeAndLog("Service: %s\n", serviceName)
	writeAndLog("Version: %s\n", buildinfo.Version)
	writeAndLog("Commit:  %s\n", buildinfo.CommitHash)
	writeAndLog("Built:   %s\n", buildinfo.BuildDate)

	reqLogger.Info("Successfully processed / request")
}

// healthHandler provides a health check endpoint.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// For frequent health checks, detailed logging per request might be too verbose.
	// Consider logging only on errors or if status is not OK.
	// log.Debug("Received health check request", "path", r.URL.Path)

	healthStatus := struct {
		Status    string `json:"status"`
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildDate string `json:"buildDate"`
	}{
		Status:    "OK",
		Version:   buildinfo.Version,
		Commit:    buildinfo.CommitHash,
		BuildDate: buildinfo.BuildDate,
	}
	respondWithJSON(w, http.StatusOK, healthStatus)
}

// Declare cfg at the package level or pass it around. For simplicity here, package level.
// However, for better testability and explicit dependency management, passing cfg would be better.
// For this iteration, let's make it accessible to rootHandler.
var cfg *config.Config

func main() {
	logging.SetupDefaultLogger("debug")
	log = logging.GetLogger("hello-tool-base-main")

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}

	var err error // Declare err here to be used for cfg loading
	cfg, err = config.LoadFromFile(cfgPath)
	if err != nil {
		log.Error("Failed to load configuration. Shutting down.", "path", cfgPath, "error", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

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

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/health", healthHandler) // Register health handler
	mux.HandleFunc("/", rootHandler)         // Keep root handler last as a catch-all for "/"

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("Server listening", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed to start or encountered an unexpected error. Shutting down.", "error", fmt.Sprintf("%+v", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	receivedSignal := <-quit
	log.Info("Received signal, initiating shutdown...", "signal", receivedSignal.String())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed (requests may have been interrupted).", "error", fmt.Sprintf("%+v", err))
	} else {
		log.Info("Server exited gracefully.")
	}
}
