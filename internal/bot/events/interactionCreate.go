package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/commands"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		err := commands.HandleCommand(s, i)
		if err != nil {
			log.Printf("[DISCORD_BOT] - ❌ Error handling command: %v", err)
		}
	case discordgo.InteractionApplicationCommandAutocomplete:
		err := commands.HandleAutocomplete(s, i)
		if err != nil {
			log.Printf("[DISCORD_BOT] - ❌ Error handling autocomplete: %v", err)
		}
	case discordgo.InteractionMessageComponent:
		err := handleMessageComponent(s, i)
		if err != nil {
			log.Printf("[DISCORD_BOT] - ❌ Error handling message component: %v", err)
		}
	}
}

func handleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	customID := i.MessageComponentData().CustomID
	
	switch customID {
	case "timezone_change_select":
		// Ensure user account exists
		account, err := services.EnsureDiscordUser(i.Member.User)
		if err != nil {
			return err
		}
		return commands.HandleTimezoneSelectMenu(s, i, account)
	}
	
	return nil
}