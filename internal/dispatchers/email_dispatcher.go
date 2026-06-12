package dispatchers

import (
	"fmt"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// EmailDispatcher handles sending reminders via email
type EmailDispatcher struct {
	mailer *services.MailerService
}

// NewEmailDispatcher creates a new email dispatcher
func NewEmailDispatcher(mailer *services.MailerService) *EmailDispatcher {
	return &EmailDispatcher{mailer: mailer}
}

// GetSupportedType returns the destination type this dispatcher supports
func (d *EmailDispatcher) GetSupportedType() models.DestinationType {
	return models.DestinationEmail
}

// Dispatch sends a reminder notification via email
func (d *EmailDispatcher) Dispatch(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) error {
	if destination.Type != models.DestinationEmail {
		return fmt.Errorf("invalid destination type for email dispatcher: %s", destination.Type)
	}

	emailVal, exists := destination.Metadata["email"]
	if !exists {
		return fmt.Errorf("email not found in destination metadata")
	}

	email, ok := emailVal.(string)
	if !ok || email == "" {
		return fmt.Errorf("email in destination metadata is not a valid string")
	}

	_, err := d.mailer.SendReminderNotificationEmail(email, reminder.Message, reminder.RemindAtUTC.String())
	return err
}
