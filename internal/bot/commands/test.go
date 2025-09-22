package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// testHandler handles the test command
func testHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	// Get the cached default timezone ID from config
	responseContent := account.ID.String()

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseContent,
		},
	})
}

// Register the test command
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "test",
			Emoji:            "🧪",
			CategoryName:     "General",
			ShortDescription: "Test command",
			FullDescription:  "A test command to check if the bot is responsive",
			Usage:            "/test",
			Example:          "/test",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "test",
			Description: "Test the bot",
			Options:     nil,
		},
		NeedsAccount: true,
		Run: testHandler,
	})
}
