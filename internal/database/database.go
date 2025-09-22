package database

import (
	"fmt"
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var repos *repositories.Repositories
var cfg = config.Load()

// Initialize sets up the database connection and runs migrations
func Initialize() error {
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbName, cfg.DbPassword)
	
	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	
	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("[DATABASE] - Failed to connect to database: %w", err)
	}

	log.Println("[DATABASE] - ✅ Connection established")

	// Run auto migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	// Seed initial data
	if err := seedData(); err != nil {
		return fmt.Errorf("[DATABASE] - Failed to seed data: %w", err)
	}
	
	// Initialize repositories after database connection is established
	repos = repositories.NewRepositories(DB)
	
	return nil
}

// runMigrations performs automatic database migrations
func runMigrations() error {	
	// AutoMigrate will create tables, missing columns and missing indexes
	err := DB.AutoMigrate(
		&models.Timezone{},
		&models.Account{},
		&models.Identity{},
	)
	
	if err != nil {
		return err
	}
	
	// Create unique constraint for provider + external_id combination
	if err := DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_identities_provider_external_id 
		ON identities(provider, external_id)
	`).Error; err != nil {
		return err
	}
	
	return nil
}

// seedData inserts initial timezone data if it doesn't exist
func seedData() error {
	// Check if timezones already exist
	var count int64
	DB.Model(&models.Timezone{}).Count(&count)
	
	if count > 0 {
		return nil
	}
		
	timezones := []models.Timezone{
		{Name: "International Date Line West", GMTOffset: -12.0},
		{Name: "Midway Island, Samoa", GMTOffset: -11.0},
		{Name: "Hawaii", GMTOffset: -10.0},
		{Name: "Alaska", GMTOffset: -9.0},
		{Name: "Pacific Time (US & Canada)", GMTOffset: -8.0},
		{Name: "Mountain Time (US & Canada)", GMTOffset: -7.0},
		{Name: "Central Time (US & Canada), Mexico City", GMTOffset: -6.0},
		{Name: "Eastern Time (US & Canada), Bogota, Lima", GMTOffset: -5.0},
		{Name: "Atlantic Time (Canada), Caracas, La Paz", GMTOffset: -4.0},
		{Name: "Newfoundland", GMTOffset: -3.5},
		{Name: "Brazil, Buenos Aires, Georgetown", GMTOffset: -3.0},
		{Name: "Mid-Atlantic", GMTOffset: -2.0},
		{Name: "Azores, Cape Verde Islands", GMTOffset: -1.0},
		{Name: "Western Europe Time, London, Lisbon, Casablanca", GMTOffset: 0.0},
		{Name: "Brussels, Copenhagen, Madrid, Paris", GMTOffset: 1.0},
		{Name: "Kaliningrad, South Africa", GMTOffset: 2.0},
		{Name: "Baghdad, Riyadh, Moscow, St. Petersburg", GMTOffset: 3.0},
		{Name: "Tehran", GMTOffset: 3.5},
		{Name: "Abu Dhabi, Muscat, Baku, Tbilisi", GMTOffset: 4.0},
		{Name: "Kabul", GMTOffset: 4.5},
		{Name: "Ekaterinburg, Islamabad, Karachi, Tashkent", GMTOffset: 5.0},
		{Name: "Bombay, Calcutta, Madras, New Delhi", GMTOffset: 5.5},
		{Name: "Kathmandu", GMTOffset: 5.75},
		{Name: "Almaty, Dhaka, Colombo", GMTOffset: 6.0},
		{Name: "Yangon, Bangkok, Hanoi, Jakarta", GMTOffset: 6.5},
		{Name: "Bangkok, Hanoi, Jakarta", GMTOffset: 7.0},
		{Name: "Beijing, Perth, Singapore, Hong Kong", GMTOffset: 8.0},
		{Name: "Tokyo, Seoul, Osaka, Sapporo, Yakutsk", GMTOffset: 9.0},
		{Name: "Darwin", GMTOffset: 9.5},
		{Name: "Eastern Australia, Guam, Vladivostok", GMTOffset: 10.0},
		{Name: "Magadan, Solomon Islands, New Caledonia", GMTOffset: 11.0},
		{Name: "Auckland, Wellington, Fiji, Kamchatka", GMTOffset: 12.0},
	}
	
	if err := DB.Create(&timezones).Error; err != nil {
		return err
	}
	
	log.Printf("[DATABASE] - ✅ Seeded %d timezones", len(timezones))
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// GetRepositories returns the repository instances
func GetRepositories() *repositories.Repositories {
	return repos
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
