package services

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
)

// VerificationService handles email verification codes
type VerificationService struct {
	verificationRepo repositories.EmailVerificationRepository
	accountRepo      repositories.AccountRepository
	identityRepo     repositories.IdentityRepository
	mailerService    *MailerService
	webAppURL        string
	verificationTTL  time.Duration // Time-to-live for verification codes (default 24 hours)
}

// NewVerificationService creates a new verification service instance
func NewVerificationService(
	verificationRepo repositories.EmailVerificationRepository,
	accountRepo repositories.AccountRepository,
	identityRepo repositories.IdentityRepository,
	mailerService *MailerService,
	webAppURL string,
) *VerificationService {
	return &VerificationService{
		verificationRepo: verificationRepo,
		accountRepo:      accountRepo,
		identityRepo:     identityRepo,
		mailerService:    mailerService,
		webAppURL:        webAppURL,
		verificationTTL:  24 * time.Hour, // Codes expire after 24 hours
	}
}

// SendAccountVerification creates a verification record for an account and emails
// the code + link. Used at registration, Discord signup, and resend.
func (v *VerificationService) SendAccountVerification(email string, accountID string) error {
	code, err := v.CreateVerification(email, accountID)
	if err != nil {
		return err
	}
	verificationLink := v.webAppURL + "/verify?email=" + email + "&code=" + code
	if _, err := v.SendVerificationEmail(email, code, verificationLink); err != nil {
		return err
	}
	return nil
}

// ResendVerification re-sends a verification email for the given address. It is
// idempotent and intentionally quiet about whether the email maps to an account,
// to avoid leaking which addresses are registered. Returns true when an email
// was actually dispatched.
func (v *VerificationService) ResendVerification(email string) (bool, error) {
	accountID, ok := v.resolveAccountID(email)
	if !ok {
		// Unknown email — silently succeed to avoid enumeration.
		return false, nil
	}

	// Already verified? Nothing to do.
	account, err := v.accountRepo.GetByID(accountID)
	if err == nil && account != nil && account.EmailVerified {
		return false, nil
	}

	if err := v.SendAccountVerification(email, accountID.String()); err != nil {
		return false, err
	}
	return true, nil
}

// resolveAccountID finds the account owning an email, via the account-level
// email first and then any prior verification record (covers accounts created
// via Discord before credentials were set).
func (v *VerificationService) resolveAccountID(email string) (uuid.UUID, bool) {
	if account, err := v.accountRepo.GetByEmail(email); err == nil && account != nil {
		return account.ID, true
	}
	if record, err := v.verificationRepo.GetByEmail(email); err == nil && record != nil {
		if id, err := uuid.Parse(record.AccountID); err == nil {
			return id, true
		}
	}
	return uuid.Nil, false
}

// GenerateVerificationCode generates a 6-digit verification code
func (v *VerificationService) GenerateVerificationCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate verification code: %w", err)
		}
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code, nil
}

// SendVerificationEmail sends a verification code to the user's email
func (v *VerificationService) SendVerificationEmail(email string, verificationCode string, verificationLink string) (string, error) {
	subject := "Verify your Chronos Reminder account"
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Verify Your Email</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h1 style="color: #4CAF50;">Verify Your Email ✉️</h1>
		<p>Welcome to Chronos Reminder!</p>
		<p>Please verify your email address by clicking the button below:</p>
		
		<p style="margin: 30px 0;">
			<a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-size: 16px; font-weight: bold;">
				Verify Email
			</a>
		</p>
		
		<p style="color: #666; margin: 20px 0;">Or use this code: <strong style="font-size: 18px; letter-spacing: 2px;">%s</strong></p>
		
		<p style="color: #999; font-size: 12px; margin-top: 30px;">
			This link will expire in 24 hours. If you didn't create this account, please ignore this email.
		</p>
	</div>
</body>
</html>
	`, verificationLink, verificationCode)

	textBody := fmt.Sprintf(`
Verify Your Email

Welcome to Chronos Reminder!

Please verify your email address by visiting this link:
%s

Or use this verification code: %s

This link will expire in 24 hours. If you didn't create this account, please ignore this email.
	`, verificationLink, verificationCode)

	return v.mailerService.SendEmail(&EmailRequest{
		To:       email,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	})
}

// CreateVerification creates a new verification record
func (v *VerificationService) CreateVerification(email string, accountID string) (string, error) {
	// Generate verification code
	code, err := v.GenerateVerificationCode()
	if err != nil {
		return "", err
	}

	// Create verification record
	expiresAt := time.Now().Add(v.verificationTTL)
	verification := &models.EmailVerification{
		Email:     email,
		AccountID: accountID,
		Code:      code,
		ExpiresAt: expiresAt,
	}

	if err := v.verificationRepo.Create(verification); err != nil {
		return "", fmt.Errorf("failed to create verification record: %w", err)
	}

	return code, nil
}

// VerifyEmail verifies an email with the provided code
func (v *VerificationService) VerifyEmail(email string, code string) (string, error) {
	// Get verification record
	verification, err := v.verificationRepo.GetByEmailAndCode(email, code)
	if err != nil {
		return "", fmt.Errorf("invalid verification code or email")
	}

	// Check if already verified
	if verification.Verified {
		return "", fmt.Errorf("verification code has already been used")
	}

	// Check if code has expired
	if time.Now().After(verification.ExpiresAt) {
		return "", fmt.Errorf("verification code has expired")
	}

	// Mark verification record as verified
	if err := v.verificationRepo.MarkAsVerified(verification.ID); err != nil {
		return "", fmt.Errorf("failed to mark email as verified: %w", err)
	}

	// Parse account ID and mark account as email verified
	accountID, err := uuid.Parse(verification.AccountID)
	if err != nil {
		return "", fmt.Errorf("invalid account ID: %w", err)
	}

	// Get account and update EmailVerified flag
	account, err := v.accountRepo.GetByID(accountID)
	if err != nil {
		return "", fmt.Errorf("failed to get account: %w", err)
	}

	if account == nil {
		return "", fmt.Errorf("account not found")
	}

	account.EmailVerified = true
	if err := v.accountRepo.Update(account); err != nil {
		return "", fmt.Errorf("failed to update account verification status: %w", err)
	}

	return verification.AccountID, nil
}

// IsEmailVerified checks if an email has been verified
func (v *VerificationService) IsEmailVerified(email string) (bool, error) {
	verified, err := v.verificationRepo.IsVerified(email)
	if err != nil {
		return false, fmt.Errorf("failed to check verification status: %w", err)
	}
	return verified, nil
}

// DeleteVerification deletes verification records (after successful verification or cleanup)
func (v *VerificationService) DeleteVerification(email string) error {
	if err := v.verificationRepo.DeleteByEmail(email); err != nil {
		log.Printf("[VERIFICATION] - ⚠️ Failed to delete verification records for %s: %v", email, err)
		return err
	}
	return nil
}
