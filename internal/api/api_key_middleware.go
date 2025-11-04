package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// contextKey is a custom type for context keys
type contextKeyAPIKeyAuth string

const APIKeyAuthKey contextKeyAPIKeyAuth = "api_key_auth"

// APIKeyAuthMiddleware creates middleware that validates API keys
func APIKeyAuthMiddleware(apiKeyService *services.APIKeyService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get API key from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, http.StatusUnauthorized, "No API key provided")
				return
			}

			// Extract Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				WriteError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			apiKey := parts[1]

			// Validate the API key
			accountID, err := apiKeyService.ValidateAPIKey(apiKey)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "Invalid API key")
				return
			}

			// Add account ID to request context
			ctx := context.WithValue(r.Context(), AccountIDKey, accountID)
			ctx = context.WithValue(ctx, APIKeyAuthKey, true)
			*r = *r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// IsAPIKeyAuth checks if the request was authenticated via API key
func IsAPIKeyAuth(r *http.Request) bool {
	val := r.Context().Value(APIKeyAuthKey)
	if val == nil {
		return false
	}
	return val.(bool)
}

// APIKeyScopeMiddleware validates that API key has the required scope for read operations
func APIKeyScopeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce scope restrictions for API keys
		if !IsAPIKeyAuth(r) {
			next.ServeHTTP(w, r)
			return
		}

		// API keys are only allowed for GET requests on specific endpoints
		if r.Method != http.MethodGet {
			WriteError(w, http.StatusForbidden, "API keys can only be used for read operations")
			return
		}

		// Check if the endpoint is allowed for API keys
		path := r.URL.Path
		allowedPaths := []string{
			"/api/reminders",           // GET /api/reminders - list all reminders
			"/api/reminders/",          // GET /api/reminders/{id} - get single reminder
		}

		allowed := false
		for _, p := range allowedPaths {
			if strings.HasPrefix(path, p) {
				allowed = true
				break
			}
		}

		if !allowed {
			WriteError(w, http.StatusForbidden, "API key does not have permission for this endpoint")
			return
		}

		next.ServeHTTP(w, r)
	})
}
