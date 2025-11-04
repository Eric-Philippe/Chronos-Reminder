package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// RateLimitMiddleware creates a middleware that enforces per-session rate limits
// Sessions are identified by their JWT token (from Authorization header or cookie)
func RateLimitMiddleware(rateLimiter *services.RateLimiterService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the session identifier from the token
			sessionID, err := extractSessionID(r)
			if err != nil {
				// If we can't extract session ID, deny the request
				writeRateLimitError(w, http.StatusUnauthorized, "Unable to identify session")
				return
			}

			// Check rate limit for this session
			allowed, remaining, resetTime, err := rateLimiter.CheckRateLimit(r.Context(), sessionID)
			if err != nil {
				// Log error but allow request to proceed (fail open)
				fmt.Printf("[RATE_LIMIT] Error checking rate limit: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rateLimiter.RequestsPerWindow()))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

			if !allowed {
				writeRateLimitError(w,
					http.StatusTooManyRequests,
					fmt.Sprintf("Rate limit exceeded. Reset at %s", resetTime.Format("2006-01-02 15:04:05 MST")),
				)
				return
			}

			// Request allowed, continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// extractSessionID extracts a unique session identifier from the request
// This uses the JWT token as the session identifier
func extractSessionID(r *http.Request) (string, error) {
	// Try to get token from Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1], nil
		}
	}

	// If no token in header, try cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", fmt.Errorf("no session token found")
}

// writeRateLimitError writes a rate limit error response
func writeRateLimitError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
		"code":  "RATE_LIMIT_EXCEEDED",
	})
}
