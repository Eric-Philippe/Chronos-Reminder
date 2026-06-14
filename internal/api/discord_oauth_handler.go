package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
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
	repos               *repositories.Repositories
	processingCodes     map[string]chan *ProcessingResult // Track codes being processed to deduplicate
	processingMu        sync.RWMutex
}

// NewDiscordOAuthHandler creates a new Discord OAuth handler
func NewDiscordOAuthHandler(discordOAuthService *services.DiscordOAuthService, repos *repositories.Repositories) *DiscordOAuthHandler {
	return &DiscordOAuthHandler{
		discordOAuthService: discordOAuthService,
		repos:               repos,
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
	accessToken, refreshToken, err := h.discordOAuthService.ExchangeCodeForToken(r.Context(), req.Code)
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
	account, token, err := h.discordOAuthService.ProcessDiscordAuth(r.Context(), userInfo, accessToken, refreshToken)
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

// LinkDiscordRequest represents the request to link a Discord identity to the
// currently authenticated account.
type LinkDiscordRequest struct {
	Code string `json:"code"`
}

// LinkDiscordResponse represents the response after linking a Discord identity.
type LinkDiscordResponse struct {
	Message         string `json:"message"`
	Username        string `json:"username"`
	MergeRequired   bool   `json:"merge_required,omitempty"`
	OtherAccountID  string `json:"other_account_id,omitempty"`
	DiscordUsername string `json:"discord_username,omitempty"`
	MergeToken      string `json:"merge_token,omitempty"` // Short-lived signed token; submit to /api/account/merge
}

// MergeAccountsRequest represents the request to merge another account into the current one.
type MergeAccountsRequest struct {
	MergeToken string `json:"merge_token"` // Signed token from the link conflict response
}

// generateMergeToken creates an HMAC-signed token encoding the allowed survivor→merged merge.
// The token is valid for 5 minutes.
func generateMergeToken(survivorID, mergedID uuid.UUID) string {
	exp := strconv.FormatInt(time.Now().Add(5*time.Minute).Unix(), 10)
	payload := survivorID.String() + ":" + mergedID.String() + ":" + exp
	secret := []byte(config.Load().JWTSecret)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + sig
}

// verifyMergeToken validates a merge token and returns survivor+merged IDs.
// Returns an error if the token is invalid or expired.
func verifyMergeToken(token string) (survivorID, mergedID uuid.UUID, err error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid merge token format")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid merge token encoding")
	}
	payload := string(payloadBytes)
	secret := []byte(config.Load().JWTSecret)
	mac := hmac.New(sha256.New, secret)
	mac.Write(payloadBytes)
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[1]), []byte(expectedSig)) {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid merge token signature")
	}
	segments := strings.Split(payload, ":")
	if len(segments) != 3 {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid merge token payload")
	}
	exp, err := strconv.ParseInt(segments[2], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return uuid.Nil, uuid.Nil, fmt.Errorf("merge token expired")
	}
	survivorID, err = uuid.Parse(segments[0])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid survivor id in merge token")
	}
	mergedID, err = uuid.Parse(segments[1])
	if err != nil {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid merged id in merge token")
	}
	return survivorID, mergedID, nil
}

// LinkDiscordIdentity links a Discord account to the authenticated account.
// Unlike DiscordCallback (which logs in / signs up), this requires an existing
// authenticated session and attaches the Discord identity to that account.
func (h *DiscordOAuthHandler) LinkDiscordIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accountID, ok := r.Context().Value(AccountIDKey).(uuid.UUID)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req LinkDiscordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		WriteError(w, http.StatusBadRequest, "Authorization code is required")
		return
	}

	// Exchange code for access token
	accessToken, refreshToken, err := h.discordOAuthService.ExchangeCodeForToken(r.Context(), req.Code)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Failed to authenticate with Discord")
		return
	}

	// Get user info from Discord
	userInfo, err := h.discordOAuthService.GetUserInfo(r.Context(), accessToken)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Failed to retrieve user information from Discord")
		return
	}

	// Link Discord identity to the authenticated account
	result, err := h.discordOAuthService.LinkDiscordToAccount(r.Context(), accountID, userInfo, accessToken, refreshToken)
	if err != nil {
		if errors.Is(err, services.ErrDiscordLinkedToOtherAccount) {
			// Offer a merge rather than a hard error
			mergeToken := generateMergeToken(accountID, result.OtherAccountID)
			WriteJSON(w, http.StatusOK, LinkDiscordResponse{
				Message:         "This Discord account belongs to another Chronos account. Confirm to merge.",
				MergeRequired:   true,
				OtherAccountID:  result.OtherAccountID.String(),
				DiscordUsername: result.OtherDiscordUsername,
				MergeToken:      mergeToken,
			})
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to link Discord account")
		return
	}

	WriteJSON(w, http.StatusOK, LinkDiscordResponse{
		Message:  "Discord account linked successfully",
		Username: userInfo.Username,
	})
}

// CompleteDiscordSetupRequest represents the request to complete Discord OAuth setup
type CompleteDiscordSetupRequest struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
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

	if strings.TrimSpace(req.Username) == "" {
		WriteError(w, http.StatusBadRequest, "Username is required")
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		WriteError(w, http.StatusBadRequest, "Password is required")
		return
	}

	if strings.TrimSpace(req.Timezone) == "" {
		req.Timezone = "UTC" // Default to UTC
	}

	fmt.Printf("[DISCORD_SETUP] Request received - Email: %s, Username: %s, Timezone: %s\n", req.Email, req.Username, req.Timezone)

	// Create app identity for the account
	token, err := h.discordOAuthService.CreateAppIdentityForDiscordAccount(
		r.Context(),
		req.AccountID,
		req.Email,
		req.Username,
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

// MergeDiscordAccounts merges another account into the currently authenticated account.
// The merge_token (issued by the Discord link conflict response) proves the merge was
// initiated legitimately and specifies which accounts are involved.
// @Route: POST /api/account/merge
func (h *DiscordOAuthHandler) MergeDiscordAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	callerID, ok := r.Context().Value(AccountIDKey).(uuid.UUID)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req MergeAccountsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(req.MergeToken) == "" {
		WriteError(w, http.StatusBadRequest, "merge_token is required")
		return
	}

	survivorID, mergedID, err := verifyMergeToken(req.MergeToken)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid or expired merge token: "+err.Error())
		return
	}

	// The caller must be the survivor specified in the token
	if callerID != survivorID {
		WriteError(w, http.StatusForbidden, "Merge token was not issued for your account")
		return
	}

	if err := services.MergeAccounts(r.Context(), h.repos, survivorID, mergedID); err != nil {
		fmt.Printf("[MERGE] Failed to merge accounts %s <- %s: %v\n", survivorID, mergedID, err)
		WriteError(w, http.StatusInternalServerError, "Failed to merge accounts")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Accounts merged successfully",
	})
}
