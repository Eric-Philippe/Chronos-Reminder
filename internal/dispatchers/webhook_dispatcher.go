package dispatchers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// WebhookDispatcher handles sending reminders via webhooks
type WebhookDispatcher struct {
	formatter  *services.WebhookFormatter
	httpClient *http.Client
}

// NewWebhookDispatcher creates a new webhook dispatcher
func NewWebhookDispatcher() *WebhookDispatcher {
	return &WebhookDispatcher{
		formatter: services.NewWebhookFormatter(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetSupportedType returns the destination type this dispatcher supports
func (d *WebhookDispatcher) GetSupportedType() models.DestinationType {
	return models.DestinationWebhook
}

// Dispatch sends the reminder via webhook
func (d *WebhookDispatcher) Dispatch(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) error {
	// Extract webhook URL from metadata
	urlVal, exists := destination.Metadata["url"]
	if !exists {
		return fmt.Errorf("webhook URL not found in metadata")
	}

	url, ok := urlVal.(string)
	if !ok {
		return fmt.Errorf("webhook URL is not a string")
	}

	// Format the payload based on the platform
	payload, err := d.formatter.FormatPayload(reminder, destination, account)
	if err != nil {
		return fmt.Errorf("failed to format webhook payload: %w", err)
	}

	// Validate the payload
	if err := d.formatter.ValidatePayload(payload); err != nil {
		return fmt.Errorf("invalid webhook payload: %w", err)
	}

	// Log the request (for debugging)
	platform := "generic"
	if platformVal, exists := destination.Metadata["platform"]; exists {
		if platformStr, ok := platformVal.(string); ok {
			platform = platformStr
		}
	}
	log.Printf("[WEBHOOK_DISPATCHER] Sending %s webhook for reminder %s to %s", platform, reminder.ID, maskURL(url))

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", d.formatter.GetContentType(destination))
	req.Header.Set("User-Agent", "Chronos-Reminder/1.0")

	// Add any custom headers from metadata
	if headersVal, exists := destination.Metadata["headers"]; exists {
		if headers, ok := headersVal.(map[string]interface{}); ok {
			for key, value := range headers {
				if valueStr, ok := value.(string); ok {
					req.Header.Set(key, valueStr)
				}
			}
		}
	}

	// Send the request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status code: %d", resp.StatusCode)
	}

	log.Printf("[WEBHOOK_DISPATCHER] Successfully sent webhook for reminder %s (status: %d)", reminder.ID, resp.StatusCode)
	return nil
}

// maskURL masks sensitive parts of a URL for logging
func maskURL(url string) string {
	// Simple masking - show only the domain
	if len(url) > 30 {
		return url[:30] + "..."
	}
	return url
}

