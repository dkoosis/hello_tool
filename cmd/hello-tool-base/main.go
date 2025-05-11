// DocEnhanced: 2025-05-11
// cmd/hello-tool-base/main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	// Roach errors for base error handling, used by apperrors
	"github.com/cockroachdb/errors"

	// Local packages
	"github.com/dkoosis/hello-tool-base/internal/apperrors" // Your custom errors package
	"github.com/dkoosis/hello-tool-base/internal/buildinfo"
	"github.com/dkoosis/hello-tool-base/internal/logging" // Your custom logging package
)

var (
	// Initialize a logger for the main package
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
	// Code    int    `json:"code,omitempty"` // Standard JSON-RPC error code if needed
}

// Helper function to respond with JSON.
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// Use our structured logger
		// Wrap the error for better internal logging context
		wrappedErr := errors.Wrap(err, "respondWithJSON: failed to encode JSON response payload")
		log.Error("Failed to encode JSON response", "error", fmt.Sprintf("%+v", wrappedErr))
		// At this point, headers might have already been sent, so sending another error response might not work as expected.
		// The primary action here is to log the encoding error.
	}
}

// Helper function to respond with a JSON error to the client.
// clientMessage is the sanitized message for the client.
// internalErr is the detailed error for server-side logging.
func respondWithError(w http.ResponseWriter, statusCode int, clientMessage string, clientDetails string, internalErr error) {
	// Log the detailed internal error using our structured logger.
	// The '%+v' formatting is key for cockroachdb/errors to include stack traces.
	log.Error("Server-side error occurred",
		"internal_error", fmt.Sprintf("%+v", internalErr), // Log full internal error detail
		"client_message", clientMessage,
		"client_details", clientDetails,
		"http_status_code", statusCode,
	)

	// Respond to the client with a sanitized message.
	respondWithJSON(w, statusCode, ClientErrorResponse{Error: clientMessage, Details: clientDetails})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	reqLogger := log.WithField("handler", "helloHandler").WithField("remote_addr", r.RemoteAddr).WithField("path", r.URL.Path)
	reqLogger.Info("Received request")

	name := r.URL.Query().Get("name")
	if name == "" {
		// Use your apperrors package to create a structured error.
		internalErr := apperrors.NewInvalidParamsError(
			"helloHandler: missing 'name' query parameter",
			nil, // No underlying cause for this specific validation error
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

	// Set content type for plain text response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Helper function to handle writes and log errors
	writeAndLog := func(format string, args ...interface{}) {
		// Only proceed if no error has occurred yet in this function.
		// This is a simple way to stop writing if a previous write failed.
		// More robust error handling might involve checking w.Header().Get("Content-Length") or similar,
		// but for this simple handler, logging the first error is sufficient.
		if _, err := fmt.Fprintf(w, format, args...); err != nil {
			// Wrap the error for context
			wrappedErr := errors.Wrapf(err, "rootHandler: failed to write response for path %s", r.URL.Path)
			// Log the error using the request-specific logger
			reqLogger.Error("Failed to write to response stream", "error", fmt.Sprintf("%+v", wrappedErr))
			// Potentially set a flag to stop further writes, though Fprintf might handle subsequent calls gracefully or fail silently.
		}
	}

	writeAndLog("Hello World Root! This is the base Go service. Try /hello?name=YourName\n")
	writeAndLog("---\n")
	writeAndLog("Service: hello-tool-base\n")
	writeAndLog("Version: %s\n", buildinfo.Version)
	writeAndLog("Commit:  %s\n", buildinfo.CommitHash) // Corrected typo from previous user input (Commit:  %s\n)
	writeAndLog("Built:   %s\n", buildinfo.BuildDate)

	reqLogger.Info("Successfully processed / request")
}

func main() {
	// Setup the default logger for the application using your logging package.
	logging.SetupDefaultLogger("debug")             // "debug" level can be configured via env var in a real app.
	log = logging.GetLogger("hello-tool-base-main") // Get a logger instance for the main package

	log.Info("Starting hello-tool-base service",
		"version", buildinfo.Version,
		"commit", buildinfo.CommitHash,
		"buildDate", buildinfo.BuildDate,
	)

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", rootHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Info("Defaulting to port", "port", port)
	}

	log.Info("Server listening", "port", port)
	// Use apperrors for server start failure.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		internalErr := apperrors.NewInternalError(
			fmt.Sprintf("Server failed to start on port %s", port),
			err, // Pass the original error as cause
			map[string]interface{}{"listen_port": port},
		)
		// Use %+v to log the full error details including stack trace.
		log.Error("Fatal: Server startup failed", "error", fmt.Sprintf("%+v", internalErr))
		os.Exit(1) // Ensure application exits on fatal error
	}
}
