// cmd/hello-tool-base/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// Import the internal buildinfo package
	// Ensure MODULE_PATH in your Makefile matches this import path's root
	"github.com/dkoosis/hello-tool-base/internal/buildinfo"
)

// GreetingResponse defines the structure for a successful greeting.
type GreetingResponse struct {
	Message string `json:"message"`
}

// ErrorResponse defines the structure for a JSON error message.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"` // Optional field for more details
}

// Helper function to respond with JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("ERROR: Failed to encode JSON response payload: %v", err)
	}
}

// Helper function to respond with a JSON error
func respondWithError(w http.ResponseWriter, code int, message string, details string) {
	log.Printf("Responding with error: code=%d, message=%s, details=%s", code, message, details)
	respondWithJSON(w, code, ErrorResponse{Error: message, Details: details})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for /hello from %s", r.RemoteAddr)

	name := r.URL.Query().Get("name")
	if name == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'name' query parameter", "The 'name' query parameter is required for a personalized greeting.")
		return
	}

	message := fmt.Sprintf("Hello, %s, from your Go Cloud Run service!", name)
	response := GreetingResponse{Message: message}
	respondWithJSON(w, http.StatusOK, response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	log.Printf("Received request for / (root) from %s", r.RemoteAddr)

	// Display build information on the root path
	fmt.Fprintln(w, "Hello World Root! This is the base Go service. Try /hello?name=YourName")
	fmt.Fprintln(w, "---")
	fmt.Fprintf(w, "Service: hello-tool-base\n")
	fmt.Fprintf(w, "Version: %s\n", buildinfo.Version)
	fmt.Fprintf(w, "Commit:  %s\n", buildinfo.CommitHash)
	fmt.Fprintf(w, "Built:   %s\n", buildinfo.BuildDate)
}

func main() {
	// Log build information at startup
	log.Printf("Starting hello-tool-base")
	log.Printf("  Version:    %s", buildinfo.Version)
	log.Printf("  CommitHash: %s", buildinfo.CommitHash)
	log.Printf("  BuildDate:  %s", buildinfo.BuildDate)

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", rootHandler) // Use the new rootHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
