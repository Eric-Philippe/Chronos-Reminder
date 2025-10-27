package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const AccountIDKey contextKey = "account_id"

// Handler wraps all API handlers
type Handler struct {
	authHandler     *AuthHandler
	reminderHandler *ReminderHandler
	sessionService  *services.SessionService
}

// NewHandler creates a new API handler
func NewHandler(
	authHandler *AuthHandler,
	reminderHandler *ReminderHandler,
	sessionService *services.SessionService,
) *Handler {
	return &Handler{
		authHandler:     authHandler,
		reminderHandler: reminderHandler,
		sessionService:  sessionService,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth routes (no auth middleware needed)
	mux.HandleFunc("POST /api/auth/register", h.authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", h.authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", h.authHandler.Logout)

	// Reminder routes (with auth middleware)
	mux.HandleFunc("GET /api/reminders/{id}", h.authMiddlewareHandler(h.reminderHandler.GetReminder))
	mux.HandleFunc("PUT /api/reminders/{id}", h.authMiddlewareHandler(h.reminderHandler.UpdateReminder))
	mux.HandleFunc("DELETE /api/reminders/{id}", h.authMiddlewareHandler(h.reminderHandler.DeleteReminder))
	mux.HandleFunc("POST /api/reminders/{id}/pause", h.authMiddlewareHandler(h.reminderHandler.PauseReminder))
	mux.HandleFunc("POST /api/reminders/{id}/resume", h.authMiddlewareHandler(h.reminderHandler.ResumeReminder))
	mux.HandleFunc("POST /api/reminders/{id}/duplicate", h.authMiddlewareHandler(h.reminderHandler.DuplicateReminder))

	return mux
}

// authMiddlewareHandler wraps a handler with JWT authentication
func (h *Handler) authMiddlewareHandler(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		var token string

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		// If no token in header, try cookie
		if token == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			WriteError(w, http.StatusUnauthorized, "No authentication token found")
			return
		}

		// Validate token and extract claims
		claims, err := h.sessionService.ValidateToken(token)
		if err != nil {
			WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Extract account ID from claims
		accountID, err := uuid.Parse(claims.AccountID)
		if err != nil {
			WriteError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Add account ID to request context
		ctx := context.WithValue(r.Context(), AccountIDKey, accountID)
		*r = *r.WithContext(ctx)

		// Call the actual handler
		handler(w, r)
	}
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
