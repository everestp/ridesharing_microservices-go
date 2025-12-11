package main

import (
	"encoding/json"
	"net/http"
)
func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json") // OK
	w.WriteHeader(status)                              // OK

	// FIXED: always return the encode error (you were returning directly but we mark it)
	err := json.NewEncoder(w).Encode(data) // FIXED
	if err != nil {                        // FIXED
		return err                          // FIXED
	}

	return nil // FIXED: explicitly return nil when success
}
// ADDED: Helper to send JSON error responses safely
func WriteJSONError(w http.ResponseWriter, status int, msg string) { // ADDED
	w.Header().Set("Content-Type", "application/json") // ADDED
	w.WriteHeader(status)                              // ADDED
	json.NewEncoder(w).Encode(map[string]string{       // ADDED
		"error": msg, // ADDED
	})
}
