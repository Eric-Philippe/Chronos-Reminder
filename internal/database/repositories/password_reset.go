package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordResetRepositoryImpl implements PasswordResetRepository
type PasswordResetRepositoryImpl struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &PasswordResetRepositoryImpl{db: db}
}

// Create creates a new password reset record
func (r *PasswordResetRepositoryImpl) Create(reset *models.PasswordReset) error {
	if reset == nil {
		return fmt.Errorf("password reset cannot be nil")
	}

	if err := r.db.Create(reset).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error creating password reset: %v", err)
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	return nil
}

// GetByID retrieves a password reset by ID
func (r *PasswordResetRepositoryImpl) GetByID(id uuid.UUID) (*models.PasswordReset, error) {
	var reset models.PasswordReset

	if err := r.db.Where("id = ?", id).First(&reset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("password reset not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting password reset by ID: %v", err)
		return nil, err
	}

	return &reset, nil
}

// GetByToken retrieves a password reset by token
func (r *PasswordResetRepositoryImpl) GetByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset

	if err := r.db.Where("token = ? AND used = false", token).First(&reset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("password reset token not found or already used")
		}
		log.Printf("[DATABASE] - ❌ Error getting password reset by token: %v", err)
		return nil, err
	}

	return &reset, nil
}

// GetByEmail retrieves the latest unused password reset for an email
func (r *PasswordResetRepositoryImpl) GetByEmail(email string) (*models.PasswordReset, error) {
	var reset models.PasswordReset

	if err := r.db.Where("email = ? AND used = false", email).Order("created_at DESC").First(&reset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("password reset not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting password reset by email: %v", err)
		return nil, err
	}

	return &reset, nil
}

// GetByAccountID retrieves the latest unused password reset for an account
func (r *PasswordResetRepositoryImpl) GetByAccountID(accountID uuid.UUID) (*models.PasswordReset, error) {
	var reset models.PasswordReset

	if err := r.db.Where("account_id = ? AND used = false", accountID).Order("created_at DESC").First(&reset).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("password reset not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting password reset by account ID: %v", err)
		return nil, err
	}

	return &reset, nil
}

// MarkAsUsed marks a password reset as used
func (r *PasswordResetRepositoryImpl) MarkAsUsed(id uuid.UUID) error {
	now := time.Now()
	if err := r.db.Model(&models.PasswordReset{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": now,
		}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error marking password reset as used: %v", err)
		return fmt.Errorf("failed to mark password reset as used: %w", err)
	}

	return nil
}

// DeleteByEmail deletes all password reset records for an email
func (r *PasswordResetRepositoryImpl) DeleteByEmail(email string) error {
	if err := r.db.Where("email = ?", email).Delete(&models.PasswordReset{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting password reset by email: %v", err)
		return fmt.Errorf("failed to delete password reset: %w", err)
	}

	return nil
}

// Delete deletes a password reset record
func (r *PasswordResetRepositoryImpl) Delete(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.PasswordReset{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting password reset: %v", err)
		return fmt.Errorf("failed to delete password reset: %w", err)
	}

	return nil
}

// DeleteExpiredTokens deletes all expired password reset tokens
func (r *PasswordResetRepositoryImpl) DeleteExpiredTokens() error {
	if err := r.db.Where("expires_at < ? AND used = false", time.Now()).Delete(&models.PasswordReset{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting expired password reset tokens: %v", err)
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}
