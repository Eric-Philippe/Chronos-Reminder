package repositories

import "gorm.io/gorm"

// Repositories contains all repository instances
type Repositories struct {
	Timezone            TimezoneRepository
	Account             AccountRepository
	Identity            IdentityRepository
	Reminder            ReminderRepository
	ReminderDestination ReminderDestinationRepository
}

// NewRepositories creates new repository instances
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Timezone:            NewTimezoneRepository(db),
		Account:             NewAccountRepository(db),
		Identity:            NewIdentityRepository(db),
		Reminder:            NewReminderRepository(db),
		ReminderDestination: NewReminderDestinationRepository(db),
	}
}
