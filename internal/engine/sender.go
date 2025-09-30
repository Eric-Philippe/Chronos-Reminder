package engine

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// DiscordSend handles sending reminders via Discord
func DiscordSend(session *discordgo.Session, reminder *models.Reminder, channelID string, account *models.Account) error {
	// Create the reminder message
	embed := &discordgo.MessageEmbed{
		Title:       "⌛ | You have a new reminder ! ⌛",
		Color:       0xCEA04D,
	}

	// Send the message
	_, err := session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return fmt.Errorf("failed to send DM  %w", err)
	}

	// Convert the due date to the user's local timezone if available
	loc, err := time.LoadLocation(account.Timezone.IANALocation)
	if err == nil {
		reminder.RemindAtUTC = reminder.RemindAtUTC.In(loc)
	}
	

	img, err := services.NewDrawService("./assets").GenerateReminderImage(services.TextOverlay{
		Label: reminder.Message,
		Date:  reminder.RemindAtUTC,
	})

	// Check for errors
	if err != nil {
		return fmt.Errorf("failed to generate reminder image: %w", err)
	}

	// Encode img (image.Image) to PNG and wrap in io.Reader
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return fmt.Errorf("failed to encode reminder image: %w", err)
	}

		// Add a button to the message
	button := discordgo.Button{
		Label:   "Snooze",
		Style:   discordgo.SecondaryButton,
		CustomID: "reminder_snooze_" + fmt.Sprint(reminder.ID),
	}
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{button},
		},
	}
	msg := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:        "reminder.png",
			ContentType: "image/png",
			Reader:      &buf,
		},
		Components: components,
	}
	_, err = session.ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		return fmt.Errorf("failed to send reminder: %w", err)
	}

	return nil
}