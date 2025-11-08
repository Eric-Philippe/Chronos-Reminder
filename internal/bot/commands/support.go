package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

func supportHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	description := "Need help? Here are the resources available:\n\n" +
		"📖 **Documentation & Contact**: https://chronosrmd.com/\n" +
		"Visit our website for full documentation, FAQ, and contact information.\n\n" +
		"💬 **Official Discord Server**: https://discord.gg/m3MsM922QD\n" +
		"Join our community to get support, share feedback, and connect with other users."

	return utils.SendEmbed(session, interaction, "💡 Need Support?", description, nil)
}

func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "support",
			Emoji:            "💡",
			CategoryName:     "General",
			ShortDescription: "Get help and support resources",
			FullDescription:  "Access documentation, FAQs, contact information, and our official Discord server",
			Usage:            "/support",
			Example:          "/support",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "support",
			Description: "Get help and support resources",
		},
		NeedsAccount: false,
		Run:          supportHandler,
	})
}