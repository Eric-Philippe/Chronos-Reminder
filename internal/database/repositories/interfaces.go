package repositories

import (
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
)

// TimezoneRepository interface defines operations for timezone data
type TimezoneRepository interface {
	GetAll() ([]models.Timezone, error)
	GetByID(id uint) (*models.Timezone, error)
	GetByName(name string) (*models.Timezone, error)
	GetDefault() (*models.Timezone, error)
}

// AccountRepository interface defines operations for account data
type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id uuid.UUID) (*models.Account, error)
	Update(account *models.Account) error
	Delete(id uuid.UUID) error
	GetWithTimezone(id uuid.UUID) (*models.Account, error)
	GetWithIdentities(id uuid.UUID) (*models.Account, error)
}

// IdentityRepository interface defines operations for identity data
type IdentityRepository interface {
	Create(identity *models.Identity) error
	GetByID(id uuid.UUID) (*models.Identity, error)
	GetByProviderAndExternalID(provider, externalID string) (*models.Identity, error)
	GetByAccountID(accountID uuid.UUID) ([]models.Identity, error)
	Update(identity *models.Identity) error
	Delete(id uuid.UUID) error
}
