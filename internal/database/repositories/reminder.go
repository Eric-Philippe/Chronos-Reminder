package repositories

import (
	"errors"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// reminderRepository implementation
type reminderRepository struct {
	db *gorm.DB
}

// NewReminderRepository creates a new reminder repository instance
func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

// Reminder Repository Implementation
func (r *reminderRepository) Create(reminder *models.Reminder) error {
	return r.db.Create(reminder).Error
}

func (r *reminderRepository) GetByID(id uuid.UUID) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.First(&reminder, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) GetByAccountID(accountID uuid.UUID) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("account_id = ?", accountID).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetByAccountIDWithDestinations(accountID uuid.UUID) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Preload("Destinations").Where("account_id = ?", accountID).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetWithDestinations(id uuid.UUID) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.Preload("Destinations").First(&reminder, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) GetWithAccount(id uuid.UUID) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.Preload("Account").Preload("Account.Timezone").First(&reminder, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) GetWithAccountAndDestinations(id uuid.UUID) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.Preload("Account").Preload("Account.Timezone").Preload("Destinations").First(&reminder, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reminder, nil
}

func (r *reminderRepository) Update(reminder *models.Reminder) error {
	return r.db.Save(reminder).Error
}

func (r *reminderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Reminder{}, "id = ?", id).Error
}

func (r *reminderRepository) GetDueReminders(beforeTime time.Time) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Preload("Account").Preload("Account.Timezone").Preload("Destinations").Where("remind_at_utc <= ?", beforeTime).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetUpcomingReminders(accountID uuid.UUID, limit int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("account_id = ? AND remind_at_utc > ?", accountID, time.Now().UTC()).
		Order("remind_at_utc ASC").
		Limit(limit).
		Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) GetRemindersByDateRange(accountID uuid.UUID, startDate, endDate time.Time) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("account_id = ? AND remind_at_utc BETWEEN ? AND ?", accountID, startDate, endDate).
		Order("remind_at_utc ASC").
		Find(&reminders).Error
	return reminders, err
}
