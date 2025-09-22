package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Identity) TableName() string {
	return "identities"
}

// Identity represents the identities table
type Identity struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID    uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	Provider     string    `gorm:"not null" json:"provider"`     // e.g., 'discord' or 'app'
	ExternalID   string    `gorm:"not null" json:"external_id"`  // e.g., discord_id or app email
	Username     *string   `json:"username"`                     // snapshot for display purposes
	Avatar       *string   `json:"avatar"`                       // optional, snapshot of Discord avatar
	PasswordHash *string   `json:"-"`                            // only for app identities, hidden in JSON
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	
	// Relationships
	Account *Account `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
}

// Providers enum
const (
	ProviderDiscord = "discord"
	ProviderApp     = "app"
)

func (i *Identity) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	i.CreatedAt = time.Now()
	return nil
}