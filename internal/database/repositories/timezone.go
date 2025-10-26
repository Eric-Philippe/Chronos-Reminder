package repositories

import (
	"errors"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"gorm.io/gorm"
)

// timezoneRepository implementation
type timezoneRepository struct {
	db *gorm.DB
}

// NewTimezoneRepository creates a new timezone repository instance
func NewTimezoneRepository(db *gorm.DB) TimezoneRepository {
	return &timezoneRepository{db: db}
}

// Timezone Repository Implementation
func (r *timezoneRepository) GetAll() ([]models.Timezone, error) {
	var timezones []models.Timezone
	err := r.db.Find(&timezones).Error
	return timezones, err
}

func (r *timezoneRepository) GetByID(id uint) (*models.Timezone, error) {
	var timezone models.Timezone
	err := r.db.First(&timezone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &timezone, nil
}

func (r *timezoneRepository) GetByName(name string) (*models.Timezone, error) {
	var timezone models.Timezone
	err := r.db.Where("name = ?", name).First(&timezone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &timezone, nil
}

func (r *timezoneRepository) GetByIANALocation(ianaLocation string) (*models.Timezone, error) {
	var timezone models.Timezone
	err := r.db.Where("iana_location = ?", ianaLocation).First(&timezone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &timezone, nil
}

func (r *timezoneRepository) GetDefault() (*models.Timezone, error) {
	var timezone models.Timezone
	defaultTzId := config.GetDatabaseConfig().DefaultTZ
	err := r.db.Where("id = ?", defaultTzId).First(&timezone).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &timezone, nil
}