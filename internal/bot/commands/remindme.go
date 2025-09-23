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

// reminderHandler handles the reminder creation command
func reminderHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options

	var message string
	var reminderTime string
	var recurrenceType string = "ONCE" // Default to ONCE

	// Parse command options
	for _, option := range options {
		switch option.Name {
		case "message":
			message = option.StringValue()
		case "time":
			reminderTime = option.StringValue()
		case "recurrence":
			if option.StringValue() != "" {
				recurrenceType = option.StringValue()
			}
		}
	}

	// Parse the reminder time
	parsedTime, err := services.ParseReminderTime(reminderTime)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Time Format", 
			fmt.Sprintf("Could not parse the time '%s'. Please use a format like '15:30', 'tomorrow 3pm', '1h 30m', or '2023-12-25 15:30'.", reminderTime))
	}

	// Get recurrence type value
	recurrenceTypeValue, exists := services.RecurrenceTypeMap[strings.ToUpper(recurrenceType)]
	if !exists {
		return utils.SendError(session, interaction, "Invalid Recurrence Type", 
			fmt.Sprintf("Invalid recurrence type '%s'. Valid options are: ONCE, YEARLY, MONTHLY, WEEKLY, DAILY, HOURLY, WORKDAYS, WEEKEND.", recurrenceType))
	}

	// Create the reminder
	reminder := &models.Reminder{
		AccountID:   account.ID,
		RemindAtUTC: parsedTime.UTC(),
		Message:     message,
		Recurrence:  int16(services.BuildRecurrenceState(recurrenceTypeValue, false)),
	}

	// Save the reminder to database
	repo := database.GetRepositories()
	if err := repo.Reminder.Create(reminder); err != nil {
		return utils.SendError(session, interaction, "Database Error", 
			"Failed to save the reminder. Please try again later.")
	}

	// Create the discord_dm destination
	destination := &models.ReminderDestination{
		ReminderID: reminder.ID,
		Type:       models.DestinationDiscordDM,
		Metadata: models.JSONB{
			"user_id": interaction.User.ID,
		},
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
	accountWithTimezone, err := repo.Account.GetWithTimezone(account.ID)
	var displayTime string
	if err == nil && accountWithTimezone != nil && accountWithTimezone.Timezone != nil {
		// Convert to user's timezone for display
		userTZ, err := time.LoadLocation(accountWithTimezone.Timezone.Name)
		if err == nil {
			displayTime = parsedTime.In(userTZ).Format("Monday, January 2, 2006 at 3:04 PM MST")
		} else {
			displayTime = parsedTime.UTC().Format("Monday, January 2, 2006 at 3:04 PM UTC")
		}
	} else {
		displayTime = parsedTime.UTC().Format("Monday, January 2, 2006 at 3:04 PM UTC")
	}

	description := fmt.Sprintf("**Message:** %s\n**Remind Time:** %s\n**Recurrence:** %s\n\n%s", 
		message, displayTime, recurrenceType, recurrenceText)

	return utils.SendSuccess(session, interaction, "Reminder Created! ⏰", description)
}

// Register the reminder command
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "remindme",
			Emoji:            "⏰",
			CategoryName:     "Reminders",
			ShortDescription: "Create a new reminder",
			FullDescription:  "Create a new reminder that will be sent to you via direct message at the specified time",
			Usage:            "/remindme message:<text> time:<when> [recurrence:<type>]",
			Example:          "/remindme message:\"Take medicine\" time:\"15:30\" recurrence:daily",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "remindme",
			Description: "Create a new reminder",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The reminder message",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "When to remind you (e.g., '15:30', 'tomorrow 3pm', '1h 30m', '2023-12-25 15:30')",
					Required:    true,
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
		Run:          reminderHandler,
	})
}
