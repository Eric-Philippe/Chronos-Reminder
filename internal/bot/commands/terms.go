package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

func termsHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	description := "📋 **Terms & Privacy Policy**: " + config.URLWebApp + "/terms\n\n" +
		"Chronos is a personal project with no company, no investors, and no monetization plans.\n\n" +
		"**Key points:**\n" +
		"🔒 Your password is bcrypt-hashed - nobody can read it, including the developer\n" +
		"🚫 Your data is never shared with or sold to any third party\n" +
		"🗑️ You can delete your account and all associated data at any time\n" +
		"💻 Self-hosting is available if you want full control over your data"

	return utils.SendEmbed(session, interaction, "📜 Terms & Privacy", description, nil)
}

func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "terms",
			Emoji:            "📜",
			CategoryName:     "General",
			ShortDescription: "View the terms and privacy policy",
			FullDescription:  "Display a summary of Chronos's terms of use and privacy policy, with a link to the full page",
			Usage:            "/terms",
			Example:          "/terms",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "terms",
			Description: "View the terms and privacy policy",
		},
		NeedsAccount: false,
		Run:          termsHandler,
	})
}
