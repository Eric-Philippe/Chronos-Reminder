package repositories

import (
	"errors"
	"log"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SchedulerNotifier interface for notifying the scheduler of reminder changes
type SchedulerNotifier interface {
	NotifyReminderCreated()
	NotifyReminderUpdated()
	NotifyReminderDeleted()
}

// reminderRepository implementation
type reminderRepository struct {
	db        *gorm.DB
	scheduler SchedulerNotifier
}

// NewReminderRepository creates a new reminder repository instance
func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db: db}
}

// SetScheduler sets the scheduler notifier for the repository
func (r *reminderRepository) SetScheduler(scheduler SchedulerNotifier) {
	r.scheduler = scheduler
}

// Reminder Repository Implementation
func (r *reminderRepository) Create(reminder *models.Reminder) error {
	err := r.db.Create(reminder).Error
	if err == nil {
		if r.scheduler != nil {
			log.Printf("[REPOSITORY] - Notifying scheduler of new reminder")
			r.scheduler.NotifyReminderCreated()
		}
	}
	return err
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
	err := r.db.Save(reminder).Error
	if err == nil {
		if r.scheduler != nil {
			log.Printf("[REPOSITORY] - Notifying scheduler of reminder update")
			r.scheduler.NotifyReminderUpdated()
		}
	}
	return err
}

func (r *reminderRepository) Delete(id uuid.UUID, notify bool) error {
	err := r.db.Delete(&models.Reminder{}, "id = ?", id).Error
	if err == nil {
		if r.scheduler != nil && notify {
			r.scheduler.NotifyReminderDeleted()
		}
	}
	return err
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

func (r *reminderRepository) GetNextReminders() ([]models.Reminder, error) {
	// First, check if the table is empty
	var count int64
	err := r.db.Model(&models.Reminder{}).Count(&count).Error
	if err != nil {
		return nil, err
	}

	if count == 0 {
		// No reminders exist
		return []models.Reminder{}, nil
	}

	// Find the next reminder(s) to process in a single query
	// Priority: past due reminders first, then earliest future reminders
	var reminders []models.Reminder
	now := time.Now().UTC()
	
	err = r.db.Preload("Account").
		Preload("Account.Timezone").
		Preload("Destinations").
		Where("remind_at_utc <= ? OR remind_at_utc = (SELECT MIN(remind_at_utc) FROM reminders WHERE remind_at_utc > ?)", now, now).
		Order("remind_at_utc ASC").
		Find(&reminders).Error

	if err != nil {
		return nil, err
	}

	// If we found reminders, prioritize past due ones
	if len(reminders) > 0 {
		// Check if we have past due reminders
		for _, reminder := range reminders {
			if reminder.RemindAtUTC.Before(now) || reminder.RemindAtUTC.Equal(now) {
				// Return only the first past due reminder
				return []models.Reminder{reminder}, nil
			}
		}
		
		// No past due reminders, return all reminders with the earliest future time
		earliestTime := reminders[0].RemindAtUTC
		var futureReminders []models.Reminder
		for _, reminder := range reminders {
			if reminder.RemindAtUTC.Equal(earliestTime) {
				futureReminders = append(futureReminders, reminder)
			}
		}
		return futureReminders, nil
	}

	// No reminders found
	return []models.Reminder{}, nil
}
