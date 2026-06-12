package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/bot/logic"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// dfmHandler handles the main dfm command
func dfmHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options
	if len(options) == 0 {
		return utils.SendError(session, interaction, "Invalid Command", "Please specify a subcommand.")
	}

	subcommand := options[0]
	switch subcommand.Name {
	case "create":
		return logic.HandleDFMCreate(session, interaction, account, subcommand.Options)
	case "list":
		return logic.HandleDFMList(session, interaction, account)
	case "check":
		return logic.HandleDFMSetChecked(session, interaction, account, subcommand.Options, true)
	case "uncheck":
		return logic.HandleDFMSetChecked(session, interaction, account, subcommand.Options, false)
	case "delete":
		return logic.HandleDFMDelete(session, interaction, account, subcommand.Options)
	case "set-reminder":
		return logic.HandleDFMSetReminder(session, interaction, account, subcommand.Options)
	case "remove-reminder":
		return logic.HandleDFMRemoveReminder(session, interaction, account)
	case "send":
		return logic.HandleDFMSend(session, interaction, account)
	default:
		return utils.SendError(session, interaction, "Unknown Subcommand", "The specified subcommand is not recognized.")
	}
}

// DFMAutocompleteHandler provides item suggestions for the check, uncheck and delete subcommands
func DFMAutocompleteHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	data := interaction.ApplicationCommandData()

	var currentInput string
	var subcommandName string
	for _, option := range data.Options {
		if option.Type == discordgo.ApplicationCommandOptionSubCommand {
			subcommandName = option.Name
			for _, subOption := range option.Options {
				if subOption.Focused {
					currentInput = strings.ToLower(strings.TrimSpace(subOption.StringValue()))
					break
				}
			}
			break
		}
	}

	respondChoices := func(choices []*discordgo.ApplicationCommandOptionChoice) error {
		return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
	}

	if subcommandName != "check" && subcommandName != "uncheck" && subcommandName != "delete" {
		return respondChoices([]*discordgo.ApplicationCommandOptionChoice{})
	}

	// Get user account
	var user *discordgo.User
	if interaction.Member != nil && interaction.Member.User != nil {
		user = interaction.Member.User
	} else if interaction.User != nil {
		user = interaction.User
	} else {
		return nil
	}

	repo := database.GetRepositories()
	identity, err := repo.Identity.GetByProviderAndExternalID(models.ProviderDiscord, user.ID)
	if err != nil || identity == nil {
		return respondChoices([]*discordgo.ApplicationCommandOptionChoice{})
	}

	note, err := repo.DFMNote.GetByAccountID(identity.AccountID)
	if err != nil || note == nil {
		return respondChoices([]*discordgo.ApplicationCommandOptionChoice{})
	}

	items, err := repo.DFMItem.GetByNoteID(note.ID)
	if err != nil {
		return respondChoices([]*discordgo.ApplicationCommandOptionChoice{})
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, item := range items {
		// Only suggest items the subcommand can act on
		if subcommandName == "check" && item.Checked {
			continue
		}
		if subcommandName == "uncheck" && !item.Checked {
			continue
		}

		displayName := item.Content
		if item.Checked {
			displayName = "[x] " + displayName
		} else {
			displayName = "[ ] " + displayName
		}
		if len(displayName) > 100 {
			displayName = displayName[:97] + "..."
		}

		if currentInput == "" || strings.Contains(strings.ToLower(item.Content), currentInput) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  displayName,
				Value: item.ID.String(),
			})
		}
	}

	if len(choices) > 25 {
		choices = choices[:25]
	}

	return respondChoices(choices)
}

// Register the dfm command
func init() {
	autocompleteFunc := AutocompleteFunc(DFMAutocompleteHandler)

	recurrenceChoices := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Daily", Value: services.GetRecurrenceTypeName(services.RecurrenceDaily)},
		{Name: "Weekly", Value: services.GetRecurrenceTypeName(services.RecurrenceWeekly)},
		{Name: "Monthly", Value: services.GetRecurrenceTypeName(services.RecurrenceMonthly)},
		{Name: "Yearly", Value: services.GetRecurrenceTypeName(services.RecurrenceYearly)},
		{Name: "Workdays", Value: services.GetRecurrenceTypeName(services.RecurrenceWorkdays)},
		{Name: "Weekend", Value: services.GetRecurrenceTypeName(services.RecurrenceWeekend)},
	}

	RegisterCommand(&Command{
		Description: Description{
			Name:             "dfm",
			Emoji:            "💭",
			CategoryName:     "Don't Forget Me",
			ShortDescription: "Manage your Don't Forget Me note",
			FullDescription:  "Keep a private note of things you don't want to forget and get the whole note sent to you on a recurring reminder",
			Usage:            "/dfm <subcommand> [options]",
			Example:          "/dfm create content:Buy groceries",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "dfm",
			Description: "Manage your Don't Forget Me note",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "create",
					Description: "Add an item to your note",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "content",
							Description: "What you don't want to forget",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "Show your note",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "check",
					Description: "Check an item of your note",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "item",
							Description:  "The item to check",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "uncheck",
					Description: "Uncheck an item of your note",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "item",
							Description:  "The item to uncheck",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "delete",
					Description: "Delete an item from your note",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "item",
							Description:  "The item to delete",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-reminder",
					Description: "Set the recurring reminder of your note",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "recurrence",
							Description: "How often the note should be sent to you",
							Required:    true,
							Choices:     recurrenceChoices,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "time",
							Description: "Time of day in HH:MM format (default 09:00)",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "date",
							Description: "Start date in YYYY-MM-DD format (default today)",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "destination",
							Description: "Where the note should be sent (default Discord DM)",
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{Name: "Discord DM", Value: "discord_dm"},
								{Name: "Email", Value: "email"},
								{Name: "Discord DM and Email", Value: "both"},
							},
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove-reminder",
					Description: "Remove the reminder of your note",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "send",
					Description: "Send your note to you right now",
				},
			},
		},
		NeedsAccount: true,
		Run:          dfmHandler,
		Autocomplete: &autocompleteFunc,
	})
}
