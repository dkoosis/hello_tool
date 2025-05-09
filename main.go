// hello-tool-base/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
		// If encoding the actual payload fails, log it but try to send a generic error.
		// This situation should be rare.
		log.Printf("ERROR: Failed to encode JSON response payload: %v", err)
		// Avoid writing to w again if headers are already sent with WriteHeader.
		// http.Error might be too late here or could corrupt the response further.
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
		// 'name' parameter is missing, and our OpenAPI spec says it's required.
		respondWithError(w, http.StatusBadRequest, "Missing 'name' query parameter", "The 'name' query parameter is required for a personalized greeting.")
		return
	}

	message := fmt.Sprintf("Hello, %s, from your Go Cloud Run service!", name)
	response := GreetingResponse{Message: message}

	// Use the helper for successful JSON response
	respondWithJSON(w, http.StatusOK, response)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Root handler
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		log.Printf("Received request for / (root) from %s", r.RemoteAddr)
		fmt.Fprintln(w, "Hello World Root! This is the base Go service. Try /hello?name=YourName")
	})

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
