package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (DFMItem) TableName() string {
	return "dfm_items"
}

// DFMItem represents a single todo entry inside a "Don't Forget Me" note
type DFMItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	NoteID    uuid.UUID `gorm:"type:uuid;not null;index" json:"note_id"`
	Content   string    `gorm:"not null" json:"content"`
	Checked   bool      `gorm:"not null;default:false" json:"checked"`
	Position  int       `gorm:"not null;default:0" json:"position"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Relationships
	Note *DFMNote `gorm:"foreignKey:NoteID;constraint:OnDelete:CASCADE" json:"note,omitempty"`
}

// BeforeCreate hooks for setting timestamps and UUIDs
func (i *DFMItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	now := time.Now()
	i.CreatedAt = now
	i.UpdatedAt = now
	return nil
}

func (i *DFMItem) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now()
	return nil
}
