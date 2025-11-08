package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/logic"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// HelpHandlerFunc is the actual help handler function that accepts command data
type HelpHandlerFunc func(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account, commands interface{}) error

// helpHandler will be set after logic package initializes
var helpHandler HelpHandlerFunc

// SetHelpHandler sets the help handler function
func SetHelpHandler(handler HelpHandlerFunc) {
	helpHandler = handler
}

// helpHandlerWrapper wraps the help handler and passes all commands
func helpHandlerWrapper(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Call the help handler with all registered commands
	return logic.HelpHandlerWithCommands(session, interaction, account, commands)
}

func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name: "help",
			Emoji: "📕",
			CategoryName: "General",
			ShortDescription: "Get help with bot commands",
			FullDescription: "Provides information about available commands and how to use them. Use this command to learn more about the bot's features and functionalities.",
			Usage: "/help [command:name]",
			Example: "/help or /help command:profile",
		},
		Data: &discordgo.ApplicationCommand{
			Name: "help",
			Description: "Get help with bot commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "command",
					Description: "The command to get help with (leave empty for general help)",
					Required: false,
				},
			},
		},
		NeedsAccount: false,
		Run: helpHandlerWrapper,
	})
}