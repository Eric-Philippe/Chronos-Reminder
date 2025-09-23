package repositories

import (
	"errors"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// reminderDestinationRepository implementation
type reminderDestinationRepository struct {
	db *gorm.DB
}

// NewReminderDestinationRepository creates a new reminder destination repository instance
func NewReminderDestinationRepository(db *gorm.DB) ReminderDestinationRepository {
	return &reminderDestinationRepository{db: db}
}

// ReminderDestination Repository Implementation
func (r *reminderDestinationRepository) Create(destination *models.ReminderDestination) error {
	return r.db.Create(destination).Error
}

func (r *reminderDestinationRepository) GetByID(id uuid.UUID) (*models.ReminderDestination, error) {
	var destination models.ReminderDestination
	err := r.db.First(&destination, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &destination, nil
}

func (r *reminderDestinationRepository) GetByReminderID(reminderID uuid.UUID) ([]models.ReminderDestination, error) {
	var destinations []models.ReminderDestination
	err := r.db.Where("reminder_id = ?", reminderID).Find(&destinations).Error
	return destinations, err
}

func (r *reminderDestinationRepository) GetByReminderIDWithReminder(reminderID uuid.UUID) ([]models.ReminderDestination, error) {
	var destinations []models.ReminderDestination
	err := r.db.Preload("Reminder").Where("reminder_id = ?", reminderID).Find(&destinations).Error
	return destinations, err
}

func (r *reminderDestinationRepository) GetByType(destinationType models.DestinationType) ([]models.ReminderDestination, error) {
	var destinations []models.ReminderDestination
	err := r.db.Where("type = ?", destinationType).Find(&destinations).Error
	return destinations, err
}

func (r *reminderDestinationRepository) Update(destination *models.ReminderDestination) error {
	return r.db.Save(destination).Error
}

func (r *reminderDestinationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ReminderDestination{}, "id = ?", id).Error
}

func (r *reminderDestinationRepository) DeleteByReminderID(reminderID uuid.UUID) error {
	return r.db.Where("reminder_id = ?", reminderID).Delete(&models.ReminderDestination{}).Error
}

func (r *reminderDestinationRepository) CreateMultiple(destinations []models.ReminderDestination) error {
	return r.db.Create(&destinations).Error
}

func (r *reminderDestinationRepository) GetByMetadataField(field string, value interface{}) ([]models.ReminderDestination, error) {
	var destinations []models.ReminderDestination
	err := r.db.Where("metadata ->> ? = ?", field, value).Find(&destinations).Error
	return destinations, err
}
