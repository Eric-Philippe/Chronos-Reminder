package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ericp/chronos-bot-reminder/internal/api"
	"github.com/ericp/chronos-bot-reminder/internal/bot"
	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/engine"
)

func main() {
	log.Println("[ALL] - ⏳ Initializing Chronos Reminder")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatalf("[DATABASE] - ❌ Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("[DATABASE] - ❌ Error closing database: %v", err)
		}
	}()

	// Get repositories
	repos := database.GetRepositories()

	// Initialize and start API server
	apiServer := api.NewServer(cfg, repos)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Fatalf("[API] - ❌ Failed to start API server: %v", err)
		}
	}()

	// Start Discord bot
	bot.StartDiscordSession()

	// Start scheduler service
	engine.StartSchedulerService()

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("[ALL] - 🛑 Gracefully shutting down...")

	// Stop scheduler
	engine.StopSchedulerService()

	// Stop Discord bot
	bot.StopDiscordSession()

	// Stop API server
	if err := apiServer.Stop(); err != nil {
		log.Printf("[API] - ❌ Error stopping API server: %v", err)
	}
}
