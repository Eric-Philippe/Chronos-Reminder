package repositories

import (
	"fmt"
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FcmTokenRepositoryImpl implements FcmTokenRepository
type FcmTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewFcmTokenRepository creates a new FCM token repository
func NewFcmTokenRepository(db *gorm.DB) FcmTokenRepository {
	return &FcmTokenRepositoryImpl{db: db}
}

// Upsert registers or updates the token for an (account, device) pair.
// A device keeps a single row; the token string itself is globally unique, so
// any stale row carrying the same token (e.g. after a reinstall on another
// account) is removed first to avoid a unique-constraint violation.
func (r *FcmTokenRepositoryImpl) Upsert(token *models.FcmToken) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Drop any other row holding this exact token value
		if err := tx.Where("token = ? AND NOT (account_id = ? AND device_id = ?)",
			token.Token, token.AccountID, token.DeviceID).
			Delete(&models.FcmToken{}).Error; err != nil {
			return err
		}

		var existing models.FcmToken
		err := tx.Where("account_id = ? AND device_id = ?", token.AccountID, token.DeviceID).
			First(&existing).Error

		switch {
		case err == gorm.ErrRecordNotFound:
			return tx.Create(token).Error
		case err != nil:
			return err
		default:
			return tx.Model(&existing).Update("token", token.Token).Error
		}
	})
}

// GetByAccountID returns all registered tokens for an account
func (r *FcmTokenRepositoryImpl) GetByAccountID(accountID uuid.UUID) ([]models.FcmToken, error) {
	var tokens []models.FcmToken
	if err := r.db.Where("account_id = ?", accountID).Find(&tokens).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error getting FCM tokens by account ID: %v", err)
		return nil, err
	}
	return tokens, nil
}

// DeleteByToken removes a token by its value (used for stale/unregistered tokens)
func (r *FcmTokenRepositoryImpl) DeleteByToken(token string) error {
	if err := r.db.Where("token = ?", token).Delete(&models.FcmToken{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting FCM token: %v", err)
		return fmt.Errorf("failed to delete FCM token: %w", err)
	}
	return nil
}

// DeleteByAccountAndDevice removes the token registered for a device (logout)
func (r *FcmTokenRepositoryImpl) DeleteByAccountAndDevice(accountID uuid.UUID, deviceID string) error {
	if err := r.db.Where("account_id = ? AND device_id = ?", accountID, deviceID).
		Delete(&models.FcmToken{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting FCM token by device: %v", err)
		return fmt.Errorf("failed to delete FCM token: %w", err)
	}
	return nil
}
