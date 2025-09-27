package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
)

// handleDeleteReminder handles the delete subcommand
func HandleDeleteReminder(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account, options []*discordgo.ApplicationCommandInteractionDataOption) error {
	if len(options) == 0 {
		return utils.SendError(session, interaction, "Missing Parameter", "Please specify a reminder to delete.")
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

	// Check if user has permission to delete this reminder
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Permission Denied", "You don't have permission to delete this reminder.")
	}

	// Show confirmation message with buttons
	embed := utils.BuildWarningEmbed(session, "Confirm Deletion", 
		fmt.Sprintf("Are you sure you want to delete the reminder:\n\n**Message:** %s\n**Remind Time:** %s\n\nThis action cannot be undone.", 
			reminder.Message, 
			reminder.RemindAtUTC.Format("Monday, January 2, 2006 at 15:04 MST")))

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "confirm_delete_" + reminderIDStr,
							Label:    "Yes, Delete",
							Style:    discordgo.DangerButton,
							Emoji:    &discordgo.ComponentEmoji{Name: "🗑️"},
						},
						discordgo.Button{
							CustomID: "cancel_delete_" + reminderIDStr,
							Label:    "Cancel",
							Style:    discordgo.SecondaryButton,
							Emoji:    &discordgo.ComponentEmoji{Name: "❌"},
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral, // Only the user can see and interact with this message
		},
	})
}

// HandleConfirmDelete handles the confirmation button click
func HandleConfirmDelete(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Extract reminder ID from custom ID
	customID := interaction.MessageComponentData().CustomID
	reminderIDStr := strings.TrimPrefix(customID, "confirm_delete_")
	
	reminderID, err := uuid.Parse(reminderIDStr)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Reminder ID", "The reminder ID is not valid.")
	}

	repo := database.GetRepositories()

	// Get the reminder to verify it still exists and user has permission
	reminder, err := repo.Reminder.GetWithAccountAndDestinations(reminderID)
	if err != nil {
		return utils.SendError(session, interaction, "Database Error", "Failed to retrieve reminder information.")
	}

	if reminder == nil {
		return utils.SendError(session, interaction, "Reminder Not Found", "The reminder has already been deleted or does not exist.")
	}

	// Re-check permissions
	if !CanAccessReminder(interaction, account, reminder) {
		return utils.SendError(session, interaction, "Permission Denied", "You don't have permission to delete this reminder.")
	}

	// Delete the reminder (destinations will be cascade deleted)
	err = repo.Reminder.Delete(reminderID, true)
	if err != nil {
		return utils.SendError(session, interaction, "Database Error", "Failed to delete the reminder. Please try again.")
	}

	// Update the message to show success
	successEmbed := utils.BuildSuccessEmbed(session, "Reminder Deleted!", 
		fmt.Sprintf("The reminder \"%s\" has been successfully deleted.", reminder.Message), nil)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{successEmbed},
			Components: []discordgo.MessageComponent{}, // Remove buttons
		},
	})
}

// HandleCancelDelete handles the cancel button click
func HandleCancelDelete(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Update the message to show cancellation
	cancelEmbed := utils.BuildInfoEmbed(session, "Deletion Cancelled", "The reminder deletion has been cancelled.")

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{cancelEmbed},
			Components: []discordgo.MessageComponent{}, // Remove buttons
		},
	})
}