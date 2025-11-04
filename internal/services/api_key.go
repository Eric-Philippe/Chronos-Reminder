package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
)

// APIKeyService handles API key operations
type APIKeyService struct {
	identityRepo repositories.IdentityRepository
	accountRepo  repositories.AccountRepository
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(
	identityRepo repositories.IdentityRepository,
	accountRepo repositories.AccountRepository,
) *APIKeyService {
	return &APIKeyService{
		identityRepo: identityRepo,
		accountRepo:  accountRepo,
	}
}

// APIKeyMetadata holds metadata about an API key for responses
type APIKeyMetadata struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Scopes    string    `json:"scopes"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
	Key       string    `json:"key,omitempty"` // Only populated on creation
}

const (
	// APIKeyPrefix is the prefix for API keys (like how stripe uses sk_live_)
	APIKeyPrefix = "ck_"
	// APIKeyLength is the length of the random part of the API key
	APIKeyLength = 48
)

// GenerateAPIKey generates a new API key
func GenerateAPIKey() (string, error) {
	randomBytes := make([]byte, APIKeyLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 for readability
	randomPart := base64.URLEncoding.EncodeToString(randomBytes)
	// Remove padding to keep it shorter
	randomPart = strings.TrimRight(randomPart, "=")

	return APIKeyPrefix + randomPart, nil
}

// HashAPIKey creates a SHA256 hash of the API key
func HashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// CreateAPIKey creates a new API key for an account
func (s *APIKeyService) CreateAPIKey(accountID uuid.UUID, name string) (*APIKeyMetadata, error) {
	// Check if account already has 5 API keys
	identities, err := s.identityRepo.GetByAccountID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identities: %w", err)
	}

	apiKeyCount := 0
	for _, identity := range identities {
		if identity.Provider == models.ProviderAPIKey {
			apiKeyCount++
		}
	}

	if apiKeyCount >= 5 {
		return nil, errors.New("maximum of 5 API keys per account")
	}

	// Generate new API key
	plainKey, err := GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Hash the key for storage
	hashedKey := HashAPIKey(plainKey)

	// Create identity record
	identity := &models.Identity{
		ID:        uuid.New(),
		AccountID: accountID,
		Provider:  models.ProviderAPIKey,
		ExternalID: fmt.Sprintf("%s_%s", name, uuid.New().String()[:8]), // Make it human-readable
		Username:  &name,
		AccessToken: &hashedKey,
		Scopes:    stringPtr("reminders.read"), // Default scope
		CreatedAt: time.Now(),
	}

	if err := s.identityRepo.Create(identity); err != nil {
		return nil, fmt.Errorf("failed to create API key identity: %w", err)
	}

	return &APIKeyMetadata{
		ID:        identity.ID.String(),
		Name:      name,
		Scopes:    "reminders.read",
		CreatedAt: identity.CreatedAt,
		Key:       plainKey, // Only return on creation
	}, nil
}

// GetAPIKeys retrieves all API keys for an account (masked)
func (s *APIKeyService) GetAPIKeys(accountID uuid.UUID) ([]APIKeyMetadata, error) {
	identities, err := s.identityRepo.GetByAccountID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch identities: %w", err)
	}

	var keys []APIKeyMetadata
	for _, identity := range identities {
		if identity.Provider == models.ProviderAPIKey {
			name := ""
			if identity.Username != nil {
				name = *identity.Username
			}

			scopes := "reminders.read"
			if identity.Scopes != nil {
				scopes = *identity.Scopes
			}

			keys = append(keys, APIKeyMetadata{
				ID:        identity.ID.String(),
				Name:      name,
				Scopes:    scopes,
				CreatedAt: identity.CreatedAt,
			})
		}
	}

	return keys, nil
}

// RevokeAPIKey revokes (deletes) an API key
func (s *APIKeyService) RevokeAPIKey(accountID uuid.UUID, keyID string) error {
	keyUUID, err := uuid.Parse(keyID)
	if err != nil {
		return fmt.Errorf("invalid key ID: %w", err)
	}

	identity, err := s.identityRepo.GetByID(keyUUID)
	if err != nil {
		return fmt.Errorf("failed to fetch identity: %w", err)
	}

	if identity == nil {
		return errors.New("API key not found")
	}

	// Verify ownership
	if identity.AccountID != accountID {
		return errors.New("unauthorized")
	}

	if identity.Provider != models.ProviderAPIKey {
		return errors.New("not an API key")
	}

	// Delete the identity
	if err := s.identityRepo.Delete(keyUUID); err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	return nil
}

// ValidateAPIKey validates an API key and returns the associated account ID
func (s *APIKeyService) ValidateAPIKey(key string) (uuid.UUID, error) {
	
	if !strings.HasPrefix(key, APIKeyPrefix) {
		return uuid.Nil, errors.New("invalid API key format")
	}

	// Hash the provided key
	hashedKey := HashAPIKey(key)

	// Find identity with matching hashed key
	identity, err := s.identityRepo.GetByAccessToken(hashedKey)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	if identity == nil {
		return uuid.Nil, errors.New("API key not found")
	}

	if identity.Provider != models.ProviderAPIKey {
		return uuid.Nil, errors.New("invalid API key")
	}

	return identity.AccountID, nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
