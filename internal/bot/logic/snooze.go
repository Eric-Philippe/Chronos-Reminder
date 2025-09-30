package logic

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
)

type SnoozeOption struct {
	Label    string
	Duration time.Duration
}

// Map of snooze options
var SnoozeOptions = []SnoozeOption{
	{"5m", 5 * time.Minute},
	{"10m", 10 * time.Minute},
	{"30m", 30 * time.Minute},
	{"1h", 1 * time.Hour},
	{"6h", 6 * time.Hour},
	{"1d", 24 * time.Hour},
}

func HandleSnooze(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Extract reminder ID from custom_id (format: reminder_request_snooze_<reminder_id>)
	customID := interaction.MessageComponentData().CustomID
	parts := strings.Split(customID, "_")
	if len(parts) != 4 {
		return utils.SendError(session, interaction, "Error", "Invalid snooze button configuration.")
	}

	reminderID, err := uuid.Parse(parts[3])
	if err != nil {
		return utils.SendError(session, interaction, "Error", "Invalid reminder ID.")
	}

	// Get the reminder with destinations and account
	repo := database.GetRepositories()
	reminder, err := repo.Reminder.GetWithAccountAndDestinations(reminderID)
	if err != nil {
		return utils.SendError(session, interaction, "Error", "Failed to retrieve reminder.")
	}

	if reminder == nil {
		return utils.SendError(session, interaction, "Error", "Reminder not found. It may have been deleted.")
	}

	// Check if user can access this reminder
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Error", "You don't have permission to snooze this reminder.")
	}

	// Check if reminder can be snoozed
	canSnooze, reason := CanSnoozeReminder(reminder)
	if !canSnooze {
		return utils.SendError(session, interaction, "Error", reason)
	}

	// Create snooze duration selection menu
	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "⏰ Snooze Reminder",
					Description: fmt.Sprintf("**Message:** %s\n\nChoose how long to snooze this reminder:", reminder.Message),
					Color:       utils.ColorInfo,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: utils.ClockLogo,
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "5 minutes",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("reminder_snooze_duration_%s_5m", reminderID.String()),
						},
						discordgo.Button{
							Label:    "10 minutes",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("reminder_snooze_duration_%s_10m", reminderID.String()),
						},
						discordgo.Button{
							Label:    "30 minutes",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("reminder_snooze_duration_%s_30m", reminderID.String()),
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "1 hour",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("reminder_snooze_duration_%s_1h", reminderID.String()),
						},
						discordgo.Button{
							Label:    "6 hours",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("reminder_snooze_duration_%s_6h", reminderID.String()),
						},
						discordgo.Button{
							Label:    "1 day",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("reminder_snooze_duration_%s_1d", reminderID.String()),
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Cancel",
							Style:    discordgo.SecondaryButton,
							CustomID: "reminder_snooze_cancel",
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	return err
}

func HandleSnoozeDuration(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Extract reminder ID and duration from custom_id (format: reminder_snooze_duration_<reminder_id>_<duration>)
	customID := interaction.MessageComponentData().CustomID
	parts := strings.Split(customID, "_")
	if len(parts) != 5 {
		return utils.SendError(session, interaction, "Error", "Invalid snooze configuration.")
	}

	reminderID, err := uuid.Parse(parts[3])
	if err != nil {
		return utils.SendError(session, interaction, "Error", "Invalid reminder ID.")
	}

	durationStr := parts[4]

	// Parse duration
	var duration time.Duration
	switch durationStr {
	case "5m":
		duration = 5 * time.Minute
	case "10m":
		duration = 10 * time.Minute
	case "30m":
		duration = 30 * time.Minute
	case "1h":
		duration = 1 * time.Hour
	case "6h":
		duration = 6 * time.Hour
	case "1d":
		duration = 24 * time.Hour
	default:
		return utils.SendError(session, interaction, "Error", "Invalid snooze duration.")
	}

	// Get the reminder
	repo := database.GetRepositories()
	reminder, err := repo.Reminder.GetWithAccountAndDestinations(reminderID)
	if err != nil {
		return utils.SendError(session, interaction, "Error", "Failed to retrieve reminder.")
	}

	if reminder == nil {
		return utils.SendError(session, interaction, "Error", "Reminder not found. It may have been deleted.")
	}

	// Verify permissions again
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Error", "You don't have permission to snooze this reminder.")
	}

	// Check if can still snooze (in case state changed)
	canSnooze, reason := CanSnoozeReminder(reminder)
	if !canSnooze {
		return utils.SendError(session, interaction, "Error", reason)
	}

	// Format duration for display
	var durationLabel string
	switch durationStr {
	case "5m":
		durationLabel = "5 minutes"
	case "10m":
		durationLabel = "10 minutes"
	case "30m":
		durationLabel = "30 minutes"
	case "1h":
		durationLabel = "1 hour"
	case "6h":
		durationLabel = "6 hours"
	case "1d":
		durationLabel = "1 day"
	}

	// Update the reminder in the database
	snoozeUntil := time.Now().UTC().Add(duration)
	err = repo.Reminder.SnoozeReminder(reminder, snoozeUntil)
	if err != nil {
		return utils.SendError(session, interaction, "Error", "Failed to update reminder.")
	}

	// Respond with success message
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "✅ Reminder Snoozed",
					Description: fmt.Sprintf("**Message:** %s\n\nThis reminder has been snoozed for **%s**.\n\nYou'll be reminded again at <t:%d:F>.", reminder.Message, durationLabel, snoozeUntil.Unix()),
					Color:       utils.ColorSuccess,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: utils.ClockLogo,
					},
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
}

func HandleSnoozeCancel(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "❌ Snooze Cancelled",
					Description: "The snooze action has been cancelled.",
					Color:       utils.ColorWarning,
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
}