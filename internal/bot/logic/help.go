package logic

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

const (
	chronosWebURL      = "https://chronosrmd.com"
	chronosLogoURL     = "https://snapfilething.homeserver-ericp.fr/uploads/logo_chronos_1762633109_5a630a73_.png"
	helpColorPrimary   = 0x5b5bff
	helpColorCategory  = 0x7b7bff
)

// CommandInfo represents command information
type CommandInfo struct {
	Name             string
	Emoji            string
	CategoryName     string
	ShortDescription string
	FullDescription  string
	Usage            string
	Example          string
	Options          []*discordgo.ApplicationCommandOption
}

// HelpHandlerWithCommands handles the help command and displays available commands
// commandsData should be []*commands.Command
func HelpHandlerWithCommands(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account, commandsData interface{}) error {
	options := interaction.ApplicationCommandData().Options
	
	// Check if a specific command was requested
	var requestedCommand string
	if len(options) > 0 && options[0].Name == "command" {
		requestedCommand = strings.ToLower(options[0].StringValue())
	}

	// Convert commandsData to CommandInfo slice using reflection
	allCommands := convertCommandsData(commandsData)

	if requestedCommand != "" {
		// Show help for a specific command
		for _, cmd := range allCommands {
			if strings.EqualFold(cmd.Name, requestedCommand) {
				return sendCommandDetailedHelp(session, interaction, cmd)
			}
		}
		// Command not found
		return sendCommandNotFound(session, interaction, requestedCommand)
	}

	// Show general help with all commands organized by category
	return sendGeneralHelp(session, interaction, allCommands)
}

// convertCommandsData converts raw commands data to CommandInfo using reflection
func convertCommandsData(commandsData interface{}) []CommandInfo {
	result := []CommandInfo{}
	
	if commandsData == nil {
		return result
	}

	// Use reflection to handle []*Command without importing the commands package
	v := reflect.ValueOf(commandsData)
	if v.Kind() != reflect.Slice {
		return result
	}

	for i := 0; i < v.Len(); i++ {
		cmd := v.Index(i).Elem() // Dereference pointer
		
		// Get Description field
		descField := cmd.FieldByName("Description")
		if !descField.IsValid() {
			continue
		}
		desc := descField.Interface()
		descValue := reflect.ValueOf(desc)

		// Get Data field (contains Options)
		dataField := cmd.FieldByName("Data")
		if !dataField.IsValid() {
			continue
		}
		data := dataField.Elem()
		optionsField := data.FieldByName("Options")

		// Extract description fields
		name := getString(descValue, "Name")
		emoji := getString(descValue, "Emoji")
		categoryName := getString(descValue, "CategoryName")
		shortDesc := getString(descValue, "ShortDescription")
		fullDesc := getString(descValue, "FullDescription")
		usage := getString(descValue, "Usage")
		example := getString(descValue, "Example")

		// Extract Options
		var options []*discordgo.ApplicationCommandOption
		if optionsField.IsValid() && optionsField.Kind() == reflect.Slice {
			optionsLen := optionsField.Len()
			for j := 0; j < optionsLen; j++ {
				optVal := optionsField.Index(j)
				if opt, ok := optVal.Interface().(*discordgo.ApplicationCommandOption); ok {
					options = append(options, opt)
				}
			}
		}

		result = append(result, CommandInfo{
			Name:             name,
			Emoji:            emoji,
			CategoryName:     categoryName,
			ShortDescription: shortDesc,
			FullDescription:  fullDesc,
			Usage:            usage,
			Example:          example,
			Options:          options,
		})
	}

	return result
}

// getString safely gets a string field from a reflected value
func getString(v reflect.Value, fieldName string) string {
	if !v.IsValid() {
		return ""
	}

	field := v.FieldByName(fieldName)
	if field.IsValid() && field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

// sendGeneralHelp sends the main help menu with all commands grouped by category
func sendGeneralHelp(session *discordgo.Session, interaction *discordgo.InteractionCreate, allCommands []CommandInfo) error {
	// Group commands by category
	categories := make(map[string][]CommandInfo)
	for _, cmd := range allCommands {
		category := cmd.CategoryName
		if category == "" {
			category = "Other"
		}
		categories[category] = append(categories[category], cmd)
	}

	// Sort categories
	var sortedCategories []string
	for category := range categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Strings(sortedCategories)

	// Build embed with all commands
	embed := &discordgo.MessageEmbed{
		Title:       "📚 Chronos Bot - Complete Help Guide",
		Description: "Welcome to **Chronos Bot**! Your ultimate reminder management solution. Select a command below to learn more, or visit our web platform to enhance your experience.",
		Color:       helpColorPrimary,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: chronosLogoURL,
		},
	}

	// Add fields for each category
	for _, category := range sortedCategories {
		cmds := categories[category]
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].Name < cmds[j].Name
		})

		fieldValue := ""
		for _, cmd := range cmds {
			emoji := cmd.Emoji
			if emoji == "" {
				emoji = "⚙️"
			}
			fieldValue += fmt.Sprintf("%s `/%s` - %s\n", emoji, cmd.Name, cmd.ShortDescription)
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📂 " + category,
			Value:  fieldValue,
			Inline: false,
		})
	}

	// Add a tips section
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:  "💡 Quick Tips",
		Value: "• Use `/help command:<name>` to learn more about a specific command\n• All times are converted to your timezone automatically\n• Visit our web platform for advanced features and analytics",
		Inline: false,
	})

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "Chronos Bot Reminder • Type /help command:<name> for more details",
		IconURL: session.State.User.AvatarURL(""),
	}

	// Create response with button
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style: discordgo.LinkButton,
							Label: "🌐 Visit Chronos Web Platform",
							URL:   chronosWebURL,
						},
					},
				},
			},
		},
	})
}

// sendCommandDetailedHelp sends detailed help for a specific command
func sendCommandDetailedHelp(session *discordgo.Session, interaction *discordgo.InteractionCreate, cmd CommandInfo) error {
	// Build the detailed description
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s %s", cmd.Emoji, strings.ToUpper(cmd.Name)),
		Description: cmd.FullDescription,
		Color:       helpColorCategory,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: chronosLogoURL,
		},
	}

	// Add category
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "📂 Category",
		Value:  cmd.CategoryName,
		Inline: true,
	})

	// Add usage
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "💬 Usage",
		Value:  fmt.Sprintf("`%s`", cmd.Usage),
		Inline: false,
	})

	// Add example
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "📝 Example",
		Value:  fmt.Sprintf("`%s`", cmd.Example),
		Inline: false,
	})

	// Add options if available
	if len(cmd.Options) > 0 {
		optionsStr := buildOptionsDescription(cmd.Options)
		if optionsStr != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "⚙️ Options",
				Value:  optionsStr,
				Inline: false,
			})
		}
	}

	// Add a helpful note about the web platform
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "🚀 Enhanced Features",
		Value:  "Visit our web platform for advanced reminder management, analytics, and more customization options!",
		Inline: false,
	})

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "Chronos Bot Reminder 1.0.0 • Learn more on our web platform",
		IconURL: session.State.User.AvatarURL(""),
	}

	// Create response with navigation buttons
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style: discordgo.LinkButton,
							Label: "🌐 Web Platform",
							URL:   chronosWebURL,
						},
					},
				},
			},
		},
	})
}

// sendCommandNotFound sends an error message when command is not found
func sendCommandNotFound(session *discordgo.Session, interaction *discordgo.InteractionCreate, commandName string) error {
	embed := &discordgo.MessageEmbed{
		Title:       "❌ Command Not Found",
		Description: fmt.Sprintf("The command `%s` was not found. Use `/help` to see all available commands.", commandName),
		Color:       0xf57c76, // Error color
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: chronosLogoURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Chronos Bot Reminder",
			IconURL: session.State.User.AvatarURL(""),
		},
	}

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// buildOptionsDescription builds a formatted string describing command options
func buildOptionsDescription(options []*discordgo.ApplicationCommandOption) string {
	if len(options) == 0 {
		return ""
	}

	var description string
	for _, option := range options {
		required := "optional"
		if option.Required {
			required = "required"
		}

		optionType := getOptionTypeName(option.Type)

		if option.Type == discordgo.ApplicationCommandOptionSubCommand || option.Type == discordgo.ApplicationCommandOptionSubCommandGroup {
			// For subcommands, just list the name and description
			description += fmt.Sprintf("• `%s` - %s\n", option.Name, option.Description)
		} else {
			// For regular options
			description += fmt.Sprintf("• `%s` (%s, %s) - %s\n", option.Name, optionType, required, option.Description)
		}
	}

	return strings.TrimSpace(description)
}

// getOptionTypeName returns a human-readable name for option types
func getOptionTypeName(optType discordgo.ApplicationCommandOptionType) string {
	switch optType {
	case discordgo.ApplicationCommandOptionString:
		return "text"
	case discordgo.ApplicationCommandOptionInteger:
		return "number"
	case discordgo.ApplicationCommandOptionBoolean:
		return "boolean"
	case discordgo.ApplicationCommandOptionUser:
		return "user"
	case discordgo.ApplicationCommandOptionChannel:
		return "channel"
	case discordgo.ApplicationCommandOptionRole:
		return "role"
	case discordgo.ApplicationCommandOptionMentionable:
		return "mentionable"
	case discordgo.ApplicationCommandOptionNumber:
		return "decimal"
	case discordgo.ApplicationCommandOptionSubCommand:
		return "subcommand"
	case discordgo.ApplicationCommandOptionSubCommandGroup:
		return "subcommand group"
	default:
		return "unknown"
	}
}