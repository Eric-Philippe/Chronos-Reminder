package commands

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// timezoneHandler handles the timezone command with subcommands
func timezoneHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options
	if len(options) == 0 {
		return timezoneListHandler(session, interaction, account)
	}

	subcommand := options[0]
	switch subcommand.Name {
	case "list":
		return timezoneListHandler(session, interaction, account)
	case "change":
		return timezoneChangeHandler(session, interaction, account)
	default:
		return timezoneListHandler(session, interaction, account)
	}
}

// timezoneListHandler handles the timezone list subcommand
func timezoneListHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, _ *models.Account) error {
	repo := database.GetRepositories()
	timezones, err := repo.Timezone.GetAll()
	if err != nil {
		return err
	}

	if len(timezones) == 0 {
		// Build embed message
		embed := &discordgo.MessageEmbed{
			Title:       "No Timezones Available",
			Description: "There are currently no timezones available.",
			Color:       0xff0000, // Red color
		}

		return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	}

	// Create a formatted list of timezones
	timezoneList := ""
	for _, tz := range timezones {
		gmtOffsetStr := ""
		if tz.GMTOffset >= 0 {
			gmtOffsetStr = fmt.Sprintf("GMT+%.1f", tz.GMTOffset)
		} else {
			gmtOffsetStr = fmt.Sprintf("GMT%.1f", tz.GMTOffset)
		}
		timezoneList += fmt.Sprintf("• %s (%s)\n", tz.Name, gmtOffsetStr)
	}

	// Build embed message
	embed := &discordgo.MessageEmbed{
		Title:       "Available Timezones",
		Description: timezoneList,
		Color:       0x00ff00, // Green color
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// timezoneChangeHandler handles the timezone change subcommand
func timezoneChangeHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, _ *models.Account) error {
	repo := database.GetRepositories()
	timezones, err := repo.Timezone.GetAll()
	if err != nil {
		return err
	}

	if len(timezones) == 0 {
		embed := &discordgo.MessageEmbed{
			Title:       "No Timezones Available",
			Description: "There are currently no timezones available to change to.",
			Color:       0xff0000, // Red color
		}

		return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	}

	// Create select menu options
	var options []discordgo.SelectMenuOption
	for _, tz := range timezones {
		gmtOffsetStr := ""
		if tz.GMTOffset >= 0 {
			gmtOffsetStr = fmt.Sprintf("GMT+%.1f", tz.GMTOffset)
		} else {
			gmtOffsetStr = fmt.Sprintf("GMT%.1f", tz.GMTOffset)
		}
		
		description := gmtOffsetStr
		if len(description) > 100 {
			description = description[:97] + "..."
		}

		options = append(options, discordgo.SelectMenuOption{
			Label:       tz.Name,
			Value:       strconv.Itoa(int(tz.ID)),
			Description: description,
		})
	}

	// Limit to 25 options (Discord's limit)
	if len(options) > 25 {
		options = options[:25]
	}

	selectMenu := &discordgo.SelectMenu{
		CustomID:    "timezone_change_select",
		Placeholder: "Choose a timezone...",
		Options:     options,
	}

	embed := utils.BuildInfoEmbed(session, "Change Timezone", "Select your new timezone from the dropdown menu below.")

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{selectMenu},
				},
			},
		},
	})
}

// HandleTimezoneSelectMenu handles the timezone selection from the dropdown
func HandleTimezoneSelectMenu(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	if len(interaction.MessageComponentData().Values) == 0 {
		return fmt.Errorf("no timezone selected")
	}

	timezoneIDStr := interaction.MessageComponentData().Values[0]
	timezoneID, err := strconv.ParseUint(timezoneIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid timezone ID: %v", err)
	}

	// Change the user's timezone
	err = services.ChangeAccountTimezone(account, uint(timezoneID))
	if err != nil {
		return utils.SendError(session, interaction, "Error", "Failed to change your timezone. Please try again.")
	}

	// Get the timezone name for confirmation
	repo := database.GetRepositories()
	timezone, err := repo.Timezone.GetByID(uint(timezoneID))
	if err != nil || timezone == nil {
		return utils.SendError(session, interaction, "Error", "Failed to retrieve the selected timezone details.")
	}

	gmtOffsetStr := ""
	if timezone.GMTOffset >= 0 {
		gmtOffsetStr = fmt.Sprintf("GMT+%.1f", timezone.GMTOffset)
	} else {
		gmtOffsetStr = fmt.Sprintf("GMT%.1f", timezone.GMTOffset)
	}

	return utils.SendSuccess(session, interaction, "Timezone Changed", fmt.Sprintf("Your timezone has been successfully changed to **%s** (%s)!", timezone.Name, gmtOffsetStr), nil)
}

// Register the timezone command with subcommands
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "timezone",
			Emoji:            "🌏",
			CategoryName:     "General",
			ShortDescription: "Manage timezones",
			FullDescription:  "List available timezones or change your current timezone",
			Usage:            "/timezone <list|change>",
			Example:          "/timezone list, /timezone change",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "timezone",
			Description: "Manage timezones",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "List all available timezones",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "change",
					Description: "Change your current timezone",
				},
			},
		},
		NeedsAccount: true,
		Run:          timezoneHandler,
		MessageComponentHandlers: []MessageComponentHandler{
			{
				CustomID: "timezone_change_select",
				Handler:  HandleTimezoneSelectMenu,
			},
		},
	})
}
