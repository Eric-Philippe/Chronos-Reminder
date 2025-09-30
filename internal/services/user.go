package services

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"gorm.io/gorm"
)

var (
	ErrUnknownTimezone = errors.New("unknown timezone ID")
)

// ChangeAccountTimezone changes the timezone of an account by timezone ID and saves it to the database
func ChangeAccountTimezone(account *models.Account, timezoneID uint) error {
	repo := database.GetRepositories()
	
	// Check that the given timezone ID exists
	timezone, err := repo.Timezone.GetByID(timezoneID)
	if err != nil {
		return err
	}
	if timezone == nil {
		return ErrUnknownTimezone
	}

	account.TimezoneID = &timezoneID

	// Save the changes to the database
	if err := database.GetDB().Save(account).Error; err != nil {
		return err
	}

	return nil
}

func GetAccountFromDiscordUser(discordUser *discordgo.User) (*models.Account, error) {
	var identity models.Identity
	err := database.GetDB().
		Preload("Account").
		Where("provider = ? AND external_id = ?", models.ProviderDiscord, discordUser.ID).
		First(&identity).Error
	
	if err != nil {
		return nil, err
	}
	
	return identity.Account, nil
}

// EnsureDiscordUser ensures a Discord user identity linked to an account exists, returns Account or error TODO: May be clever to store the result in a cache
func EnsureDiscordUser(discordUser *discordgo.User) (*models.Account, error) {
	var identity models.Identity
	err := database.GetDB().
		Preload("Account").
		// Preload the user's reminders
		Preload("Account.Reminders").
		// Preload the timezone for the account
		Preload("Account.Timezone").
		// Preload the identities for the account
		Where("provider = ? AND external_id = ?", models.ProviderDiscord, discordUser.ID).
		First(&identity).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Identity not found, create a new account and identity
			return createFromDiscordUser(discordUser)
		}
		// Some other error occurred
		return nil, err
	}
	
	// Identity found with preloaded account
	return identity.Account, nil
}

// CreateFromDiscordUser creates a Discord user identity linked to an account, returns Account or error
func createFromDiscordUser(discordUser *discordgo.User) (*models.Account, error) {
	// Get the cached default timezone ID from config
	defaultTimezoneID := config.GetDefaultTimezoneID()

	var account *models.Account
	
	// Use a database transaction to ensure both operations succeed or both fail
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		account = &models.Account{
			TimezoneID: defaultTimezoneID,
		}

		// Create the account within the transaction
		if err := tx.Create(account).Error; err != nil {
			return err
		}

		identity := &models.Identity{
			AccountID:  account.ID,
			Provider:   models.ProviderDiscord,
			ExternalID: discordUser.ID,
			Username:   &discordUser.Username,
			Avatar:     &discordUser.Avatar,
		}

		// Create the identity within the transaction
		if err := tx.Create(identity).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func DiscordUserUsesApp(account *models.Account) bool {
	for _, identity := range account.Identities {
		if identity.Provider == models.ProviderApp {
			return true
		}
	}
	return false
}