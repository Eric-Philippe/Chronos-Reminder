package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
)

// HandleShowReminder handles the show subcommand
func HandleShowReminder(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	if len(options) == 0 {
		return utils.SendError(session, interaction, "Missing Parameter", "Please specify a reminder to show.")
	}

	reminderIDStr := options[0].StringValue()
	reminderID, err := uuid.Parse(reminderIDStr)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Reminder ID", "The provided reminder ID is not valid.")
	}

	repo := database.GetRepositories()

	// Get the reminder with account and destinations
	reminder, err := repo.Reminder.GetWithAccountAndDestinations(reminderID)
	if err != nil {
		return utils.SendError(session, interaction, "Database Error", "Failed to retrieve reminder information.")
	}

	if reminder == nil {
		return utils.SendError(session, interaction, "Reminder Not Found", "The specified reminder does not exist.")
	}

	// Check if user has permission to access this reminder
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Permission Denied", "You don't have permission to access this reminder.")
	}

	// Build the reminder embed
	embed := BuildReminderEmbed(session, reminder)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
