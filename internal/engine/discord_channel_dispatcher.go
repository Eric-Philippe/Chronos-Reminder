package engine

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

type DiscordChannelDispatcher struct {
	session *discordgo.Session
}

func NewDiscordChannelDispatcher() *DiscordChannelDispatcher {
	return &DiscordChannelDispatcher{
		session: bot.GetDiscordSession(),
	}
}

// GetSupportedType returns the destination type this dispatcher supports
func (d *DiscordChannelDispatcher) GetSupportedType() models.DestinationType {
	return models.DestinationDiscordChannel
}

// Dispatch sends a reminder message to a Discord channel
func (d *DiscordChannelDispatcher) Dispatch(reminder *models.Reminder, destination *models.ReminderDestination) error {
	// Validate destination type
	if destination.Type != models.DestinationDiscordChannel {
		return fmt.Errorf("invalid destination type for Discord channel dispatcher: %s", destination.Type)
	}

	// Build the data from the metadata
	// {"guild_id": "912661874871533588", "channel_id": "913222458251837441"}
	channelID, exists := destination.Metadata["channel_id"]
	if !exists {
		return fmt.Errorf("channel_id not found in destination metadata")
	}

	channelIDStr, ok := channelID.(string)
	if !ok {
		return fmt.Errorf("channel_id is not a string: %v", channelID)
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

	// Send the message to the specified channel
	_, err := d.session.ChannelMessageSendEmbed(channelIDStr, embed)
	if err != nil {
		return fmt.Errorf("failed to send message to channel %s: %w", channelIDStr, err)
	}

	return nil
}