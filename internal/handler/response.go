package handler

import (
	"encoding/json"
	"net/http"
)

// errorResponse is the standard error payload.
type errorResponse struct {
	Error string `json:"error"`
}

// respondJSON writes status + JSON body.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// respondError writes a JSON error payload.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, errorResponse{Error: message})
}
