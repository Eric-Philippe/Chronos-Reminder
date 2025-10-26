package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/ericp/chronos-bot-reminder/internal/docs"
	"github.com/ericp/chronos-bot-reminder/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the API server
type Server struct {
	mux    *WrappedMux
	port   string
	server *http.Server
	cfg    *config.Config
}

// NewServer creates a new API server instance
func NewServer(cfg *config.Config, repos *repositories.Repositories) *Server {
	// Initialize services
	authService := services.NewAuthService(
		repos.Account,
		repos.Identity,
		repos.Timezone,
	)

	sessionService := services.NewSessionService(
		repos.Identity,
		repos.Account,
	)

	discordOAuthService := services.NewDiscordOAuthService(
		cfg.DiscordClientID,
		cfg.DiscordClientSecret,
		cfg.DiscordRedirectURI,
		repos.Identity,
		repos.Account,
		repos.Timezone,
		sessionService,
	)

	// Initialize handlers
	authHandler := NewAuthHandler(authService, sessionService)
	discordOAuthHandler := NewDiscordOAuthHandler(discordOAuthService)

	// Create wrapped mux with CORS middleware
	wrappedMux := NewWrappedMux()
	wrappedMux.Use(CORSMiddleware(cfg))

	// Register all routes
	registerSwaggerRoutes(wrappedMux)
	registerAuthRoutes(wrappedMux, authHandler)
	registerDiscordOAuthRoutes(wrappedMux, discordOAuthHandler)

	return &Server{
		mux:  wrappedMux,
		port: cfg.APIPort,
		cfg:  cfg,
	}
}

// registerSwaggerRoutes registers Swagger documentation routes
func registerSwaggerRoutes(mux *WrappedMux) {
	// Swagger UI main page
	swaggerHandler := httpSwagger.Handler()

	mux.Handle("GET /swagger/", swaggerHandler)
	mux.Handle("GET /swagger/*", swaggerHandler)

	// Swagger JSON endpoint
	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(docs.ReadDoc()))
	})

	// Redirect /swagger/index.html to /swagger/
	mux.HandleFunc("GET /swagger/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
}

// registerAuthRoutes registers authentication routes
func registerAuthRoutes(mux *WrappedMux, authHandler *AuthHandler) {
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)
}

// registerDiscordOAuthRoutes registers Discord OAuth routes
func registerDiscordOAuthRoutes(mux *WrappedMux, discordOAuthHandler *DiscordOAuthHandler) {
	mux.HandleFunc("POST /api/auth/discord/callback", discordOAuthHandler.DiscordCallback)
	mux.HandleFunc("POST /api/auth/discord/setup", discordOAuthHandler.CompleteDiscordSetup)
}

// Start starts the API server and listens for incoming requests
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.mux,
	}

	log.Printf("[API] - 🚀 Starting API server on port %s\n", s.port)
	log.Printf("[API] - 📡 Server running at http://localhost:%s\n", s.port)
	log.Printf("[API] - 📚 Swagger documentation available at http://localhost:%s/swagger/\n", s.port)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("[API] - ❌ Failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the API server
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	log.Println("[API] - 🛑 Shutting down API server...")
	return s.server.Close()
}

// GetPort returns the port the server is listening on
func (s *Server) GetPort() string {
	return s.port
}
