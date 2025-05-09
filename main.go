// hello-tool-base/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GreetingResponse struct {
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /hello") // Simple log line for now

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World" // Default name
	}

	message := fmt.Sprintf("Hello, %s, from your Go Cloud Run service!", name)
	response := GreetingResponse{Message: message}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: Failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Root handler
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		log.Println("Received request for / (root)")
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
