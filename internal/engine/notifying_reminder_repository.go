package engine

import (
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
)

// NotifyingReminderRepository wraps a reminder repository and notifies the scheduler of changes
type NotifyingReminderRepository struct {
	repo      repositories.ReminderRepository
	scheduler *Scheduler
}

// NewNotifyingReminderRepository creates a new notifying repository wrapper
func NewNotifyingReminderRepository(repo repositories.ReminderRepository, scheduler *Scheduler) *NotifyingReminderRepository {
	return &NotifyingReminderRepository{
		repo:      repo,
		scheduler: scheduler,
	}
}

// Create creates a reminder and notifies the scheduler
func (n *NotifyingReminderRepository) Create(reminder *models.Reminder) error {
	err := n.repo.Create(reminder)
	if err == nil {
		n.scheduler.NotifyReminderCreated()
	}
	return err
}

// Update updates a reminder and notifies the scheduler
func (n *NotifyingReminderRepository) Update(reminder *models.Reminder) error {
	err := n.repo.Update(reminder)
	if err == nil {
		n.scheduler.NotifyReminderUpdated()
	}
	return err
}

// Delete deletes a reminder and notifies the scheduler
func (n *NotifyingReminderRepository) Delete(id uuid.UUID, notify bool) error {
	err := n.repo.Delete(id, false) // Pass false to avoid double notification
	if err == nil {
		n.scheduler.NotifyReminderDeleted()
	}
	return err
}

// All other methods are pass-through without notifications
func (n *NotifyingReminderRepository) GetByID(id uuid.UUID) (*models.Reminder, error) {
	return n.repo.GetByID(id)
}

func (n *NotifyingReminderRepository) GetByAccountID(accountID uuid.UUID) ([]models.Reminder, error) {
	return n.repo.GetByAccountID(accountID)
}

func (n *NotifyingReminderRepository) GetByAccountIDWithDestinations(accountID uuid.UUID) ([]models.Reminder, error) {
	return n.repo.GetByAccountIDWithDestinations(accountID)
}

func (n *NotifyingReminderRepository) GetWithDestinations(id uuid.UUID) (*models.Reminder, error) {
	return n.repo.GetWithDestinations(id)
}

func (n *NotifyingReminderRepository) GetWithAccount(id uuid.UUID) (*models.Reminder, error) {
	return n.repo.GetWithAccount(id)
}

func (n *NotifyingReminderRepository) GetWithAccountAndDestinations(id uuid.UUID) (*models.Reminder, error) {
	return n.repo.GetWithAccountAndDestinations(id)
}

func (n *NotifyingReminderRepository) GetDueReminders(beforeTime time.Time) ([]models.Reminder, error) {
	return n.repo.GetDueReminders(beforeTime)
}

func (n *NotifyingReminderRepository) GetUpcomingReminders(accountID uuid.UUID, limit int) ([]models.Reminder, error) {
	return n.repo.GetUpcomingReminders(accountID, limit)
}

func (n *NotifyingReminderRepository) GetRemindersByDateRange(accountID uuid.UUID, startDate, endDate time.Time) ([]models.Reminder, error) {
	return n.repo.GetRemindersByDateRange(accountID, startDate, endDate)
}

func (n *NotifyingReminderRepository) GetNextReminders() ([]models.Reminder, error) {
	return n.repo.GetNextReminders()
}
