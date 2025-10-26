package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// ProcessingResult holds the result of an OAuth callback processing
type ProcessingResult struct {
	Response interface{}
	Error    string
}

// DiscordOAuthHandler handles Discord OAuth operations
type DiscordOAuthHandler struct {
	discordOAuthService *services.DiscordOAuthService
	processingCodes     map[string]chan *ProcessingResult // Track codes being processed to deduplicate
	processingMu        sync.RWMutex
}

// NewDiscordOAuthHandler creates a new Discord OAuth handler
func NewDiscordOAuthHandler(discordOAuthService *services.DiscordOAuthService) *DiscordOAuthHandler {
	return &DiscordOAuthHandler{
		discordOAuthService: discordOAuthService,
		processingCodes:     make(map[string]chan *ProcessingResult),
	}
}

// OAuthCallbackRequest represents the query parameters from Discord OAuth callback
type OAuthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// OAuthCallbackResponse represents the response after OAuth callback
type OAuthCallbackResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	Message   string `json:"message"`
}

// OAuthSetupRequiredResponse represents a response when app identity setup is needed
type OAuthSetupRequiredResponse struct {
	Status           string `json:"status"`
	Message          string `json:"message"`
	AccountID        string `json:"account_id"`
	DiscordEmail     string `json:"discord_email"`
	DiscordUsername  string `json:"discord_username"`
	NeedsSetup       bool   `json:"needs_setup"`
}

// DiscordCallback handles the Discord OAuth callback
func (h *DiscordOAuthHandler) DiscordCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req OAuthCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		WriteError(w, http.StatusBadRequest, "Authorization code is required")
		return
	}

	// Check if this code is already being processed
	h.processingMu.Lock()
	if existingChannel, alreadyProcessing := h.processingCodes[req.Code]; alreadyProcessing {
		h.processingMu.Unlock()
		// Wait for the first request to complete and reuse its result
		result := <-existingChannel
		if result == nil || result.Error != "" {
			WriteError(w, http.StatusUnauthorized, "Failed to authenticate with Discord")
			return
		}
		// Return the cached response
		WriteJSON(w, http.StatusOK, result.Response)
		return
	}

	// Create a channel for this code's result
	resultChannel := make(chan *ProcessingResult, 1)
	h.processingCodes[req.Code] = resultChannel
	h.processingMu.Unlock()

	// Ensure we clean up the channel when done
	defer func() {
		h.processingMu.Lock()
		delete(h.processingCodes, req.Code)
		h.processingMu.Unlock()
		close(resultChannel)
	}()

	// Exchange code for access token
	accessToken, err := h.discordOAuthService.ExchangeCodeForToken(r.Context(), req.Code)
	if err != nil {
		resultChannel <- &ProcessingResult{Error: err.Error()} // Signal failure to other waiters
		WriteError(w, http.StatusUnauthorized, "Failed to authenticate with Discord")
		return
	}

	// Get user info from Discord
	userInfo, err := h.discordOAuthService.GetUserInfo(r.Context(), accessToken)
	if err != nil {
		resultChannel <- &ProcessingResult{Error: err.Error()} // Signal failure to other waiters
		WriteError(w, http.StatusUnauthorized, "Failed to retrieve user information from Discord")
		return
	}

	// Process Discord auth (create or login)
	account, token, err := h.discordOAuthService.ProcessDiscordAuth(r.Context(), userInfo)
	if err != nil {
		resultChannel <- &ProcessingResult{Error: err.Error()} // Signal failure to other waiters
		WriteError(w, http.StatusInternalServerError, "Failed to process authentication")
		return
	}

	// Check if setup is required (token will be "SETUP_REQUIRED" in that case)
	if token == "SETUP_REQUIRED" {
		fmt.Printf("[DISCORD_CALLBACK] Setup required detected, returning setup response for account: %s\n", account.ID)

		setupResp := OAuthSetupRequiredResponse{
			Status:          "setup_required",
			Message:         "Please set up your app identity to complete registration",
			AccountID:       account.ID.String(),
			DiscordEmail:    userInfo.Email,
			DiscordUsername: userInfo.Username,
			NeedsSetup:      true,
		}

		fmt.Printf("[DISCORD_CALLBACK] Setup response: %+v\n", setupResp)
		
		// Send response to all waiters (first request + any duplicates)
		resultChannel <- &ProcessingResult{Response: setupResp}
		
		WriteJSON(w, http.StatusOK, setupResp)
		return
	}

	// Get email - prefer the email from Discord, otherwise use the first identity's external ID
	email := userInfo.Email
	if email == "" && len(account.Identities) > 0 {
		email = account.Identities[0].ExternalID
	}

	// Get username
	username := userInfo.Username
	if username == "" && len(account.Identities) > 0 && account.Identities[0].Username != nil {
		username = *account.Identities[0].Username
	}

	// Set HTTP-only secure cookie for token
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// For now, we'll return a simple expires_at time (30 days from now)
	// In a production app, you'd calculate this from the JWT token
	expiresAt := "" // Will be set by frontend from token or calculated separately

	resp := OAuthCallbackResponse{
		ID:        account.ID.String(),
		Email:     email,
		Username:  username,
		Token:     token,
		ExpiresAt: expiresAt,
		Message:   "Authentication successful",
	}

	// Send response to all waiters (first request + any duplicates)
	resultChannel <- &ProcessingResult{Response: resp}

	WriteJSON(w, http.StatusOK, resp)
}

// CompleteDiscordSetupRequest represents the request to complete Discord OAuth setup
type CompleteDiscordSetupRequest struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Timezone  string `json:"timezone"`
}

// CompleteDiscordSetup completes the app identity setup for a Discord-only user
func (h *DiscordOAuthHandler) CompleteDiscordSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CompleteDiscordSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.AccountID) == "" {
		WriteError(w, http.StatusBadRequest, "Account ID is required")
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		WriteError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if strings.TrimSpace(req.Timezone) == "" {
		req.Timezone = "UTC" // Default to UTC
	}

	// Create app identity for the account
	token, err := h.discordOAuthService.CreateAppIdentityForDiscordAccount(
		r.Context(),
		req.AccountID,
		req.Email,
		req.Password,
		req.Timezone,
	)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to complete setup: "+err.Error())
		return
	}

	// Get the account to retrieve username from Discord identity
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Invalid account ID")
		return
	}

	account, err := h.discordOAuthService.GetAccount(r.Context(), accountID)
	if err != nil || account == nil {
		WriteError(w, http.StatusInternalServerError, "Failed to load account")
		return
	}

	// Find Discord identity to get username
	var username string
	for _, identity := range account.Identities {
		if identity.Provider == "discord" && identity.Username != nil {
			username = *identity.Username
			break
		}
	}

	// Calculate expiration time (token was created with 30 days duration)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	expiresAtStr := expiresAt.Format("2006-01-02T15:04:05Z07:00")

	// Return successful response with token and user data
	resp := OAuthCallbackResponse{
		ID:        req.AccountID,
		Email:     req.Email,
		Username:  username,
		Token:     token,
		ExpiresAt: expiresAtStr,
		Message:   "Setup completed successfully",
	}

	// Set HTTP-only secure cookie for token
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 days
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	WriteJSON(w, http.StatusOK, resp)
}
