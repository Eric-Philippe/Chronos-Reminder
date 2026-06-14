package services

import (
	"context"
	"fmt"
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


// MergeAccounts merges mergedID into survivorID in a single transaction.
// All data (reminders, identities, FCM tokens, DFM note/items) is moved to the
// survivor. The merged account is deleted. If the survivor already has
// email/password credentials, they are kept; otherwise the merged account's
// credentials are adopted.
func MergeAccounts(ctx context.Context, repos *repositories.Repositories, survivorID, mergedID uuid.UUID) error {
	if survivorID == mergedID {
		return fmt.Errorf("cannot merge account into itself")
	}

	db := database.GetDB()

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Re-point all reminders from merged to survivor
		if err := tx.Model(&models.Reminder{}).
			Where("account_id = ?", mergedID).
			Update("account_id", survivorID).Error; err != nil {
			return fmt.Errorf("re-pointing reminders: %w", err)
		}

		// Re-point FCM tokens
		if err := tx.Model(&models.FcmToken{}).
			Where("account_id = ?", mergedID).
			Update("account_id", survivorID).Error; err != nil {
			return fmt.Errorf("re-pointing fcm tokens: %w", err)
		}

		// Re-point identities (discord / mobile / api_key rows of merged account)
		var mergedIdentities []models.Identity
		if err := tx.Where("account_id = ?", mergedID).Find(&mergedIdentities).Error; err != nil {
			return fmt.Errorf("fetching merged identities: %w", err)
		}

		for i := range mergedIdentities {
			id := &mergedIdentities[i]
			if id.Provider == models.ProviderMobile {
				// mobile external_id == account.id; fix it to survivor's id
				survivorIDStr := survivorID.String()
				id.ExternalID = survivorIDStr
				// Check whether survivor already has a mobile identity
				var existing models.Identity
				err := tx.Where("account_id = ? AND provider = ?", survivorID, models.ProviderMobile).First(&existing).Error
				if err == nil {
					// Survivor already has mobile identity — drop this duplicate
					if err := tx.Delete(id).Error; err != nil {
						return fmt.Errorf("deleting duplicate mobile identity: %w", err)
					}
					continue
				}
			}
			id.AccountID = survivorID
			if err := tx.Save(id).Error; err != nil {
				return fmt.Errorf("moving identity %s: %w", id.ID, err)
			}
		}

		// Handle DFM notes — one per account
		var survivorNote, mergedNote models.DFMNote
		hasSurvivorNote := tx.Where("account_id = ?", survivorID).First(&survivorNote).Error == nil
		hasMergedNote := tx.Where("account_id = ?", mergedID).First(&mergedNote).Error == nil

		if hasMergedNote {
			if !hasSurvivorNote {
				// Just re-point the merged note
				mergedNote.AccountID = survivorID
				if err := tx.Save(&mergedNote).Error; err != nil {
					return fmt.Errorf("re-pointing dfm note: %w", err)
				}
			} else {
				// Both have notes: append merged items to survivor, then delete merged note
				var mergedItems []models.DFMItem
				if err := tx.Where("note_id = ?", mergedNote.ID).Order("position").Find(&mergedItems).Error; err != nil {
					return fmt.Errorf("fetching merged dfm items: %w", err)
				}
				// Find max position in survivor's note
				var maxPos int
				tx.Model(&models.DFMItem{}).Where("note_id = ?", survivorNote.ID).
					Select("COALESCE(MAX(position), -1)").Scan(&maxPos)

				for i := range mergedItems {
					item := &mergedItems[i]
					item.NoteID = survivorNote.ID
					item.Position = maxPos + 1 + i
					item.ID = uuid.New()
					if err := tx.Create(item).Error; err != nil {
						return fmt.Errorf("moving dfm item: %w", err)
					}
				}
				// Delete merged note (items were already re-created under survivor)
				if err := tx.Where("note_id = ?", mergedNote.ID).Delete(&models.DFMItem{}).Error; err != nil {
					return fmt.Errorf("deleting old dfm items: %w", err)
				}
				if err := tx.Delete(&mergedNote).Error; err != nil {
					return fmt.Errorf("deleting merged dfm note: %w", err)
				}
			}
		}

		// Adopt merged account's credentials only if survivor has none
		survivor, err := repos.Account.GetByID(survivorID)
		if err != nil || survivor == nil {
			return fmt.Errorf("loading survivor account: %w", err)
		}
		merged, err := repos.Account.GetByID(mergedID)
		if err != nil || merged == nil {
			return fmt.Errorf("loading merged account: %w", err)
		}

		if survivor.Email == nil && merged.Email != nil {
			if err := tx.Model(&models.Account{}).Where("id = ?", survivorID).Updates(map[string]interface{}{
				"email":          merged.Email,
				"password_hash":  merged.PasswordHash,
				"username":       merged.Username,
				"email_verified": merged.EmailVerified,
			}).Error; err != nil {
				return fmt.Errorf("adopting credentials: %w", err)
			}
		}

		// Delete transient rows for merged account
		tx.Where("account_id = ?", mergedID.String()).Delete(&models.EmailVerification{})
		tx.Where("account_id = ?", mergedID).Delete(&models.PasswordReset{})

		// Delete the merged account (FK cascade covers any leftovers)
		if err := tx.Where("id = ?", mergedID).Delete(&models.Account{}).Error; err != nil {
			return fmt.Errorf("deleting merged account: %w", err)
		}

		// Invalidate caches for both accounts
		if err := invalidateCacheByAccountID(survivorID); err != nil {
			log.Printf("[MERGE] Warning: failed to invalidate survivor cache: %v", err)
		}
		if err := invalidateCacheByAccountID(mergedID); err != nil {
			log.Printf("[MERGE] Warning: failed to invalidate merged cache: %v", err)
		}

		return nil
	})
}

// invalidateCacheByAccountID clears all cache entries for an account.
// If the account no longer exists (already deleted), only the ID-keyed entry is cleared.
func invalidateCacheByAccountID(id uuid.UUID) error {
	db := database.GetDB()
	var account models.Account
	if err := db.Preload("Identities").First(&account, "id = ?", id).Error; err != nil {
		// Account already deleted — clear only the ID-based key
		return database.DeleteCache(GetAccountCacheKeyByID(id.String()))
	}
	return InvalidateAccountCache(&account)
}
