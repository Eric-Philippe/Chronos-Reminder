// Command migrate-credentials performs the one-shot cutover that moves the
// email / password / username credential off the legacy `app` identity rows and
// onto the accounts table.
//
// It is idempotent: once there are no `app` identities left it is a no-op.
//
//	go run ./cmd/migrate-credentials
package main

import (
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database"
	"gorm.io/gorm"
)

func main() {
	log.Println("[MIGRATE] - ⏳ Account credential migration starting")

	// Loads .env and config the same way the main binary does.
	config.Load()

	// Connects and runs AutoMigrate, which adds the new
	// accounts.email / accounts.username / accounts.password_hash columns.
	if err := database.Initialize(); err != nil {
		log.Fatalf("[MIGRATE] - ❌ Failed to initialize database: %v", err)
	}
	defer func() { _ = database.Close() }()

	db := database.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Copy credentials from each `app` identity up to its account.
		//    Read via raw SQL because the Go model no longer knows the `app`
		//    provider (its Scan would reject those rows).
		res := tx.Exec(`
			UPDATE accounts a
			SET email         = i.external_id,
			    password_hash = i.password_hash,
			    username      = COALESCE(a.username, i.username)
			FROM identities i
			WHERE i.account_id = a.id
			  AND i.provider = 'app'
		`)
		if res.Error != nil {
			return res.Error
		}
		log.Printf("[MIGRATE] - ✅ Copied credentials onto %d account(s)", res.RowsAffected)

		// 2. Delete the now-redundant `app` identity rows.
		res = tx.Exec(`DELETE FROM identities WHERE provider = 'app'`)
		if res.Error != nil {
			return res.Error
		}
		log.Printf("[MIGRATE] - ✅ Removed %d legacy app identity row(s)", res.RowsAffected)

		// 3. Drop the deprecated password_hash column from identities entirely.
		if err := tx.Exec(`ALTER TABLE identities DROP COLUMN IF EXISTS password_hash`).Error; err != nil {
			return err
		}
		log.Println("[MIGRATE] - ✅ Dropped identities.password_hash column")

		return nil
	})
	if err != nil {
		log.Fatalf("[MIGRATE] - ❌ Migration failed (rolled back): %v", err)
	}

	log.Println("[MIGRATE] - 🎉 Account credential migration complete")
}
