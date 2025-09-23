package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Identity) TableName() string {
	return "identities"
}

// ProviderType represents the provider enum
type ProviderType string

// Providers enum
const (
	ProviderDiscord ProviderType = "discord"
	ProviderApp     ProviderType = "app"
)

// Value implements the driver.Valuer interface for database storage
func (p ProviderType) Value() (driver.Value, error) {
	return string(p), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (p *ProviderType) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	
	switch v := value.(type) {
	case string:
		*p = ProviderType(v)
	case []byte:
		*p = ProviderType(v)
	default:
		return fmt.Errorf("cannot scan %T into ProviderType", value)
	}
	
	// Validate the scanned value
	if *p != ProviderDiscord && *p != ProviderApp {
		return fmt.Errorf("invalid provider type: %s", *p)
	}
	
	return nil
}

// String returns the string representation of ProviderType
func (p ProviderType) String() string {
	return string(p)
}

// IsValid checks if the provider type is valid
func (p ProviderType) IsValid() bool {
	return p == ProviderDiscord || p == ProviderApp
}

// Identity represents the identities table
type Identity struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"account_id"`
	Provider     ProviderType `gorm:"type:provider_type;not null" json:"provider"` // enum type from database
	ExternalID   string    `gorm:"not null" json:"external_id"`  // e.g., discord_id or app email
	Username     *string   `json:"username"`                     // snapshot for display purposes
	Avatar       *string   `json:"avatar"`                       // optional, snapshot of Discord avatar
	PasswordHash *string   `json:"-"`                            // only for app identities, hidden in JSON
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`
	
	// Relationships
	Account *Account `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
}

func (i *Identity) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	i.CreatedAt = time.Now()
	return nil
}