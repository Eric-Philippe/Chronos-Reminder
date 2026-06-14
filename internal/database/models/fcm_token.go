package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (FcmToken) TableName() string {
	return "fcm_tokens"
}

// FcmToken represents a Firebase Cloud Messaging registration token for an
// Android device. A single account may have several tokens (one per device).
type FcmToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	Token     string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	DeviceID  string    `gorm:"type:text;not null" json:"device_id"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()" json:"updated_at"`

	// Relationships
	Account *Account `gorm:"foreignKey:AccountID;references:ID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
}

// BeforeCreate hooks for setting UUIDs and timestamps
func (t *FcmToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	return nil
}

// BeforeUpdate keeps UpdatedAt fresh
func (t *FcmToken) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
