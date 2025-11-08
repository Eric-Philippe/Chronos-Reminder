package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// parseDurationString parses duration strings with various formats (e.g., "10s", "5m")
// Only accepts seconds and minutes, not hours
func parseDurationString(input string) (time.Duration, error) {
	// Try parsing as standard Go duration (but reject hours)
	if d, err := time.ParseDuration(input); err == nil {
		if d >= time.Hour {
			return 0, fmt.Errorf("hours are not supported for timers, use /remindme or /remindus for longer reminders")
		}
		return d, nil
	}

	return 0, fmt.Errorf("unable to parse duration: %s", input)
}

// formatDuration formats a duration in a human-readable way
func formatDurationHourglass(d time.Duration) string {
	if d < 0 {
		return "-" + formatDurationHourglass(-d)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}
	return result
}

func hourglassHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options

	var duration string
	var message string

	// Parse command options
	for _, option := range options {
		switch option.Name {
		case "duration":
			duration = option.StringValue()
		case "message":
			message = option.StringValue()
		}
	}

	// Parse the duration
	parsedDuration, err := parseDurationString(duration)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Duration Format",
			fmt.Sprintf("Could not parse duration '%s'. Use formats like '10s' or '5m'. For longer reminders, use /remindme or /remindus commands.", duration))
	}

	// Validate duration is not too long (max 30 minutes for in-memory timer)
	if parsedDuration > 30*time.Minute {
		return utils.SendError(session, interaction, "Duration Too Long",
			"The timer duration cannot exceed 30 minutes. For longer reminders, use /remindme or /remindus commands.")
	}

	// Validate duration is not zero or negative
	if parsedDuration <= 0 {
		return utils.SendError(session, interaction, "Invalid Duration",
			"The timer duration must be greater than 0.")
	}

	// Get the user ID
	var userID string
	if interaction.Member != nil && interaction.Member.User != nil {
		userID = interaction.Member.User.ID
	} else if interaction.User != nil {
		userID = interaction.User.ID
	}

	// Get the channel ID
	channelID := interaction.ChannelID

	// Send immediate acknowledgement
	endTime := time.Now().Add(parsedDuration)
	description := fmt.Sprintf("**Message:** %s\n**Duration:** %s\n**End Time:** <t:%d:t>",
		message, formatDurationHourglass(parsedDuration), endTime.Unix())

	if err := utils.SendEmbed(session, interaction, "⏳ Timer Started", description, nil); err != nil {
		return err
	}

	// Start timer in a goroutine
	go func() {
		time.Sleep(parsedDuration)

		// Send follow-up message when timer ends
		embed := &discordgo.MessageEmbed{
			Title:       "⏳ Timer Finished!",
			Description: fmt.Sprintf("<@%s> - %s", userID, message),
			Color:       utils.ColorSuccess,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: utils.ClockLogo,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Chronos Bot Reminder",
			},
		}

		_, err := session.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			fmt.Printf("Error sending timer completion message: %v\n", err)
		}
	}()

	return nil
}

func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "hourglass",
			Emoji:            "⏳",
			CategoryName:     "General",
			ShortDescription: "Start a short timer (max 30 minutes)!",
			FullDescription:  "Start a quick in-memory timer that will notify you when it ends. For longer or persistent reminders, use /remindme or /remindus.",
			Usage:            "/hourglass duration:10m message:Take a break!",
			Example:          "/hourglass duration:25m message:Focus time!",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "hourglass",
			Description: "Start a short timer",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "duration",
					Description: "Duration in seconds or minutes (max 30min) - e.g., '10s', '5m'",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The message to send when the timer ends",
					Required:    true,
				},
			},
		},
		NeedsAccount: false,
		Run:          hourglassHandler,
	})
}