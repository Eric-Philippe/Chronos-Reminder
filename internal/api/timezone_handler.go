package api

import (
	"net/http"

	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
)

// TimezoneHandler handles timezone-related HTTP requests
type TimezoneHandler struct {
	timezoneRepo repositories.TimezoneRepository
}

// NewTimezoneHandler creates a new timezone handler
func NewTimezoneHandler(timezoneRepo repositories.TimezoneRepository) *TimezoneHandler {
	return &TimezoneHandler{
		timezoneRepo: timezoneRepo,
	}
}

// GetAvailableTimezones retrieves all available timezones
// @Route: GET /api/timezones
// @Description: Returns a list of all available timezones. No authentication required.
func (h *TimezoneHandler) GetAvailableTimezones(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	timezones, err := h.timezoneRepo.GetAll()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve timezones")
		return
	}

	if len(timezones) == 0 {
		WriteJSON(w, http.StatusOK, []interface{}{})
		return
	}

	WriteJSON(w, http.StatusOK, timezones)
}
