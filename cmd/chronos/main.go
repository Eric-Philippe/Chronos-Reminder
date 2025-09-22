package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ericp/chronos-bot-reminder/internal/bot"
	"github.com/ericp/chronos-bot-reminder/internal/database"
)

func main() {
	// Initialize database
	log.Println("[ALL] - ⏳ Initializing Chronos Reminder")
	
	if err := database.Initialize(); err != nil {
		log.Fatalf("[DATABASE] - ❌ Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("[DATABASE] - ❌ Error closing database: %v", err)
		}
	}()

	// Start Discord bot
	bot.StartDiscordSession()

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("[ALL] - 🛑 Shutting down...")
	bot.StopDiscordSession()
}
