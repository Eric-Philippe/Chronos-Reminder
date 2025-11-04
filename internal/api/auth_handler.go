package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// Ensure RegisterUserRequest is imported from services package

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService         *services.AuthService
	sessionService      *services.SessionService
	verificationService *services.VerificationService
	webAppURL           string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	authService *services.AuthService,
	sessionService *services.SessionService,
	verificationService *services.VerificationService,
	webAppURL string,
) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		sessionService:      sessionService,
		verificationService: verificationService,
		webAppURL:           webAppURL,
	}
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timezone string `json:"timezone"` // IANA timezone identifier (e.g., "America/New_York")
}

// RegisterResponse represents the registration response payload
type RegisterResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// VerifyEmailRequest represents the email verification request payload
type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// VerifyEmailResponse represents the email verification response payload
type VerifyEmailResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	Message   string `json:"message"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	Message   string `json:"message"`
}

// Register handles user registration for the app provider
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if err := validateRegisterRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert API request to service request
	serviceReq := &services.RegisterUserRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Timezone: req.Timezone,
	}

	// Register the user
	account, err := h.authService.RegisterUser(r.Context(), serviceReq)
	if err != nil {
		// Check for specific error types
		if strings.Contains(err.Error(), "already exists") {
			WriteError(w, http.StatusConflict, err.Error())
			return
		}
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	// Create verification record and send verification email
	verificationCode, err := h.verificationService.CreateVerification(req.Email, account.ID.String())
	if err != nil {
		// Log error but don't fail registration
		WriteError(w, http.StatusInternalServerError, "Failed to create verification code")
		return
	}

	// Build verification link (frontend will handle the redirect)
	verificationLink := h.webAppURL + "/verify?email=" + req.Email + "&code=" + verificationCode

	// Send verification email
	_, err = h.verificationService.SendVerificationEmail(req.Email, verificationCode, verificationLink)
	if err != nil {
		// Log error but don't fail registration
		WriteError(w, http.StatusInternalServerError, "Failed to send verification email")
		return
	}

	// Get the email identity for response
	var email string
	if len(account.Identities) > 0 {
		email = account.Identities[0].ExternalID
	}

	var username string
	if len(account.Identities) > 0 && account.Identities[0].Username != nil {
		username = *account.Identities[0].Username
	}

	resp := RegisterResponse{
		ID:       account.ID.String(),
		Email:    email,
		Username: username,
		Message:  "Account created successfully. Please check your email to verify your account.",
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// Login handles user login for the app provider
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if err := validateLoginRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Convert API request to service request
	serviceReq := &services.LoginRequest{
		Email:      req.Email,
		Password:   req.Password,
		RememberMe: req.RememberMe,
	}

	// Authenticate user
	sessionData, token, err := h.sessionService.LoginUser(r.Context(), serviceReq)
	if err != nil {
		// Check for specific error types
		if strings.Contains(err.Error(), "email not verified") {
			WriteError(w, http.StatusForbidden, "Email not verified. Please check your email to verify your account.")
			return
		}
		if strings.Contains(err.Error(), "invalid email or password") {
			WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		if strings.Contains(err.Error(), "not found") {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	// Set HTTP-only secure cookie for token
	maxAge := 24 * 3600 // 24 hours
	if req.RememberMe {
		maxAge = 30 * 24 * 3600 // 30 days
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	resp := LoginResponse{
		ID:        sessionData.AccountID.String(),
		Email:     sessionData.Email,
		Username:  sessionData.Username,
		Token:     token,
		ExpiresAt: sessionData.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		Message:   "Login successful",
	}

	WriteJSON(w, http.StatusOK, resp)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get token from cookie to extract account ID
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		WriteError(w, http.StatusBadRequest, "No session found")
		return
	}

	// Validate token and get claims
	claims, err := h.sessionService.ValidateToken(cookie.Value)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	// Parse account ID
	accountID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	// Logout user (invalidate session)
	if err := h.sessionService.LogoutUser(accountID); err != nil {
		// Log but don't fail - cookie will be cleared anyway
	}

	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

// VerifyEmail handles email verification with verification code
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if strings.TrimSpace(req.Email) == "" {
		WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	if strings.TrimSpace(req.Code) == "" {
		WriteError(w, http.StatusBadRequest, "Verification code is required")
		return
	}

	// Verify email
	accountID, err := h.verificationService.VerifyEmail(req.Email, req.Code)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, "Invalid or expired verification code")
		return
	}

	// Create session and login user
	loginReq := &services.LoginWithIDRequest{
		AccountID: accountID,
	}

	sessionData, token, err := h.sessionService.LoginUserWithID(r.Context(), loginReq)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create session after verification")
		return
	}

	// Set HTTP-only secure cookie for token
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 3600, // 24 hours
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Delete verification records
	_ = h.verificationService.DeleteVerification(req.Email)

	resp := VerifyEmailResponse{
		ID:        sessionData.AccountID.String(),
		Email:     sessionData.Email,
		Username:  sessionData.Username,
		Token:     token,
		ExpiresAt: sessionData.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		Message:   "Email verified successfully",
	}

	WriteJSON(w, http.StatusOK, resp)
}

// validateLoginRequest validates the login request
func validateLoginRequest(req *LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return ErrEmailRequired
	}

	if !isValidEmail(req.Email) {
		return ErrInvalidEmail
	}

	if strings.TrimSpace(req.Password) == "" {
		return ErrPasswordRequired
	}

	return nil
}

// validateRegisterRequest validates the registration request
func validateRegisterRequest(req *RegisterRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return ErrEmailRequired
	}

	if !isValidEmail(req.Email) {
		return ErrInvalidEmail
	}

	if strings.TrimSpace(req.Username) == "" {
		return ErrUsernameRequired
	}

	if len(req.Username) > 128 {
		return ErrUsernameTooLong
	}

	if strings.TrimSpace(req.Password) == "" {
		return ErrPasswordRequired
	}

	if len(req.Password) < 8 {
		return ErrPasswordTooShort
	}

	if strings.TrimSpace(req.Timezone) == "" {
		return ErrTimezoneRequired
	}

	return nil
}

// isValidEmail performs a basic email validation
func isValidEmail(email string) bool {
	// Simple email validation - checks for @ and domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return false
	}
	// Check for at least one dot in domain
	if !strings.Contains(parts[1], ".") {
		return false
	}
	return true
}
