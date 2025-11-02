package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// WebhookPayload represents a generic webhook payload
type WebhookPayload struct {
	Content string                 `json:"content,omitempty"`
	Data    map[string]interface{} `json:"-"` // Platform-specific data
}

// WebhookFormatter handles formatting reminder messages for different webhook platforms
type WebhookFormatter struct{}

// NewWebhookFormatter creates a new webhook formatter
func NewWebhookFormatter() *WebhookFormatter {
	return &WebhookFormatter{}
}

// FormatPayload formats a reminder message for the specified webhook platform
func (f *WebhookFormatter) FormatPayload(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) ([]byte, error) {
	// Get platform from metadata, default to generic
	platform := models.WebhookPlatformGeneric
	if platformVal, exists := destination.Metadata["platform"]; exists {
		if platformStr, ok := platformVal.(string); ok {
			platform = models.WebhookPlatform(platformStr)
		}
	}

	switch platform {
	case models.WebhookPlatformDiscord:
		return f.formatDiscordWebhook(reminder, destination, account)
	case models.WebhookPlatformSlack:
		return f.formatSlackWebhook(reminder, destination, account)
	case models.WebhookPlatformGeneric:
		return f.formatGenericWebhook(reminder, destination, account)
	default:
		return f.formatGenericWebhook(reminder, destination, account)
	}
}

// formatDiscordWebhook formats a reminder for Discord webhook
func (f *WebhookFormatter) formatDiscordWebhook(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) ([]byte, error) {
	// Build Discord embed
	embed := map[string]interface{}{
		"title":       "⏰ Reminder",
		"description": reminder.Message,
		"color":       3447003, // Blue color
		"timestamp":   time.Now().Format(time.RFC3339),
		"footer": map[string]interface{}{
			"text": "Chronos Reminder",
		},
	}

	// Add fields for additional information
	fields := []map[string]interface{}{
		{
			"name":   "Scheduled Time",
			"value":   reminder.RemindAtUTC.Format("Monday, January 2, 2006 at 15:04 MST"),
			"inline": false,
		},
	}

	if reminder.SnoozedAtUTC != nil {
		fields = append(fields, map[string]interface{}{
			"name":   "Snoozed Until",
			"value":   reminder.SnoozedAtUTC.Format("Monday, January 2, 2006 at 15:04 MST"),
			"inline": false,
		})
	}

	embed["fields"] = fields

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	// Add optional username override
	if username, exists := destination.Metadata["username"]; exists {
		if usernameStr, ok := username.(string); ok && usernameStr != "" {
			payload["username"] = usernameStr
		}
	}

	// Add optional avatar URL override
	if avatarURL, exists := destination.Metadata["avatar_url"]; exists {
		if avatarURLStr, ok := avatarURL.(string); ok && avatarURLStr != "" {
			payload["avatar_url"] = avatarURLStr
		}
	}

	return json.Marshal(payload)
}

// formatSlackWebhook formats a reminder for Slack webhook
func (f *WebhookFormatter) formatSlackWebhook(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) ([]byte, error) {
	// Build Slack block kit message
	blocks := []map[string]interface{}{
		{
			"type": "header",
			"text": map[string]interface{}{
				"type": "plain_text",
				"text": "⏰ Reminder",
			},
		},
		{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Message:*\n%s", reminder.Message),
			},
		},
		{
			"type": "section",
			"fields": []map[string]interface{}{
				{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*Scheduled Time:*\n%s", reminder.RemindAtUTC.Format("Monday, January 2, 2006 at 15:04 MST")),
				},
			},
		},
	}

	// Add snoozed information if applicable
	if reminder.SnoozedAtUTC != nil {
		blocks = append(blocks, map[string]interface{}{
			"type": "section",
			"text": map[string]interface{}{
				"type": "mrkdwn",
				"text": fmt.Sprintf("*Snoozed Until:*\n%s", reminder.SnoozedAtUTC.Format("Monday, January 2, 2006 at 15:04 MST")),
			},
		})
	}

	payload := map[string]interface{}{
		"blocks": blocks,
	}

	// Add optional channel override
	if channel, exists := destination.Metadata["channel"]; exists {
		if channelStr, ok := channel.(string); ok && channelStr != "" {
			payload["channel"] = channelStr
		}
	}

	// Add optional username override
	if username, exists := destination.Metadata["username"]; exists {
		if usernameStr, ok := username.(string); ok && usernameStr != "" {
			payload["username"] = usernameStr
		}
	}

	// Add optional icon emoji
	if iconEmoji, exists := destination.Metadata["icon_emoji"]; exists {
		if iconEmojiStr, ok := iconEmoji.(string); ok && iconEmojiStr != "" {
			payload["icon_emoji"] = iconEmojiStr
		}
	}

	return json.Marshal(payload)
}

// formatGenericWebhook formats a reminder for generic webhook
func (f *WebhookFormatter) formatGenericWebhook(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) ([]byte, error) {
	// Simple JSON payload with all reminder information
	payload := map[string]interface{}{
		"id":           reminder.ID.String(),
		"message":      reminder.Message,
		"remind_at":    reminder.RemindAtUTC.Format(time.RFC3339),
		"created_at":   reminder.CreatedAt.Format(time.RFC3339),
		"recurrence":   reminder.Recurrence,
	}

	if reminder.SnoozedAtUTC != nil {
		payload["snoozed_at"] = reminder.SnoozedAtUTC.Format(time.RFC3339)
	}

	if reminder.NextFireUTC != nil {
		payload["next_fire"] = reminder.NextFireUTC.Format(time.RFC3339)
	}

	if account != nil {
		payload["account_id"] = account.ID.String()
	}

	// Add any custom fields from metadata
	if customFields, exists := destination.Metadata["custom_fields"]; exists {
		if fields, ok := customFields.(map[string]interface{}); ok {
			for key, value := range fields {
				payload[key] = value
			}
		}
	}

	return json.Marshal(payload)
}

// GetContentType returns the appropriate Content-Type header for the webhook platform
func (f *WebhookFormatter) GetContentType(destination *models.ReminderDestination) string {
	// All platforms use JSON
	return "application/json"
}

// ValidatePayload validates that a payload can be sent successfully
func (f *WebhookFormatter) ValidatePayload(payload []byte) error {
	// Basic validation: ensure it's valid JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(payload, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON payload: %w", err)
	}

	// Check if payload is not empty
	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}

	return nil
}

// PrettyPrint returns a human-readable version of the payload for debugging
func (f *WebhookFormatter) PrettyPrint(payload []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, payload, "", "  "); err != nil {
		return string(payload)
	}
	return prettyJSON.String()
}
