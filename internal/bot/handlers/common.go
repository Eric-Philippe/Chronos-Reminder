package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

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

// BuildReminderEmbed creates a detailed embed for a reminder with its destinations
func BuildReminderEmbed(session *discordgo.Session, reminder *models.Reminder) *discordgo.MessageEmbed {
	// Format the reminder time
	remindTimeStr := reminder.RemindAtUTC.Format("Monday, January 2, 2006 at 15:04")

	// Determine status
	status := "✅ Active"
	if services.IsPaused(int(reminder.Recurrence)) {
		status = "⏸️ Paused"
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📝 Reminder Details",
		Description: fmt.Sprintf("**Message:** %s", reminder.Message),
		Color:       utils.ColorInfo,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: utils.ClockLogo,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "⏰ Remind Time",
				Value:  remindTimeStr,
				Inline: true,
			},
			{
				Name:   "📊 Status",
				Value:  status,
				Inline: true,
			},
			{
				Name:   "🆔 Reminder ID",
				Value:  reminder.ID.String(),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: session.State.User.AvatarURL(""),
			Text:    "Chronos Bot Reminder",
		},
	}

	// Add owner information if different from viewer
	if reminder.Account != nil {
		ownerInfo := fmt.Sprintf("Account ID: %s", reminder.Account.ID.String())
		// Try to get Discord identity for better display
		repo := database.GetRepositories()
		identities, err := repo.Identity.GetByAccountID(reminder.Account.ID)
		if err == nil && len(identities) > 0 {
			// Look for Discord identity
			for _, identity := range identities {
				if identity.Provider == models.ProviderDiscord {
					ownerInfo = fmt.Sprintf("<@%s>", identity.ExternalID)
					break
				}
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "👤 Created By",
			Value:  ownerInfo,
			Inline: true,
		})
	}

	// If the reminder has a recurrence different from one-time, show it
	recurrenceType := services.GetRecurrenceType(int(reminder.Recurrence))
	if recurrenceType != services.RecurrenceOnce {
		recurrenceStr := services.GetRecurrenceTypeLabel(recurrenceType)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🔁 Recurrence",
			Value:  recurrenceStr,
			Inline: true,
		})
	}

	// Add destinations
	if len(reminder.Destinations) > 0 {
		for i, dest := range reminder.Destinations {
			destField := buildDestinationField(dest, i+1)
			embed.Fields = append(embed.Fields, destField)
		}
	}

	return embed
}

// buildDestinationField creates an embed field for a destination
func buildDestinationField(dest models.ReminderDestination, index int) *discordgo.MessageEmbedField {
	fieldName := fmt.Sprintf("📍 Destination %d", index)
	var fieldValue string

	switch dest.Type {
	case models.DestinationDiscordDM:
		if userID, exists := dest.Metadata["user_id"]; exists {
			if userIDStr, ok := userID.(string); ok {
				fieldValue = fmt.Sprintf("**Type:** Discord DM\n**User:** <@%s>", userIDStr)
			} else {
				fieldValue = "**Type:** Discord DM\n**User:** Unknown"
			}
		} else {
			fieldValue = "**Type:** Discord DM\n**User:** Invalid configuration"
		}

	case models.DestinationDiscordChannel:
		channelInfo := "**Type:** Discord Channel\n"
		
		if channelID, exists := dest.Metadata["channel_id"]; exists {
			if channelIDStr, ok := channelID.(string); ok {
				channelInfo += fmt.Sprintf("**Channel:** <#%s>\n", channelIDStr)
			} else {
				channelInfo += "**Channel:** Unknown\n"
			}
		} else {
			channelInfo += "**Channel:** Invalid configuration\n"
		}

		if guildID, exists := dest.Metadata["guild_id"]; exists {
			if guildIDStr, ok := guildID.(string); ok {
				channelInfo += fmt.Sprintf("**Server ID:** %s", guildIDStr)
			}
		}

		fieldValue = channelInfo

	case models.DestinationWebhook:
		if url, exists := dest.Metadata["url"]; exists {
			if urlStr, ok := url.(string); ok {
				// Mask the webhook URL for security
				maskedURL := maskWebhookURL(urlStr)
				fieldValue = fmt.Sprintf("**Type:** Webhook\n**URL:** %s", maskedURL)
			} else {
				fieldValue = "**Type:** Webhook\n**URL:** Invalid configuration"
			}
		} else {
			fieldValue = "**Type:** Webhook\n**URL:** Invalid configuration"
		}

		// Add any additional webhook metadata
		if name, exists := dest.Metadata["name"]; exists {
			if nameStr, ok := name.(string); ok {
				fieldValue += fmt.Sprintf("\n**Name:** %s", nameStr)
			}
		}

	default:
		fieldValue = fmt.Sprintf("**Type:** Unknown (%s)\n**Configuration:** %v", dest.Type, dest.Metadata)
	}

	return &discordgo.MessageEmbedField{
		Name:   fieldName,
		Value:  fieldValue,
		Inline: false,
	}
}

// maskWebhookURL masks sensitive parts of a webhook URL
func maskWebhookURL(url string) string {
	// Split by '/' and mask the token part
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		// Find the token part (usually the last part and quite long)
		for i := len(parts) - 1; i >= 0; i-- {
			if len(parts[i]) > 20 { // Webhook tokens are typically long
				parts[i] = parts[i][:8] + "..." + parts[i][len(parts[i])-4:]
				break
			}
		}
	}
	return strings.Join(parts, "/")
}

// CanAccessReminder checks if the user has permission to access the reminder
func CanAccessReminder(interaction *discordgo.InteractionCreate, account *models.Account, reminder *models.Reminder) bool {
	// User can always access their own reminders
	if reminder.AccountID == account.ID {
		return true
	}

	// If in a server and user is admin, check if the reminder has a server destination
	if interaction.GuildID != "" && interaction.Member != nil {
		// Check if user has administrator permissions
		permissions := interaction.Member.Permissions
		isAdmin := (permissions & discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator

		if isAdmin {
			// Check if reminder has a destination for this server
			for _, dest := range reminder.Destinations {
				if dest.Type == models.DestinationDiscordChannel {
					if guildID, exists := dest.Metadata["guild_id"]; exists {
						if guildIDStr, ok := guildID.(string); ok && guildIDStr == interaction.GuildID {
							return true
						}
					}
				}
			}
		}
	}

	return false
}