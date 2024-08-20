package http

import (
	"encoding/json"
	"net/http"

	"github.com/willejs/ports-service/internal/controller"
)

// PortHandler handles HTTP requests for ports.
type PortHandler struct {
	controller *controller.PortController
}

// NewPortHandler creates a new PortHandler.
func NewPortHandler(controller *controller.PortController) *PortHandler {
	return &PortHandler{controller: controller}
}

// ListPorts handles the /ports endpoint.
func (h *PortHandler) ListPorts(w http.ResponseWriter, r *http.Request) {
	ports, err := h.controller.ListAllPorts()
	if err != nil {
		http.Error(w, "Failed to list ports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Encode the ports to JSON and check for errors
	if err := json.NewEncoder(w).Encode(ports); err != nil {
		// return an error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
