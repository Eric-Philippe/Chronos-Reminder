package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (ReminderError) TableName() string {
	return "reminder_errors"
}

// ReminderError represents the reminder_errors table
type ReminderError struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReminderID            uuid.UUID `gorm:"type:uuid;not null;index" json:"reminder_id"`
	ReminderDestinationID uuid.UUID `gorm:"type:uuid;not null;index" json:"reminder_destination_id"`
	Timestamp             time.Time `gorm:"not null;default:now()" json:"timestamp"`
	Stacktrace            string    `gorm:"type:text;not null" json:"stacktrace"`
	Fixed                 bool      `gorm:"not null;default:false" json:"fixed"`
	
	// Relationships
	Reminder            *Reminder            `gorm:"foreignKey:ReminderID;constraint:OnDelete:CASCADE" json:"reminder,omitempty"`
	ReminderDestination *ReminderDestination `gorm:"foreignKey:ReminderDestinationID;constraint:OnDelete:CASCADE" json:"reminder_destination,omitempty"`
}

// BeforeCreate hooks for setting timestamps and UUIDs
func (re *ReminderError) BeforeCreate(tx *gorm.DB) error {
	if re.ID == uuid.Nil {
		re.ID = uuid.New()
	}
	if re.Timestamp.IsZero() {
		re.Timestamp = time.Now()
	}
	return nil
}
