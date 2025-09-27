package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// HandlePauseReminder handles the pause subcommand
func HandlePauseReminder(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	if len(options) == 0 {
		return utils.SendError(session, interaction, "Missing Parameter", "Please specify a reminder to pause.")
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

	// Check if user has permission to modify this reminder
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Permission Denied", "You don't have permission to modify this reminder.")
	}

	// Check if it's a one-time reminder
	recurrenceType := services.GetRecurrenceType(int(reminder.Recurrence))
	if recurrenceType == services.RecurrenceOnce {
		return utils.SendError(session, interaction, "Cannot Pause One-Time Reminder", 
			"One-time reminders cannot be paused. You can delete them instead if needed.")
	}

	// Check if already paused
	if services.IsPaused(int(reminder.Recurrence)) {
		return utils.SendInfo(session, interaction, "Already Paused", 
			fmt.Sprintf("The reminder \"%s\" is already paused.", reminder.Message))
	}

	// Update the recurrence to include the pause bit
	reminder.Recurrence = int16(services.SetPauseState(int(reminder.Recurrence), true))

	// Save the updated reminder
	err = repo.Reminder.Update(reminder, true)
	if err != nil {
		return utils.SendError(session, interaction, "Database Error", "Failed to pause the reminder. Please try again.")
	}

	// Send success message
	successEmbed := utils.BuildSuccessEmbed(session, "Reminder Paused! ⏸️", 
		fmt.Sprintf("The reminder \"%s\" has been successfully paused. You can unpause it anytime using `/reminders unpause`.", 
			reminder.Message), nil)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{successEmbed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// HandleUnpauseReminder handles the unpause subcommand
func HandleUnpauseReminder(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	if len(options) == 0 {
		return utils.SendError(session, interaction, "Missing Parameter", "Please specify a reminder to restart.")
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

	// Check if user has permission to modify this reminder
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Permission Denied", "You don't have permission to modify this reminder.")
	}

	// Check if it's a one-time reminder
	recurrenceType := services.GetRecurrenceType(int(reminder.Recurrence))
	if recurrenceType == services.RecurrenceOnce {
		return utils.SendError(session, interaction, "Cannot Restart One-Time Reminder", 
			"One-time reminders cannot be paused or restarted.")
	}

	// Check if not paused
	if !services.IsPaused(int(reminder.Recurrence)) {
		return utils.SendInfo(session, interaction, "Not Paused", 
			fmt.Sprintf("The reminder \"%s\" is not currently paused.", reminder.Message))
	}

	// Update the recurrence to remove the pause bit
	reminder.Recurrence = int16(services.SetPauseState(int(reminder.Recurrence), false))

	// Recalculate the next occurrence from now to avoid catching up
	nextTime, err := services.RecalculateNextOccurrence(reminder.RemindAtUTC, int(reminder.Recurrence))
	if err != nil {
		return utils.SendError(session, interaction, "Calculation Error", "Failed to recalculate the next reminder time.")
	}

	// Update the reminder time
	reminder.RemindAtUTC = nextTime

	// Save the updated reminder
	err = repo.Reminder.Update(reminder, true)
	if err != nil {
		return utils.SendError(session, interaction, "Database Error", "Failed to unpause the reminder. Please try again.")
	}

	// Send success message
	successEmbed := utils.BuildSuccessEmbed(session, "Reminder Resumed! ▶️", 
		fmt.Sprintf("The reminder \"%s\" has been successfully resumed and is now active again.", 
			reminder.Message), nil)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{successEmbed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}