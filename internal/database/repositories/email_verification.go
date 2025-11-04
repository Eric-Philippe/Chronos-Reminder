package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailVerificationRepositoryImpl implements EmailVerificationRepository
type EmailVerificationRepositoryImpl struct {
	db *gorm.DB
}

// NewEmailVerificationRepository creates a new email verification repository
func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &EmailVerificationRepositoryImpl{db: db}
}

// Create creates a new email verification record
func (r *EmailVerificationRepositoryImpl) Create(verification *models.EmailVerification) error {
	if verification == nil {
		return fmt.Errorf("verification cannot be nil")
	}

	if err := r.db.Create(verification).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error creating email verification: %v", err)
		return fmt.Errorf("failed to create email verification: %w", err)
	}

	return nil
}

// GetByID retrieves an email verification by ID
func (r *EmailVerificationRepositoryImpl) GetByID(id uuid.UUID) (*models.EmailVerification, error) {
	var verification models.EmailVerification

	if err := r.db.Where("id = ?", id).First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email verification not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting email verification by ID: %v", err)
		return nil, err
	}

	return &verification, nil
}

// GetByEmailAndCode retrieves an email verification by email and code (regardless of verified status)
func (r *EmailVerificationRepositoryImpl) GetByEmailAndCode(email string, code string) (*models.EmailVerification, error) {
	var verification models.EmailVerification

	// Remove verified=false condition so we can check if code was already used
	if err := r.db.Where("email = ? AND code = ?", email, code).First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email verification not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting email verification: %v", err)
		return nil, err
	}

	return &verification, nil
}

// GetByEmail retrieves the latest email verification record for an email
func (r *EmailVerificationRepositoryImpl) GetByEmail(email string) (*models.EmailVerification, error) {
	var verification models.EmailVerification

	if err := r.db.Where("email = ?", email).Order("created_at DESC").First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email verification not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting email verification by email: %v", err)
		return nil, err
	}

	return &verification, nil
}

// GetByAccountID retrieves an email verification by account ID
func (r *EmailVerificationRepositoryImpl) GetByAccountID(accountID string) (*models.EmailVerification, error) {
	var verification models.EmailVerification

	if err := r.db.Where("account_id = ?", accountID).Order("created_at DESC").First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email verification not found")
		}
		log.Printf("[DATABASE] - ❌ Error getting email verification by account ID: %v", err)
		return nil, err
	}

	return &verification, nil
}

// MarkAsVerified marks an email verification as verified
func (r *EmailVerificationRepositoryImpl) MarkAsVerified(id uuid.UUID) error {
	now := time.Now()
	if err := r.db.Model(&models.EmailVerification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"verified":    true,
			"verified_at": now,
		}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error marking email verification as verified: %v", err)
		return fmt.Errorf("failed to mark email verification as verified: %w", err)
	}

	return nil
}

// IsVerified checks if an email has been verified
func (r *EmailVerificationRepositoryImpl) IsVerified(email string) (bool, error) {
	var verification models.EmailVerification

	if err := r.db.Where("email = ? AND verified = true", email).First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		log.Printf("[DATABASE] - ❌ Error checking email verification status: %v", err)
		return false, err
	}

	return true, nil
}

// DeleteByEmail deletes all email verification records for an email
func (r *EmailVerificationRepositoryImpl) DeleteByEmail(email string) error {
	if err := r.db.Where("email = ?", email).Delete(&models.EmailVerification{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting email verification by email: %v", err)
		return fmt.Errorf("failed to delete email verification: %w", err)
	}

	return nil
}

// Delete deletes an email verification record
func (r *EmailVerificationRepositoryImpl) Delete(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.EmailVerification{}).Error; err != nil {
		log.Printf("[DATABASE] - ❌ Error deleting email verification: %v", err)
		return fmt.Errorf("failed to delete email verification: %w", err)
	}

	return nil
}
