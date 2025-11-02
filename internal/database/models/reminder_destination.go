package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ReminderDestination) TableName() string {
	return "reminder_destinations"
}

// DestinationType represents the destination_type enum
type DestinationType string

// Destination types enum
const (
	DestinationDiscordDM      DestinationType = "discord_dm"
	DestinationDiscordChannel DestinationType = "discord_channel"
	DestinationWebhook        DestinationType = "webhook"
)

// WebhookPlatform represents the optional platform for webhook destinations
type WebhookPlatform string

// Webhook platform types
const (
	WebhookPlatformGeneric WebhookPlatform = "generic"
	WebhookPlatformDiscord WebhookPlatform = "discord"
	WebhookPlatformSlack   WebhookPlatform = "slack"
)

// IsValid checks if the webhook platform is valid
func (w WebhookPlatform) IsValid() bool {
	return w == WebhookPlatformGeneric || w == WebhookPlatformDiscord || w == WebhookPlatformSlack
}

// String returns the string representation of WebhookPlatform
func (w WebhookPlatform) String() string {
	return string(w)
}

// Value implements the driver.Valuer interface for database storage
func (d DestinationType) Value() (driver.Value, error) {
	return string(d), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (d *DestinationType) Scan(value interface{}) error {
	if value == nil {
		*d = ""
		return nil
	}
	
	switch v := value.(type) {
	case string:
		*d = DestinationType(v)
	case []byte:
		*d = DestinationType(v)
	default:
		return fmt.Errorf("cannot scan %T into DestinationType", value)
	}
	
	// Validate the scanned value
	if !d.IsValid() {
		return fmt.Errorf("invalid destination type: %s", *d)
	}
	
	return nil
}

// String returns the string representation of DestinationType
func (d DestinationType) String() string {
	return string(d)
}

// IsValid checks if the destination type is valid
func (d DestinationType) IsValid() bool {
	return d == DestinationDiscordDM || d == DestinationDiscordChannel || d == DestinationWebhook
}

// JSONB is a custom type for JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}
	
	return json.Unmarshal(bytes, j)
}

// ReminderDestination represents the reminder_destinations table
type ReminderDestination struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReminderID uuid.UUID       `gorm:"type:uuid;not null;index" json:"reminder_id"`
	Type       DestinationType `gorm:"type:destination_type;not null" json:"type"`
	Metadata   JSONB           `gorm:"type:jsonb;not null" json:"metadata"`
	
	// Relationships
	Reminder *Reminder `gorm:"foreignKey:ReminderID;constraint:OnDelete:CASCADE" json:"reminder,omitempty"`
}

// BeforeCreate hooks for setting UUIDs
func (rd *ReminderDestination) BeforeCreate(tx *gorm.DB) error {
	if rd.ID == uuid.Nil {
		rd.ID = uuid.New()
	}
	return nil
}

// Validation methods for metadata based on destination type
func (rd *ReminderDestination) ValidateMetadata() error {
	switch rd.Type {
	case DestinationDiscordDM:
		if _, exists := rd.Metadata["user_id"]; !exists {
			return fmt.Errorf("discord_dm destination requires user_id in metadata")
		}
	case DestinationDiscordChannel:
		if _, exists := rd.Metadata["guild_id"]; !exists {
			return fmt.Errorf("discord_channel destination requires guild_id in metadata")
		}
		if _, exists := rd.Metadata["channel_id"]; !exists {
			return fmt.Errorf("discord_channel destination requires channel_id in metadata")
		}
	case DestinationWebhook:
		if _, exists := rd.Metadata["url"]; !exists {
			return fmt.Errorf("webhook destination requires url in metadata")
		}
		
		// Validate optional platform field
		if platformVal, exists := rd.Metadata["platform"]; exists {
			platformStr, ok := platformVal.(string)
			if !ok {
				return fmt.Errorf("webhook platform must be a string")
			}
			platform := WebhookPlatform(platformStr)
			if !platform.IsValid() {
				return fmt.Errorf("invalid webhook platform: %s (must be 'generic', 'discord', or 'slack')", platformStr)
			}
			
			// Platform-specific validation
			switch platform {
			case WebhookPlatformDiscord:
				// Discord webhooks can optionally have username and avatar_url
				// No strict requirements beyond the URL
			case WebhookPlatformSlack:
				// Slack webhooks can optionally have channel override
				// No strict requirements beyond the URL
			case WebhookPlatformGeneric:
				// Generic webhooks have no additional requirements
			}
		}
	default:
		return fmt.Errorf("invalid destination type: %s", rd.Type)
	}
	return nil
}

// BeforeUpdate and BeforeCreate validation
func (rd *ReminderDestination) BeforeUpdate(tx *gorm.DB) error {
	return rd.ValidateMetadata()
}

func (rd *ReminderDestination) BeforeSave(tx *gorm.DB) error {
	return rd.ValidateMetadata()
}
