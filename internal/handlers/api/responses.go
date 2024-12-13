package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// jsonResponse represents a standardized JSON response.
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

// writeJSON writes a JSON response to the HTTP response writer.
func writeJSON(w http.ResponseWriter, status int, data interface{}, logger *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Println("Failed to write JSON response:", err)
	}
}
