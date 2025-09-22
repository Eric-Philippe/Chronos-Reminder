package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// pingHandler handles the ping command
func pingHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	repo := database.GetRepositories()
	timezones, err := repo.Timezone.GetAll()
	if err != nil {
		return err
	}
	

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Pong! Available timezones: %d", len(timezones)),
		},
	})
}

// Register the ping command
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "ping",
			Emoji:            "🏓",
			CategoryName:     "General",
			ShortDescription: "Ping the bot",
			FullDescription:  "Ping the bot and get a response",
			Usage:            "/ping",
			Example:          "/ping",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Ping the bot",
			Options:     nil,
		},
		NeedsAccount: false,
		Run: pingHandler,
	})
}
