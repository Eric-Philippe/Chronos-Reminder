package bot

import (
	"errors"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/commands"
	"github.com/ericp/chronos-bot-reminder/internal/bot/events"
)

var ErrMissingToken = errors.New("[DISCORD_BOT] - missing Discord bot token")

var (
	discord *discordgo.Session
)

// GetDiscordSession return a singleton ready Discord session
func GetDiscordSession() *discordgo.Session {
	if discord == nil {
		session, err := newDiscordSession(os.Getenv("DISCORD_BOT_TOKEN"))
		if err != nil {
			log.Fatalf("[DISCORD_BOT] - ❌ Cannot create Discord session: %v", err)
		}

		session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

		session.AddHandler(events.InteractionCreate)
		session.AddHandler(events.Ready)

		err = session.Open()
		if err != nil {
			log.Fatalf("[DISCORD_BOT] - ❌ Cannot open Discord session: %v", err)
		}

		// Register commands after session is open
		commandsLength, err := commands.RegisterCommands(session)
		if err != nil {
			log.Fatalf("[DISCORD_BOT] - ❌ Cannot register commands: %v", err)
		}

		log.Printf("[DISCORD_BOT] - ✅ Registered %d commands", commandsLength)

		discord = session
	}
	return discord
}

// newDiscordSession creates a new Discord session
func newDiscordSession(token string) (*discordgo.Session, error) {
	if token == "" {
		return nil, ErrMissingToken
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// StartDiscordSession initializes and returns the Discord session
func StartDiscordSession() *discordgo.Session {
	return GetDiscordSession()
}

// StopDiscordSession gracefully closes the Discord session
func StopDiscordSession() {
	if discord != nil {
		err := discord.Close()
		if err != nil {
			log.Printf("[DISCORD_BOT] - ❌ Error closing Discord session: %v", err)
		}
		discord = nil
	}
}
