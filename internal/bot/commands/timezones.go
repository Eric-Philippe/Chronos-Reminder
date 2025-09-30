package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/logic"
)

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
		Run:          logic.TimezoneHandler,
	})
}
