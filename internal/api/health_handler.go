package api

import (
	"net/http"

	"github.com/ericp/chronos-bot-reminder/internal/config"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Root handles GET / and identifies the service.
func (h *HealthHandler) Root(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"service": "Chronos Reminder API",
		"version": config.Version,
		"docs":    config.URLDocs,
		"status":  config.URLStatus,
	})
}

// Health returns a simple health check response
// @Summary Health check
// @Description Returns the health status of the API
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{} "Server is healthy"
// @Router /api/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"service": "chronos-reminder-api",
		"version": config.Version,
	})
}
