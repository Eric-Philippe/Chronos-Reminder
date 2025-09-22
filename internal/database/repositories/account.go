package repositories

import (
	"errors"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// accountRepository implementation
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository instance
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

// Account Repository Implementation
func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Account{}, "id = ?", id).Error
}

func (r *accountRepository) GetWithTimezone(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("Timezone").First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetWithIdentities(id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("Identities").Preload("Timezone").First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
