package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (EmailVerification) TableName() string {
	return "email_verifications"
}

// EmailVerification represents the email_verifications table
type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID string    `gorm:"type:text;not null;index" json:"account_id"`
	Email     string    `gorm:"type:text;not null;index" json:"email"`
	Code      string    `gorm:"type:text;not null" json:"code"`
	Verified  bool      `gorm:"type:boolean;default:false" json:"verified"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	VerifiedAt *time.Time `gorm:"type:timestamp;default:null" json:"verified_at,omitempty"`

	// Relationships
	Account *Account `gorm:"foreignKey:AccountID;references:ID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
}

// BeforeCreate hooks for setting timestamps and UUIDs
func (e *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	e.CreatedAt = time.Now()
	return nil
}
