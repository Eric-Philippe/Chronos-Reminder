package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// ticHandler handles the tic command
func ticHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	return utils.SendInfo(session, interaction, "The bot is alive!", "⏰ Tac !")
}

// Register the tic command
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "tic",
			Emoji:            "⏰",
			CategoryName:     "General",
			ShortDescription: "Ping the bot",
			FullDescription:  "Ping the bot and get a response",
			Usage:            "/tic",
			Example:          "/tic",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "tic",
			Description: "Tick the bot",
			Options:     nil,
		},
		NeedsAccount: false,
		Run: 	  ticHandler,
	})
}
