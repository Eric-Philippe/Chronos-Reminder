package engine

import (
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// WebhookDispatcher handles sending reminders via webhooks (placeholder)
type WebhookDispatcher struct{}

// NewWebhookDispatcher creates a new webhook dispatcher
func NewWebhookDispatcher() *WebhookDispatcher {
	return &WebhookDispatcher{}
}

// GetSupportedType returns the destination type this dispatcher supports
func (d *WebhookDispatcher) GetSupportedType() models.DestinationType {
	return models.DestinationWebhook
}

// Dispatch is a placeholder for webhook dispatching
func (d *WebhookDispatcher) Dispatch(reminder *models.Reminder, destination *models.ReminderDestination) error {
	log.Printf("[WEBHOOK_DISPATCHER] - Placeholder: would send reminder %s via webhook", reminder.ID)
	// TODO: Implement webhook dispatching
	return nil
}

