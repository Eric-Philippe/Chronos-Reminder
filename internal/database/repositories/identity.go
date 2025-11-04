package repositories

import (
	"errors"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// identityRepository implementation
type identityRepository struct {
	db *gorm.DB
}

// NewIdentityRepository creates a new identity repository instance
func NewIdentityRepository(db *gorm.DB) IdentityRepository {
	return &identityRepository{db: db}
}

// Identity Repository Implementation
func (r *identityRepository) Create(identity *models.Identity) error {
	return r.db.Create(identity).Error
}

func (r *identityRepository) GetByID(id uuid.UUID) (*models.Identity, error) {
	var identity models.Identity
	err := r.db.First(&identity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &identity, nil
}

func (r *identityRepository) GetByProviderAndExternalID(provider models.ProviderType, externalID string) (*models.Identity, error) {
	var identity models.Identity
	err := r.db.Where("provider = ? AND external_id = ?", provider, externalID).First(&identity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &identity, nil
}

func (r *identityRepository) GetByAccountID(accountID uuid.UUID) ([]models.Identity, error) {
	var identities []models.Identity
	err := r.db.Where("account_id = ?", accountID).Find(&identities).Error
	return identities, err
}

func (r *identityRepository) Update(identity *models.Identity) error {
	return r.db.Save(identity).Error
}

func (r *identityRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Identity{}, "id = ?", id).Error
}

func (r *identityRepository) GetByAccessToken(hashedToken string) (*models.Identity, error) {
	var identity models.Identity
	err := r.db.Where("access_token = ?", hashedToken).First(&identity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &identity, nil
}
