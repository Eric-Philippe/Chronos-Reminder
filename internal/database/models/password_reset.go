package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (PasswordReset) TableName() string {
	return "password_resets"
}

// PasswordReset represents the password_resets table
type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	Email     string    `gorm:"type:text;not null;index" json:"email"`
	Token     string    `gorm:"type:text;not null;unique;index" json:"token"`
	Used      bool      `gorm:"type:boolean;default:false" json:"used"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UsedAt    *time.Time `gorm:"type:timestamp;default:null" json:"used_at,omitempty"`

	// Relationships
	Account *Account `gorm:"foreignKey:AccountID;references:ID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
}

// BeforeCreate hooks for setting timestamps and UUIDs
func (p *PasswordReset) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = time.Now()
	return nil
}
