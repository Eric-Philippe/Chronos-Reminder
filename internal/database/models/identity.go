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

// Providers enum.
// NOTE: the legacy 'app' provider has been removed — email/password now live on
// the accounts table. The 'app' value may still exist in the DB enum but is
// never written or read by the application.
const (
	ProviderDiscord ProviderType = "discord"
	ProviderAPIKey  ProviderType = "api_key"
	ProviderMobile  ProviderType = "mobile"
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
	if *p != ProviderDiscord && *p != ProviderAPIKey && *p != ProviderMobile {
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
	return p == ProviderDiscord || p == ProviderAPIKey || p == ProviderMobile
}

// Identity represents the identities table
type Identity struct {
	ID           uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"account_id"`
	Provider     ProviderType `gorm:"type:provider_type;not null" json:"provider"` // enum type from database
	ExternalID   string    `gorm:"not null" json:"external_id"`  // discord_id, account-id (mobile), or api-key id
	Username     *string   `json:"username"`                     // snapshot for display purposes
	Avatar       *string   `json:"avatar"`                       // optional, snapshot of Discord avatar
	AccessToken  *string   `json:"-"`                            // Discord OAuth access token or API key hash, hidden in JSON
	RefreshToken *string   `json:"-"`                            // Discord OAuth refresh token, hidden in JSON
	Scopes       *string   `gorm:"type:text" json:"scopes,omitempty"`        // comma-separated scopes for API keys (e.g., "reminders.read")
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