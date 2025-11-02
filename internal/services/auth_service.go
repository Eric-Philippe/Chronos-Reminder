package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	accountRepo    repositories.AccountRepository
	identityRepo   repositories.IdentityRepository
	timezoneRepo   repositories.TimezoneRepository
}

// NewAuthService creates a new auth service instance
func NewAuthService(
	accountRepo repositories.AccountRepository,
	identityRepo repositories.IdentityRepository,
	timezoneRepo repositories.TimezoneRepository,
) *AuthService {
	return &AuthService{
		accountRepo:  accountRepo,
		identityRepo: identityRepo,
		timezoneRepo: timezoneRepo,
	}
}

// RegisterUserRequest represents user registration data
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
}

// RegisterUser creates a new app provider account with the given credentials
func (s *AuthService) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*models.Account, error) {
	if req == nil {
		return nil, errors.New("registration request is nil")
	}

	// Check if email already exists
	existingIdentity, err := s.identityRepo.GetByProviderAndExternalID(models.ProviderApp, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if existingIdentity != nil {
		return nil, errors.New("email already exists")
	}

	// Get timezone by IANA identifier
	timezone, err := s.timezoneRepo.GetByIANALocation(req.Timezone)
	if err != nil {
		return nil, fmt.Errorf("error fetching timezone: %w", err)
	}
	if timezone == nil {
		return nil, fmt.Errorf("timezone %s not found", req.Timezone)
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create new account
	account := &models.Account{
		ID:         uuid.New(),
		TimezoneID: &timezone.ID,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	// Create identity with app provider
	identity := &models.Identity{
		ID:           uuid.New(),
		AccountID:    account.ID,
		Provider:     models.ProviderApp,
		ExternalID:   req.Email,
		Username:     &req.Username,
		PasswordHash: &hashedPassword,
	}

	if err := s.identityRepo.Create(identity); err != nil {
		// Clean up the created account on identity creation failure
		s.accountRepo.Delete(account.ID)
		return nil, fmt.Errorf("error creating identity: %w", err)
	}

	// Load the full account with identities
	account.Identities = []models.Identity{*identity}
	account.Timezone = timezone

	return account, nil
}

// hashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
