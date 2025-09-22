package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("[DISCORD_BOT] - 🤖 Bot is ready! Logged in as: %s#%s", event.User.Username, event.User.Discriminator)
}