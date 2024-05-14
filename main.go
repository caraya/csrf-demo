package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// WebhookData represents the JSON payload expected by the webhook
type WebhookData struct {
	Message string `json:"message"`
}

// webhookHandler handles incoming webhook requests
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON payload
	var data WebhookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}

	// Print the received data to stdout
	fmt.Printf("Received data: %s\n", data.Message)

	// Send a response
	fmt.Fprintln(w, "Webhook received!")
}

// withCORS is a middleware that sets CORS headers allowing all origins
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/webhook", webhookHandler).Methods("POST", "OPTIONS")

	// Apply the CORS middleware to all routes
	http.Handle("/", withCORS(r))

	fmt.Println("Starting server on :9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
