package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// Description represents the description object for a command
type Description struct {
	Name             string `json:"name"`
	Emoji            string `json:"emoji"`
	CategoryName     string `json:"categoryName"`
	ShortDescription string `json:"shortDescription"`
	FullDescription  string `json:"fullDescription"`
	Usage            string `json:"usage"`
	Example          string `json:"example"`
}

// AutocompleteFunc represents the autocomplete function signature
type AutocompleteFunc func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error

// RunFunc represents the run function signature
type RunFunc func(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error

// Command represents a bot command with its description, data, and handlers
type Command struct {
	Description  Description                           `json:"description"`
	Data         *discordgo.ApplicationCommand         `json:"data"`
	Autocomplete *AutocompleteFunc                     `json:"autocomplete,omitempty"`
	NeedsAccount bool                                  `json:"needsAccount"`
	Run          RunFunc                               `json:"run"`
}

// commandRegistry holds registered commands by name
var commandRegistry = make(map[string]*Command)
var commands []*Command

// RegisterCommand registers a single command
func RegisterCommand(command *Command) {
	commands = append(commands, command)
	commandRegistry[command.Data.Name] = command
}

// RegisterCommands registers all commands with Discord and builds the command registry
func RegisterCommands(session *discordgo.Session) (int, error) {
	var applicationCommands []*discordgo.ApplicationCommand
	
	for _, cmd := range commands {
		applicationCommands = append(applicationCommands, cmd.Data)
	}
	
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", applicationCommands)
	if err != nil {
		return 0, err
	}

	return len(applicationCommands), nil
}

// HandleCommand routes an interaction to the appropriate command handler
func HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return nil
	}
	
	commandName := interaction.ApplicationCommandData().Name
	command, exists := commandRegistry[commandName]
	if !exists || command.Run == nil {
		return nil
	}

	var account *models.Account
	var err error
	
	if command.NeedsAccount {
		var user *discordgo.User
		if interaction.Member != nil && interaction.Member.User != nil {
			user = interaction.Member.User
		} else if interaction.User != nil {
			user = interaction.User
		} else {
			return nil // No user information available
		}
		
		account, err = services.EnsureDiscordUser(user)
		if err != nil {
			return err
		}
	}

	return command.Run(session, interaction, account)
}

// HandleAutocomplete routes an autocomplete interaction to the appropriate command handler
func HandleAutocomplete(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	if interaction.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return nil
	}
	
	commandName := interaction.ApplicationCommandData().Name
	command, exists := commandRegistry[commandName]
	if !exists || command.Autocomplete == nil {
		return nil
	}

	return (*command.Autocomplete)(session, interaction)
}

// dateAutocompleteHandler handles autocomplete for the date field
func DateAutocompleteHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	data := interaction.ApplicationCommandData()
	var currentInput string

	// Find the date option that's being typed
	for _, option := range data.Options {
		if option.Name == "date" && option.Focused {
			currentInput = strings.ToLower(strings.TrimSpace(option.StringValue()))
			break
		}
	}

	// Predefined suggestions
	suggestions := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Today", Value: "today"},
		{Name: "Tomorrow", Value: "tomorrow"},
		{Name: "Next Week", Value: "next week"},
		{Name: "Next Month", Value: "next month"},
	}

	// Filter suggestions based on current input
	var filteredSuggestions []*discordgo.ApplicationCommandOptionChoice
	for _, suggestion := range suggestions {
		if currentInput == "" || strings.Contains(strings.ToLower(suggestion.Name), currentInput) {
			filteredSuggestions = append(filteredSuggestions, suggestion)
		}
	}

	// Limit to 25 suggestions (Discord's limit)
	if len(filteredSuggestions) > 25 {
		filteredSuggestions = filteredSuggestions[:25]
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: filteredSuggestions,
		},
	})
}

// RemindersAutocompleteHandler handles autocomplete for the reminder selection
func RemindersAutocompleteHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	data := interaction.ApplicationCommandData()
	var currentInput string
	var subcommandName string

	// Find the focused option within the subcommand
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

	// Only provide autocomplete for delete, show, pause, and unpause subcommands
	if subcommandName != "delete" && subcommandName != "show" && subcommandName != "pause" && subcommandName != "unpause" {
		return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
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

	// Get user account from database
	repo := database.GetRepositories()
	identity, err := repo.Identity.GetByProviderAndExternalID(models.ProviderDiscord, user.ID)
	if err != nil || identity == nil {
		return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
	}

	account, err := repo.Account.GetByID(identity.AccountID)
	if err != nil || account == nil {
		return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		})
	}

	// Get reminders based on permissions
	var allReminders []models.Reminder
	var err2 error

	// Get user's own reminders
	userReminders, err2 := repo.Reminder.GetByAccountIDWithDestinations(account.ID)
	if err2 == nil {
		allReminders = append(allReminders, userReminders...)
	}

	// If in a server and user is admin, get server reminders
	if interaction.GuildID != "" && interaction.Member != nil {
		permissions := interaction.Member.Permissions
		isAdmin := (permissions & discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator
		
		if isAdmin {
			// Get all reminders that have a destination for this server
			serverDestinations, err3 := repo.ReminderDestination.GetByMetadataField("guild_id", interaction.GuildID)
			if err3 == nil {
				for _, dest := range serverDestinations {
					if dest.Type == models.DestinationDiscordChannel {
						serverReminder, err4 := repo.Reminder.GetWithAccountAndDestinations(dest.ReminderID)
						if err4 == nil && serverReminder != nil {
							// Check if not already in the list (avoid duplicates)
							found := false
							for _, existing := range allReminders {
								if existing.ID == serverReminder.ID {
									found = true
									break
								}
							}
							if !found {
								allReminders = append(allReminders, *serverReminder)
							}
						}
					}
				}
			}
		}
	}

	// Filter reminders based on input and create choices
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, reminder := range allReminders {
		// For pause/unpause commands, filter out one-time reminders
		if (subcommandName == "pause" || subcommandName == "unpause") {
			recurrenceType := services.GetRecurrenceType(int(reminder.Recurrence))
			if recurrenceType == services.RecurrenceOnce {
				continue // Skip one-time reminders
			}
		}

		// Additional filtering for pause/unpause based on current state
		if subcommandName == "pause" {
			// For pause: only show active (non-paused) reminders
			if services.IsPaused(int(reminder.Recurrence)) {
				continue // Skip already paused reminders
			}
		} else if subcommandName == "unpause" {
			// For unpause: only show paused reminders
			if !services.IsPaused(int(reminder.Recurrence)) {
				continue // Skip non-paused reminders
			}
		}

		// Create a display name that includes both message and time
		displayTime := reminder.RemindAtUTC.Format("Jan 2, 2006 15:04")
		displayName := fmt.Sprintf("[%s] %s", displayTime, reminder.Message)
		
		// Add status indicator for pause/unpause commands
		if subcommandName == "pause" || subcommandName == "unpause" {
			recurrenceName := services.GetRecurrenceTypeLabel(services.GetRecurrenceType(int(reminder.Recurrence)))
			if services.IsPaused(int(reminder.Recurrence)) {
				displayName += " [⏸️ Paused]"
			} else {
				displayName += fmt.Sprintf(" [🔁 %s]", recurrenceName)
			}
		}
		
		// Truncate if too long (Discord has a limit)
		if len(displayName) > 100 {
			displayName = displayName[:97] + "..."
		}

		// Filter based on input
		if currentInput == "" || strings.Contains(strings.ToLower(displayName), currentInput) || strings.Contains(strings.ToLower(reminder.Message), currentInput) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  displayName,
				Value: reminder.ID.String(),
			})
		}
	}

	// Limit to 25 choices (Discord's limit)
	if len(choices) > 25 {
		choices = choices[:25]
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}