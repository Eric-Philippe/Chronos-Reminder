package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

func hasPerm(perms int64, perm int64) bool {
    return perms&perm == perm
}

// remindUsHandler handles the remind us command
func remindUsHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Check if the command is being used in a server (not DM)
	if interaction.GuildID == "" {
		return utils.SendError(session, interaction, "Server Required", 
			"The `/remindus` command can only be used in a server, not in direct messages. Use `/remindme` for personal reminders.")
	}

	options := interaction.ApplicationCommandData().Options

	var message string
	var dateStr string
	var timeStr string
	var channelID string
	var roleID string
	var recurrenceType string = "ONCE" // Default to ONCE

	// Parse command options
	for _, option := range options {
		switch option.Name {
		case "message":
			message = option.StringValue()
		case "date":
			dateStr = option.StringValue()
		case "time":
			timeStr = option.StringValue()
		case "channel":
			if channel := option.ChannelValue(session); channel != nil {
				channelID = channel.ID
			} else if option.Value != nil {
				// Fallback: try to get channel ID directly from the option value
				if channelIDStr, ok := option.Value.(string); ok {
					channelID = channelIDStr
				}
			}
		case "role":
			if role := option.RoleValue(session, interaction.GuildID); role != nil {
				roleID = role.ID
			} else if option.Value != nil {
				// Fallback: try to get role ID directly from the option value
				if roleIDStr, ok := option.Value.(string); ok {
					roleID = roleIDStr
				}
			}
		case "recurrence":
			if option.StringValue() != "" {
				recurrenceType = option.StringValue()
			}
		}
	}

	// Validate that a channel was selected
	if channelID == "" {
		return utils.SendError(session, interaction, "Channel Required", 
			"Please select a channel where the reminder should be sent.")
	}

	// Validate required fields
	if message == "" {
		return utils.SendError(session, interaction, "Message Required", 
			"Please provide a message for the reminder.")
	}
	
	if dateStr == "" {
		return utils.SendError(session, interaction, "Date Required", 
			"Please provide a date for the reminder.")
	}
	
	if timeStr == "" {
		return utils.SendError(session, interaction, "Time Required", 
			"Please provide a time for the reminder.")
	}

	// Get the user ID for permission checking
	var userID string
	if interaction.Member != nil && interaction.Member.User != nil {
		userID = interaction.Member.User.ID
	} else if interaction.User != nil {
		userID = interaction.User.ID
	} else {
		return utils.SendError(session, interaction, "User Information Missing", 
			"Could not determine user information for permission check.")
	}

	// Verify the user has manage channel permissions, administrator permissions, or is the server owner
	channelPerms, err := session.UserChannelPermissions(userID, channelID)
	if err != nil {
		return utils.SendError(session, interaction, "Permission Check Failed", 
			"Could not verify your permissions for the selected channel.")
	}

	userPerms := interaction.Member.Permissions

	// Check if user is server owner
	guild, err := session.Guild(interaction.GuildID)
	isAllowed := err == nil && (guild.OwnerID == userID || hasPerm(userPerms, discordgo.PermissionAdministrator)|| hasPerm(userPerms, discordgo.PermissionManageChannels) || hasPerm(channelPerms, discordgo.PermissionManageChannels))

	if !isAllowed {
		return utils.SendError(session, interaction, "Insufficient Permissions", 
			"You need 'Manage Channel', 'Administrator' permission, or be the server owner to create reminders in the selected channel.")
	}

	// If a role is specified, validate role mention permissions
	if roleID != "" {
		// Check if bot has permission to mention roles
		// First check guild-wide permissions (for Administrator)
		botMember, err := session.GuildMember(interaction.GuildID, session.State.User.ID)
		if err != nil {
			return utils.SendError(session, interaction, "Bot Permission Check Failed", 
				"Could not verify bot's permissions to mention roles.")
		}

		// Check guild-wide permissions first
		var botHasPermission bool
		for _, botRoleID := range botMember.Roles {
			role, err := session.State.Role(interaction.GuildID, botRoleID)
			if err != nil {
				// If role not in state, try to fetch from API
				roles, apiErr := session.GuildRoles(interaction.GuildID)
				if apiErr != nil {
					continue // Skip this role if we can't fetch it
				}
				// Find the role in the API response
				for _, r := range roles {
					if r.ID == botRoleID {
						role = r
						break
					}
				}
				if role == nil {
					continue // Skip if role not found
				}
			}
			
			if role.Permissions&discordgo.PermissionAdministrator != 0 || 
				role.Permissions&discordgo.PermissionMentionEveryone != 0 || 
				role.Permissions&discordgo.PermissionManageRoles != 0 {
				botHasPermission = true
				break
			}
		}

		// If no guild-wide permission found, check channel-specific permissions
		if !botHasPermission {
			botPerms, err := session.UserChannelPermissions(session.State.User.ID, channelID)
			if err == nil && (botPerms&discordgo.PermissionMentionEveryone != 0 || 
				botPerms&discordgo.PermissionManageRoles != 0 || 
				botPerms&discordgo.PermissionAdministrator != 0) {
				botHasPermission = true
			}
		}

		if !botHasPermission {
			return utils.SendError(session, interaction, "Bot Insufficient Permissions", 
				"The bot needs 'Mention Everyone', 'Manage Roles', or 'Administrator' permission to mention roles in reminders.")
		}

		// Check if user has permission to manage the specified role (unless they're owner/admin)
		if isAllowed {
			// Get the role to check hierarchy
			role, err := session.State.Role(interaction.GuildID, roleID)
			if err != nil {
				// Try to fetch from API if not in state
				roles, err := session.GuildRoles(interaction.GuildID)
				if err != nil {
					return utils.SendError(session, interaction, "Role Validation Failed", 
						"Could not validate the specified role.")
				}
				for _, r := range roles {
					if r.ID == roleID {
						role = r
						break
					}
				}
			}

			if role == nil {
				return utils.SendError(session, interaction, "Invalid Role", 
					"The specified role could not be found.")
			}

			// Check if user has manage roles permission
			if !hasPerm(userPerms, discordgo.PermissionManageRoles) {
				return utils.SendError(session, interaction, "Role Permission Required", 
					"You need 'Manage Roles' permission to mention roles in reminders.")
			}

			// Check role hierarchy - user's highest role must be higher than the role they want to mention
			member, err := session.GuildMember(interaction.GuildID, userID)
			if err != nil {
				return utils.SendError(session, interaction, "Member Information Missing", 
					"Could not verify your role hierarchy.")
			}

			userHighestRole := 0
			for _, userRoleID := range member.Roles {
				userRole, err := session.State.Role(interaction.GuildID, userRoleID)
				if err == nil && userRole.Position > userHighestRole {
					userHighestRole = userRole.Position
				}
			}

			if role.Position >= userHighestRole {
				return utils.SendError(session, interaction, "Role Hierarchy Error", 
					"You can only mention roles that are lower in the hierarchy than your highest role.")
			}
		}
	}


	// Parse the reminder date and time in user's timezone
	parsedTime, err := services.ParseReminderDateTimeInTimezone(dateStr, timeStr, account.Timezone.IANALocation)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Date/Time Format", 
			fmt.Sprintf("Could not parse the date '%s' and time '%s'. Please check your date and time formats.", dateStr, timeStr))
	}

		location, err := time.LoadLocation(account.Timezone.IANALocation)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Timezone", 
			fmt.Sprintf("Could not load timezone '%s'. Please check your timezone settings.", account.Timezone.IANALocation))
	}
	now := time.Now().In(location)
	// If the parsed reminder time is before the current time, return an error
	if parsedTime.Before(now) {
		return utils.SendError(session, interaction, "Invalid Date/Time", 
			"The specified date and time is in the past. Please provide a future date and time for the reminder.")
	}

	// Get recurrence type value
	recurrenceTypeValue, exists := services.RecurrenceTypeMap[strings.ToUpper(recurrenceType)]
	if !exists {
		return utils.SendError(session, interaction, "Invalid Recurrence Type", 
			fmt.Sprintf("Invalid recurrence type '%s'. Valid options are: ONCE, YEARLY, MONTHLY, WEEKLY, DAILY, HOURLY, WORKDAYS, WEEKEND.", recurrenceType))
	}

	// Create the reminder with UTC time
	reminder := &models.Reminder{
		AccountID:   account.ID,
		RemindAtUTC: parsedTime.UTC(),
		Message:     message,
		Recurrence:  int16(services.BuildRecurrenceState(recurrenceTypeValue, false)),
	}

	repo := database.GetRepositories()

	// Save the reminder to database
	if err := repo.Reminder.Create(reminder, true); err != nil {
		return utils.SendError(session, interaction, "Database Error", 
			"Failed to save the reminder. Please try again later.")
	}

	// Create the discord_channel destination
	destinationMetadata := models.JSONB{
		"guild_id":   interaction.GuildID,
		"channel_id": channelID,
	}
	
	// Add role mention if specified
	if roleID != "" {
		destinationMetadata["mention_role_id"] = roleID
	}
	
	destination := &models.ReminderDestination{
		ReminderID: reminder.ID,
		Type:       models.DestinationDiscordChannel,
		Metadata:   destinationMetadata,
	}

	if err := repo.ReminderDestination.Create(destination); err != nil {
		// If destination creation fails, we should clean up the reminder
		repo.Reminder.Delete(reminder.ID, true)
		return utils.SendError(session, interaction, "Database Error", 
			"Failed to set up reminder destination. Please try again later.")
	}

	// Format response message
	var recurrenceText string
	if recurrenceType == "ONCE" {
		recurrenceText = "This is a one-time reminder."
	} else {
		recurrenceText = fmt.Sprintf("This reminder will repeat: %s", strings.ToLower(recurrenceType))
	}

	// Load account timezone for display
	var displayTime string
	if account != nil && account.Timezone != nil {
		// Display the local time as entered by the user
		displayTime = parsedTime.Format("Monday, January 2, 2006 at 15:04")
	} else {
		// Display in the same timezone as the parsed time was created
		displayTime = parsedTime.Format("Monday, January 2, 2006 at 15:04")
	}

	description := fmt.Sprintf("**Content:** %s\n**Remind Time:** %s\n**Channel:** <#%s>", 
		message, displayTime, channelID)
	
	// Add role mention info if specified
	if roleID != "" {
		description += fmt.Sprintf("\n**Role Mention:** <@&%s>", roleID)
	}

	return utils.SendEmbed(session, interaction, "Channel Reminder Created! 📢", description, &recurrenceText)
}

func init() {
	autocompleteFunc := AutocompleteFunc(DateAutocompleteHandler)

	RegisterCommand(&Command{
		Description: Description{
			Name:             "remindus",
			Emoji:            "📢",
			CategoryName:     "Reminders",
			ShortDescription: "Create a new reminder in a channel",
			FullDescription:  "Create a new reminder that will be sent in a specified channel at the specified date and time. Requires 'Manage Channel', 'Administrator' permission, or server ownership.",
			Usage:            "/remindus message:<text> date:<date> time:<time> channel:<channel> [role:<role>] [recurrence:<type>]",
			Example:          "/remindus message:\"Team meeting\" date:\"25/12/2024\" time:\"10:00\" channel:#general role:@developers recurrence:weekly",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "remindus",
			Description: "Create a new reminder in a channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The reminder message",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date",
					Description: "The date for the reminder (e.g., 'today', 'tomorrow', '25/12/2024', '2024-12-25')",
					Required:    true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "The time for the reminder (e.g., '15:30', '3pm', '9:30am')",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel to send the reminder in",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role to mention in the reminder (optional)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "recurrence",
					Description: "How often to repeat (default: once)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Once",
							Value: "ONCE",
						},
						{
							Name:  "Hourly",
							Value: "HOURLY",
						},
						{
							Name:  "Daily",
							Value: "DAILY",
						},
						{
							Name:  "Weekly",
							Value: "WEEKLY",
						},
						{
							Name:  "Monthly",
							Value: "MONTHLY",
						},
						{
							Name:  "Yearly",
							Value: "YEARLY",
						},
						{
							Name:  "Workdays (Mon-Fri)",
							Value: "WORKDAYS",
						},
						{
							Name:  "Weekends (Sat-Sun)",
							Value: "WEEKEND",
						},
					},
				},
			},
		},
		NeedsAccount: true,
		Run:          remindUsHandler,
		Autocomplete: &autocompleteFunc,
	})
}