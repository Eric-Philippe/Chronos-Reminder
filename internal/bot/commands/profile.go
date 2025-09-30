package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/logic"
)

// Register the profile command
func init() {
	RegisterCommand(&Command{
		Description: Description{
			Name:             "profile",
			Emoji:            "👤",
			CategoryName:     "User",
			ShortDescription: "View user profile",
			FullDescription:  "Display a user's profile with their avatar, creation date, reminder count, and platform badges. Use without parameters to view your own profile, or specify a user to view theirs.",
			Usage:            "/profile [user:@user]",
			Example:          "/profile or /profile user:@username",
		},
		Data: &discordgo.ApplicationCommand{
			Name:        "profile",
			Description: "View a user's profile information",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user whose profile to view (leave empty for your own profile)",
					Required:    false,
				},
			},
		},
		NeedsAccount: true,
		Run:          logic.ProfileHandler,
	})
}