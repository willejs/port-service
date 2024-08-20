package http

import (
	"encoding/json"
	"net/http"
)

// create a simple /ready handler
func (h *PortHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]bool{"ready": true}); err != nil {
		// Return an error
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
