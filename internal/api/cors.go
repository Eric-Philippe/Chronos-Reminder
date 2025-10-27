package api

import (
	"net/http"

	"github.com/ericp/chronos-bot-reminder/internal/config"
)

// RestrictedRoutes defines routes that only allow specific origins (website only)
var RestrictedRoutes = map[string]bool{
	"/api/auth/register": true, // Website only
	"/api/auth/login":    true, // Website only
	"/api/auth/logout":   true, // Website only
	// Add more website-restricted routes here
}

// ProtectedRoutes defines routes that require authentication
var ProtectedRoutes = map[string]bool{
	"/api/reminders":       true, // Get all reminders
	"/api/reminders/{id}":  true, // Get single reminder
	"/api/reminders/errors": true, // Get reminders with errors
	"/api/account":         true, // Get account info
	// Add more authenticated routes here
}

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			isRestrictedRoute := RestrictedRoutes[r.URL.Path]
			isProtectedRoute := ProtectedRoutes[r.URL.Path]

			// Determine allowed origin based on route type and environment
			var allowedOrigin string

			if isRestrictedRoute || isProtectedRoute {
				// Restricted/Protected routes: require specific origin (configured in API_CORS)
				if cfg.Environment == "DEV" {
					// In DEV mode, allow any origin
					if origin != "" {
						allowedOrigin = origin
					} else {
						allowedOrigin = "*"
					}
				} else {
					// In production, check against configured origin
					configuredOrigin := cfg.APICors
					if configuredOrigin != "" && configuredOrigin != "*" {
						if origin == configuredOrigin {
							allowedOrigin = origin
						}
					} else if configuredOrigin == "*" {
						allowedOrigin = "*"
					}
				}
			} else {
				// Default: treat as protected routes (auth endpoints)
				if cfg.Environment == "DEV" {
					// In DEV mode, allow any origin
					if origin != "" {
						allowedOrigin = origin
					} else {
						allowedOrigin = "*"
					}
				} else {
					// In production, check against configured origin
					configuredOrigin := cfg.APICors
					if configuredOrigin != "" && configuredOrigin != "*" {
						if origin == configuredOrigin {
							allowedOrigin = origin
						}
					} else if configuredOrigin == "*" {
						allowedOrigin = "*"
					}
				}
			}

			// Always set CORS headers if we have an allowed origin
			if allowedOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// WrappedMux wraps http.ServeMux with middleware support
type WrappedMux struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

// NewWrappedMux creates a new wrapped mux
func NewWrappedMux() *WrappedMux {
	return &WrappedMux{
		mux:        http.NewServeMux(),
		middleware: []func(http.Handler) http.Handler{},
	}
}

// Use adds middleware to the mux
func (wm *WrappedMux) Use(middleware func(http.Handler) http.Handler) {
	wm.middleware = append(wm.middleware, middleware)
}

// HandleFunc registers a handler function with middleware
func (wm *WrappedMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	wm.mux.HandleFunc(pattern, handler)
}

// Handle registers a handler with middleware
func (wm *WrappedMux) Handle(pattern string, handler http.Handler) {
	wm.mux.Handle(pattern, handler)
}

// ServeHTTP implements http.Handler interface with middleware chain
func (wm *WrappedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Apply middleware in reverse order
	handler := http.Handler(wm.mux)
	for i := len(wm.middleware) - 1; i >= 0; i-- {
		handler = wm.middleware[i](handler)
	}
	handler.ServeHTTP(w, r)
}
