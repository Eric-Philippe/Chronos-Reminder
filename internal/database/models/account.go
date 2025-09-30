package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Account) TableName() string {
	return "accounts"
}

// Account represents the accounts table
type Account struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TimezoneID *uint     `gorm:"index" json:"timezone_id"`
	CreatedAt  time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null;default:now()" json:"updated_at"`
	
	// Relationships
	Timezone   *Timezone  `gorm:"foreignKey:TimezoneID" json:"timezone,omitempty"`
	Identities []Identity `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"identities,omitempty"`
	Reminders  []Reminder `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"reminders,omitempty"`
}

// BeforeCreate hooks for setting timestamps and UUIDs
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	return nil
}

func (a *Account) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}