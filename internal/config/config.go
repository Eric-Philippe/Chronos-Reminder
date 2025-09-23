package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Config structure for application settings
type Config struct {
    DiscordBotToken string
    LoggingLevel    string
    DefaultTZ       string
    APIPort         string
    DbHost          string
    DbPort          string
    DbUser          string
    DbPassword      string
    DbName          string
}

var (
	defaultTimezoneID *uint
	defaultTZMutex    sync.RWMutex
)

// SetDefaultTimezoneID allows updating the cached default timezone ID
func SetDefaultTimezoneID(id uint) {
	defaultTZMutex.Lock()
	defer defaultTZMutex.Unlock()
	defaultTimezoneID = &id
}

// Load reads configuration from environment variables or .env file
func Load() *Config {
    // Load .env file, ignore error if file not found
    if err := godotenv.Load(); err != nil {
        log.Println("[ALL] - ⚠️ No .env file found, relying on environment variables")
    }

    cfg := &Config{
        DiscordBotToken: getEnv("DISCORD_BOT_TOKEN", ""),
        LoggingLevel:    getEnv("LOGGING_LEVEL", "INFO"),
        DefaultTZ:       getEnv("DEFAULT_TZ", "15"), // Default to UTC+1
        APIPort:         getEnv("API_PORT", "8080"),
        DbHost:          getEnv("DB_HOST", "localhost"),
        DbPort:          getEnv("DB_PORT", "5432"),
        DbUser:          getEnv("DB_USER", "user"),
        DbPassword:      getEnv("DB_PASSWORD", "password"),
        DbName:          getEnv("DB_NAME", "ChronosReminder"),
    }

    return cfg
}

func GetDatabaseConfig() *Config {
    return &Config{
        DbHost:     os.Getenv("DB_HOST"),
        DbPort:     os.Getenv("DB_PORT"),
        DbUser:     os.Getenv("DB_USER"),
        DbPassword: os.Getenv("DB_PASSWORD"),
        DbName:     os.Getenv("DB_NAME"),
    }
}

func IsDebugMode() bool {
    cfg := Load()
    return cfg.LoggingLevel == "DEBUG"
}

// GetDefaultTimezoneID returns the cached default timezone ID as *uint
func GetDefaultTimezoneID() *uint {
	defaultTZMutex.RLock()
	defer defaultTZMutex.RUnlock()
	
	if defaultTimezoneID == nil {
		cfg := Load()
		if id, err := strconv.ParseUint(cfg.DefaultTZ, 10, 32); err == nil {
			uintID := uint(id)
			defaultTimezoneID = &uintID
		} else {
			log.Printf("[DATABASE] - ⚠️ Invalid default timezone ID in config: %s, error: %v", cfg.DefaultTZ, err)
			// Fallback to timezone ID 1 if parsing fails
			fallbackID := uint(1)
			defaultTimezoneID = &fallbackID
		}
	}
	
	// Return a copy to avoid external modification
	if defaultTimezoneID != nil {
		id := *defaultTimezoneID
		return &id
	}
	return nil
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
