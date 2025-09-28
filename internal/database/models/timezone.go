package models

// TableName methods to ensure proper table naming
func (Timezone) TableName() string {
	return "timezones"
}

// Timezone represents the timezone table
type Timezone struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"not null;size:100" json:"name"`
	GMTOffset    float64 `gorm:"not null;type:decimal(4,2)" json:"gmt_offset"`
	IANALocation string  `gorm:"not null;size:50" json:"iana_location"` // IANA timezone identifier
	
	// Relationships
	Accounts []Account `gorm:"foreignKey:TimezoneID" json:"accounts,omitempty"`
}