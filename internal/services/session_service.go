package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"os"

	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// SessionCacheDuration is how long to keep session data cached in Redis
	SessionCacheDuration = 24 * time.Hour
	// SessionRefreshThreshold is when we refresh the session token
	SessionRefreshThreshold = 1 * time.Hour
	// SessionCacheKeyFormat is the format for session cache keys
	SessionCacheKeyFormat = "session:%s"
)

// SessionService handles session and authentication operations
type SessionService struct {
	identityRepo repositories.IdentityRepository
	accountRepo  repositories.AccountRepository
}

// NewSessionService creates a new session service instance
func NewSessionService(
	identityRepo repositories.IdentityRepository,
	accountRepo repositories.AccountRepository,
) *SessionService {
	return &SessionService{
		identityRepo: identityRepo,
		accountRepo:  accountRepo,
	}
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

// LoginWithIDRequest represents a login request using account ID
type LoginWithIDRequest struct {
	AccountID string `json:"account_id"`
}

// SessionToken represents a JWT token payload
type SessionToken struct {
	AccountID  string    `json:"account_id"`
	IdentityID string    `json:"identity_id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	ExpiresAt  time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}

// SessionData represents cached session information
type SessionData struct {
	AccountID  uuid.UUID `json:"account_id"`
	IdentityID uuid.UUID `json:"identity_id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	ExpiresAt  time.Time `json:"expires_at"`
	RememberMe bool      `json:"remember_me"`
}

// LoginUser authenticates a user and returns a session token
func (s *SessionService) LoginUser(ctx context.Context, req *LoginRequest) (*SessionData, string, error) {
	if req == nil {
		return nil, "", errors.New("login request is nil")
	}

	// Find identity by email and provider
	identity, err := s.identityRepo.GetByProviderAndExternalID(models.ProviderApp, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("error finding identity: %w", err)
	}

	if identity == nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Verify password
	if identity.PasswordHash == nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := VerifyPassword(*identity.PasswordHash, req.Password); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Get full account with timezone
	account, err := s.accountRepo.GetWithTimezone(identity.AccountID)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching account: %w", err)
	}

	if account == nil {
		return nil, "", errors.New("account not found")
	}

	// Check if email has been verified on the account
	if !account.EmailVerified {
		return nil, "", errors.New("email not verified")
	}

	// Determine session duration based on remember_me flag
	var sessionDuration time.Duration
	if req.RememberMe {
		sessionDuration = 30 * 24 * time.Hour // 30 days
	} else {
		sessionDuration = 24 * time.Hour // 24 hours
	}

	// Create JWT token
	token, err := s.generateToken(account, identity, sessionDuration)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	// Get username
	var username string
	if identity.Username != nil {
		username = *identity.Username
	}

	// Create session data
	sessionData := &SessionData{
		AccountID:  account.ID,
		IdentityID: identity.ID,
		Email:      identity.ExternalID,
		Username:   username,
		ExpiresAt:  time.Now().Add(sessionDuration),
		RememberMe: req.RememberMe,
	}

	// Cache session in Redis
	if err := s.cacheSession(sessionData); err != nil {
		// Log but don't fail - session will still work with the token
		fmt.Printf("[SESSION] Warning: Failed to cache session: %v\n", err)
	}

	return sessionData, token, nil
}

// ValidateToken validates a JWT token and returns the session claims
func (s *SessionService) ValidateToken(tokenString string) (*SessionToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SessionToken{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*SessionToken)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if token is expired
	if claims.RegisteredClaims.ExpiresAt != nil && claims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// RefreshToken refreshes an existing token
func (s *SessionService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Calculate remaining time
	now := time.Now()
	expiresAt := claims.RegisteredClaims.ExpiresAt
	if expiresAt == nil {
		return tokenString, nil // No expiration, no need to refresh
	}

	timeRemaining := expiresAt.Time.Sub(now)

	// Only refresh if less than threshold remains
	if timeRemaining > SessionRefreshThreshold {
		return tokenString, nil // Token is still valid, no need to refresh
	}

	// Generate new token with extended expiration
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, SessionToken{
		AccountID:  claims.AccountID,
		IdentityID: claims.IdentityID,
		Email:      claims.Email,
		Username:   claims.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	return newToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// LogoutUser invalidates a session
func (s *SessionService) LogoutUser(accountID uuid.UUID) error {
	return s.invalidateSession(accountID)
}

// LoginUserWithID authenticates a user using account ID (used for email verification)
func (s *SessionService) LoginUserWithID(ctx context.Context, req *LoginWithIDRequest) (*SessionData, string, error) {
	if req == nil {
		return nil, "", errors.New("login request is nil")
	}

	if req.AccountID == "" {
		return nil, "", errors.New("account ID is required")
	}

	// Parse account ID
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, "", errors.New("invalid account ID")
	}

	// Get full account with timezone and identities
	account, err := s.accountRepo.GetWithIdentities(accountID)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching account: %w", err)
	}

	if account == nil {
		return nil, "", errors.New("account not found")
	}

	// Get the app identity - find the first app provider identity
	var identity *models.Identity
	for _, id := range account.Identities {
		if id.Provider == models.ProviderApp {
			identity = &id
			break
		}
	}

	if identity == nil {
		return nil, "", errors.New("identity not found")
	}

	// Check if email has been verified on the account
	if !account.EmailVerified {
		return nil, "", errors.New("email not verified")
	}

	// Create JWT token (24 hour session)
	token, err := s.generateToken(account, identity, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	// Get username
	var username string
	if identity.Username != nil {
		username = *identity.Username
	}

	// Create session data
	sessionData := &SessionData{
		AccountID:  account.ID,
		IdentityID: identity.ID,
		Email:      identity.ExternalID,
		Username:   username,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		RememberMe: false,
	}

	// Cache session in Redis
	if err := s.cacheSession(sessionData); err != nil {
		// Log but don't fail - session will still work with the token
		fmt.Printf("[SESSION] Warning: Failed to cache session: %v\n", err)
	}

	return sessionData, token, nil
}

// generateToken creates a JWT token
func (s *SessionService) generateToken(account *models.Account, identity *models.Identity, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	username := ""
	if identity.Username != nil {
		username = *identity.Username
	}

	claims := SessionToken{
		AccountID:  account.ID.String(),
		IdentityID: identity.ID.String(),
		Email:      identity.ExternalID,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// cacheSession stores session data in Redis
func (s *SessionService) cacheSession(session *SessionData) error {
	cacheKey := fmt.Sprintf(SessionCacheKeyFormat, session.AccountID.String())
	return database.SetCache(cacheKey, session, SessionCacheDuration)
}

// GetCachedSession retrieves cached session data
func (s *SessionService) GetCachedSession(accountID uuid.UUID) (*SessionData, error) {
	cacheKey := fmt.Sprintf(SessionCacheKeyFormat, accountID.String())
	var session SessionData
	err := database.GetCache(cacheKey, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// invalidateSession removes session data from Redis
func (s *SessionService) invalidateSession(accountID uuid.UUID) error {
	cacheKey := fmt.Sprintf(SessionCacheKeyFormat, accountID.String())
	return database.DeleteCache(cacheKey)
}

// GenerateSessionID generates a random session ID
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateTokenForAccount creates a JWT token directly for an account and identity
// This is used for OAuth flows where we don't have a password
func (s *SessionService) generateTokenForAccount(account *models.Account, identity *models.Identity, duration time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	username := ""
	if identity.Username != nil {
		username = *identity.Username
	}

	claims := SessionToken{
		AccountID:  account.ID.String(),
		IdentityID: identity.ID.String(),
		Email:      identity.ExternalID,
		Username:   username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
