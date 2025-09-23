package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Reminder) TableName() string {
	return "reminders"
}

// Reminder represents the reminder table
type Reminder struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID    uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	RemindAtUTC  time.Time `gorm:"not null" json:"remind_at_utc"`
	Message      string    `gorm:"not null" json:"message"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	Recurrence   int16     `gorm:"not null;default:0" json:"recurrence"`
	
	// Relationships
	Account      *Account               `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
	Destinations []ReminderDestination  `gorm:"foreignKey:ReminderID;constraint:OnDelete:CASCADE" json:"destinations,omitempty"`
}

// BeforeCreate hooks for setting timestamps and UUIDs
func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	r.CreatedAt = time.Now()
	return nil
}
