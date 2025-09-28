package repositories

import (
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// reminderErrorRepository implementation
type reminderErrorRepository struct {
	db *gorm.DB
}

// NewReminderErrorRepository creates a new reminder error repository instance
func NewReminderErrorRepository(db *gorm.DB) ReminderErrorRepository {
	return &reminderErrorRepository{db: db}
}

// Create creates a new reminder error
func (r *reminderErrorRepository) Create(reminderError *models.ReminderError) error {
	return r.db.Create(reminderError).Error
}

// GetByID retrieves a reminder error by ID
func (r *reminderErrorRepository) GetByID(id uuid.UUID) (*models.ReminderError, error) {
	var reminderError models.ReminderError
	err := r.db.Preload("Reminder").
		Preload("ReminderDestination").
		First(&reminderError, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reminderError, nil
}

// GetByReminderID retrieves all reminder errors for a specific reminder
func (r *reminderErrorRepository) GetByReminderID(reminderID uuid.UUID) ([]models.ReminderError, error) {
	var reminderErrors []models.ReminderError
	err := r.db.Preload("Reminder").
		Preload("ReminderDestination").
		Where("reminder_id = ?", reminderID).
		Order("timestamp DESC").
		Find(&reminderErrors).Error
	return reminderErrors, err
}

// GetByReminderDestinationID retrieves all reminder errors for a specific reminder destination
func (r *reminderErrorRepository) GetByReminderDestinationID(reminderDestinationID uuid.UUID) ([]models.ReminderError, error) {
	var reminderErrors []models.ReminderError
	err := r.db.Preload("Reminder").
		Preload("ReminderDestination").
		Where("reminder_destination_id = ?", reminderDestinationID).
		Order("timestamp DESC").
		Find(&reminderErrors).Error
	return reminderErrors, err
}

// GetByDateRange retrieves all reminder errors within a date range
func (r *reminderErrorRepository) GetByDateRange(startDate, endDate time.Time) ([]models.ReminderError, error) {
	var reminderErrors []models.ReminderError
	err := r.db.Preload("Reminder").
		Preload("ReminderDestination").
		Where("timestamp BETWEEN ? AND ?", startDate, endDate).
		Order("timestamp DESC").
		Find(&reminderErrors).Error
	return reminderErrors, err
}

// Delete removes a reminder error by ID
func (r *reminderErrorRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ReminderError{}, "id = ?", id).Error
}

// GetUnfixedByReminderID retrieves all unfixed errors for a specific reminder
func (r *reminderErrorRepository) GetUnfixedByReminderID(reminderID uuid.UUID) ([]models.ReminderError, error) {
	var reminderErrors []models.ReminderError
	err := r.db.Preload("Reminder").
		Preload("ReminderDestination").
		Where("reminder_id = ? AND fixed = false", reminderID).
		Order("timestamp DESC").
		Find(&reminderErrors).Error
	return reminderErrors, err
}

// GetUnfixedByReminderDestinationID retrieves all unfixed errors for a specific reminder destination
func (r *reminderErrorRepository) GetUnfixedByReminderDestinationID(reminderDestinationID uuid.UUID) ([]models.ReminderError, error) {
	var reminderErrors []models.ReminderError
	err := r.db.Preload("Reminder").
		Preload("ReminderDestination").
		Where("reminder_destination_id = ? AND fixed = false", reminderDestinationID).
		Order("timestamp DESC").
		Find(&reminderErrors).Error
	return reminderErrors, err
}

// MarkAsFixed marks a reminder error as fixed
func (r *reminderErrorRepository) MarkAsFixed(id uuid.UUID) error {
	return r.db.Model(&models.ReminderError{}).Where("id = ?", id).Update("fixed", true).Error
}

// MarkMultipleAsFixed marks multiple reminder errors as fixed
func (r *reminderErrorRepository) MarkMultipleAsFixed(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&models.ReminderError{}).Where("id IN ?", ids).Update("fixed", true).Error
}
