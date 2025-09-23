package engine

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// DiscordDMDispatcher handles sending reminders via Discord DM
type DiscordDMDispatcher struct {
	session *discordgo.Session
}

// NewDiscordDMDispatcher creates a new Discord DM dispatcher
func NewDiscordDMDispatcher() *DiscordDMDispatcher {
	return &DiscordDMDispatcher{
		session: bot.GetDiscordSession(),
	}
}

// GetSupportedType returns the destination type this dispatcher supports
func (d *DiscordDMDispatcher) GetSupportedType() models.DestinationType {
	return models.DestinationDiscordDM
}

// Dispatch sends a reminder message via Discord DM
func (d *DiscordDMDispatcher) Dispatch(reminder *models.Reminder, destination *models.ReminderDestination) error {
	// Validate destination type
	if destination.Type != models.DestinationDiscordDM {
		return fmt.Errorf("invalid destination type for Discord DM dispatcher: %s", destination.Type)
	}

	// Extract user ID from metadata
	userID, exists := destination.Metadata["user_id"]
	if !exists {
		return fmt.Errorf("user_id not found in destination metadata")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return fmt.Errorf("user_id is not a string: %v", userID)
	}

	// Create DM channel with the user
	dmChannel, err := d.session.UserChannelCreate(userIDStr)
	if err != nil {
		return fmt.Errorf("failed to create DM channel with user %s: %w", userIDStr, err)
	}

	// Create the reminder message
	embed := &discordgo.MessageEmbed{
		Title:       "⏰ Reminder",
		Description: reminder.Message,
		Color:       0x00ff00, // Green color
		Timestamp:   reminder.RemindAtUTC.Format("2006-01-02T15:04:05Z"),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Chronos Reminder",
		},
	}

	// Send the message
	_, err = d.session.ChannelMessageSendEmbed(dmChannel.ID, embed)
	if err != nil {
		return fmt.Errorf("failed to send DM to user %s: %w", userIDStr, err)
	}

	return nil
}
