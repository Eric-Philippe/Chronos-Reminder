package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/bot/handlers"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// remindersHandler handles the main reminders command
func remindersHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options
	if len(options) == 0 {
		return utils.SendError(session, interaction, "Invalid Command", "Please specify a subcommand.")
	}

	subcommand := options[0]
	switch subcommand.Name {
	case "list":
		return handlers.HandleListReminders(session, interaction, account, subcommand.Options)
	case "show":
		return handlers.HandleShowReminder(session, interaction, account, subcommand.Options)
	case "pause":
		return handlers.HandlePauseReminder(session, interaction, account, subcommand.Options)
	case "unpause":
		return handlers.HandleUnpauseReminder(session, interaction, account, subcommand.Options)
	case "delete":
		return handlers.HandleDeleteReminder(session, interaction, account, subcommand.Options)
	default:
		return utils.SendError(session, interaction, "Unknown Subcommand", "The specified subcommand is not recognized.")
	}
}

// Register the reminders command
func init() {
	autocompleteFunc := AutocompleteFunc(handlers.RemindersAutocompleteHandler)
	
	RegisterCommand(&Command{
		Description: Description{
			Name:             "reminders",
			Emoji:            "📝",
			CategoryName:     "Reminders",
			ShortDescription: "Manage your reminders",
			FullDescription:  "List, show, pause, restart, or delete your existing reminders",
			Usage:            "/reminders <subcommand> [options]",
			Example:          "/reminders delete reminder:<reminder>",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "reminders",
			Description: "Manage your reminders",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "List all your reminders",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "show",
					Description: "Show details of a specific reminder",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "reminder",
							Description:  "The reminder to show",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "pause",
					Description: "Pause a reminder",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "reminder",
							Description:  "The reminder to pause",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "restart",
					Description: "Restart a reminder",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "reminder",
							Description:  "The reminder to restart",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "delete",
					Description: "Delete a reminder",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "reminder",
							Description:  "The reminder to delete",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
			},
		},
		NeedsAccount: true,
		Run:          remindersHandler,
		Autocomplete: &autocompleteFunc,
		MessageComponentHandlers: []MessageComponentHandler{
			{
				CustomID: "confirm_delete_",
				Handler:  handlers.HandleConfirmDelete,
			},
			{
				CustomID: "cancel_delete_",
				Handler:  handlers.HandleCancelDelete,
			},
			{
				CustomID: "show_reminder_",
				Handler:  handlers.HandleShowReminderFromList,
			},
			{
				CustomID: "back_to_list",
				Handler:  handlers.HandleBackToList,
			},
		},
	})
}