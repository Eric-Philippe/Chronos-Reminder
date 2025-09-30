package repositories

import (
	"time"

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

// IdentityRepository defines the interface for identity database operations
type IdentityRepository interface {
	Create(identity *models.Identity) error
	GetByID(id uuid.UUID) (*models.Identity, error)
	GetByProviderAndExternalID(provider models.ProviderType, externalID string) (*models.Identity, error)
	GetByAccountID(accountID uuid.UUID) ([]models.Identity, error)
	Update(identity *models.Identity) error
	Delete(id uuid.UUID) error
}

// ReminderRepository interface defines operations for reminder data
type ReminderRepository interface {
	Create(reminder *models.Reminder, notify bool) error
	GetByID(id uuid.UUID) (*models.Reminder, error)
	GetByAccountID(accountID uuid.UUID) ([]models.Reminder, error)
	GetByAccountIDWithDestinations(accountID uuid.UUID) ([]models.Reminder, error)
	GetWithDestinations(id uuid.UUID) (*models.Reminder, error)
	GetWithAccount(id uuid.UUID) (*models.Reminder, error)
	GetWithAccountAndDestinations(id uuid.UUID) (*models.Reminder, error)
	Update(reminder *models.Reminder, notify bool) error
	Delete(id uuid.UUID, notify bool) error
	GetNextReminders() ([]models.Reminder, error)
	GetNextsRemindersToDelete() ([]models.Reminder, error)
	Reschedule(id uuid.UUID, newTime time.Time, notify bool) error
	RescheduleReminder(reminder *models.Reminder, newTime time.Time, notify bool) error
	Snooze(id uuid.UUID, snoozeUntil time.Time) error
	SnoozeReminder(reminder *models.Reminder, snoozeUntil time.Time) error
}

// ReminderDestinationRepository interface defines operations for reminder destination data
type ReminderDestinationRepository interface {
	Create(destination *models.ReminderDestination) error
	GetByID(id uuid.UUID) (*models.ReminderDestination, error)
	GetByReminderID(reminderID uuid.UUID) ([]models.ReminderDestination, error)
	GetByReminderIDWithReminder(reminderID uuid.UUID) ([]models.ReminderDestination, error)
	GetByType(destinationType models.DestinationType) ([]models.ReminderDestination, error)
	Update(destination *models.ReminderDestination) error
	Delete(id uuid.UUID) error
	DeleteByReminderID(reminderID uuid.UUID) error
	CreateMultiple(destinations []models.ReminderDestination) error
	GetByMetadataField(field string, value interface{}) ([]models.ReminderDestination, error)
}

// ReminderErrorRepository interface defines operations for reminder error data
type ReminderErrorRepository interface {
	Create(reminderError *models.ReminderError) error
	GetByID(id uuid.UUID) (*models.ReminderError, error)
	GetByReminderID(reminderID uuid.UUID) ([]models.ReminderError, error)
	GetByReminderDestinationID(reminderDestinationID uuid.UUID) ([]models.ReminderError, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.ReminderError, error)
	Delete(id uuid.UUID) error
	GetUnfixedByReminderID(reminderID uuid.UUID) ([]models.ReminderError, error)
	GetUnfixedByReminderDestinationID(reminderDestinationID uuid.UUID) ([]models.ReminderError, error)
	MarkAsFixed(id uuid.UUID) error
	MarkMultipleAsFixed(ids []uuid.UUID) error
}
