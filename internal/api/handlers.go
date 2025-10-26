package api

import (
	"encoding/json"
	"net/http"
)

// Handler wraps all API handlers
type Handler struct {
	authHandler *AuthHandler
}

// NewHandler creates a new API handler
func NewHandler(authHandler *AuthHandler) *Handler {
	return &Handler{
		authHandler: authHandler,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /api/auth/register", h.authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", h.authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", h.authHandler.Logout)

	return mux
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, map[string]string{
		"error": message,
	})
}
