package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
)

// PasswordResetService handles password reset operations
type PasswordResetService struct {
	passwordResetRepo repositories.PasswordResetRepository
	identityRepo      repositories.IdentityRepository
	mailerService     *MailerService
	resetTokenTTL     time.Duration // Time-to-live for reset tokens (default 24 hours)
}

// NewPasswordResetService creates a new password reset service instance
func NewPasswordResetService(
	passwordResetRepo repositories.PasswordResetRepository,
	identityRepo repositories.IdentityRepository,
	mailerService *MailerService,
) *PasswordResetService {
	return &PasswordResetService{
		passwordResetRepo: passwordResetRepo,
		identityRepo:      identityRepo,
		mailerService:     mailerService,
		resetTokenTTL:     24 * time.Hour, // Tokens expire after 24 hours
	}
}

// GenerateResetToken generates a cryptographically secure reset token
func (p *PasswordResetService) GenerateResetToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}
	return hex.EncodeToString(token), nil
}

// RequestPasswordReset creates a password reset token and sends an email to the user
func (p *PasswordResetService) RequestPasswordReset(email string) error {
	// Check if identity exists with this email
	identity, err := p.identityRepo.GetByProviderAndExternalID(models.ProviderApp, email)
	if err != nil {
		return fmt.Errorf("email not found")
	}

	if identity == nil {
		// Return generic message for security (don't reveal if email exists)
		return nil
	}

	// Generate reset token
	token, err := p.GenerateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset record
	expiresAt := time.Now().Add(p.resetTokenTTL)
	passwordReset := &models.PasswordReset{
		AccountID: identity.AccountID,
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := p.passwordResetRepo.Create(passwordReset); err != nil {
		return fmt.Errorf("failed to create password reset record: %w", err)
	}

	// Build reset link (frontend will handle the redirect)
	resetLink := fmt.Sprintf("https://chronosrmd.com/reset-password?token=%s&email=%s", token, email)

	// Send reset email
	_, err = p.SendPasswordResetEmail(email, resetLink)
	if err != nil {
		log.Printf("[PASSWORD_RESET] - ⚠️ Failed to send password reset email: %v", err)
		// Don't fail the request even if email fails - token is still created
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email to the user
func (p *PasswordResetService) SendPasswordResetEmail(email string, resetLink string) (string, error) {
	subject := "Reset your Chronos Reminder password"
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Reset Your Password</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h1 style="color: #4CAF50;">Reset Your Password 🔐</h1>
		<p>We received a request to reset your Chronos Reminder account password.</p>
		<p>Click the button below to reset your password:</p>
		
		<p style="margin: 30px 0;">
			<a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block; font-size: 16px; font-weight: bold;">
				Reset Password
			</a>
		</p>
		
		<p style="color: #666; margin: 20px 0;">Or copy this link: <br><span style="word-break: break-all;">%s</span></p>
		
		<p style="color: #999; font-size: 12px; margin-top: 30px;">
			This link will expire in 24 hours. If you didn't request a password reset, please ignore this email or contact our support team.
		</p>
	</div>
</body>
</html>
	`, resetLink, resetLink)

	textBody := fmt.Sprintf(`
Reset Your Password

We received a request to reset your Chronos Reminder account password.

Please visit this link to reset your password:
%s

This link will expire in 24 hours. If you didn't request a password reset, please ignore this email.
	`, resetLink)

	return p.mailerService.SendEmail(&EmailRequest{
		To:       email,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	})
}

// VerifyResetToken verifies a password reset token
func (p *PasswordResetService) VerifyResetToken(email string, token string) (*models.PasswordReset, error) {
	// Get password reset record
	passwordReset, err := p.passwordResetRepo.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired reset token")
	}

	// Check if token belongs to the email
	if passwordReset.Email != email {
		return nil, fmt.Errorf("token does not match email")
	}

	// Check if token has been used
	if passwordReset.Used {
		return nil, fmt.Errorf("reset token has already been used")
	}

	// Check if token has expired
	if time.Now().After(passwordReset.ExpiresAt) {
		return nil, fmt.Errorf("reset token has expired")
	}

	return passwordReset, nil
}

// ResetPassword resets the user's password using a valid reset token
func (p *PasswordResetService) ResetPassword(email string, token string, newPassword string) error {
	// Verify token first
	passwordReset, err := p.VerifyResetToken(email, token)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Get identity for app provider
	identity, err := p.identityRepo.GetByProviderAndExternalID(models.ProviderApp, email)
	if err != nil {
		return fmt.Errorf("identity not found: %w", err)
	}

	if identity == nil {
		return fmt.Errorf("identity not found")
	}

	// Update password
	identity.PasswordHash = &hashedPassword
	if err := p.identityRepo.Update(identity); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark reset token as used
	if err := p.passwordResetRepo.MarkAsUsed(passwordReset.ID); err != nil {
		log.Printf("[PASSWORD_RESET] - ⚠️ Failed to mark reset token as used: %v", err)
	}

	// Delete all other unused reset tokens for this email for security
	if err := p.passwordResetRepo.DeleteByEmail(email); err != nil {
		log.Printf("[PASSWORD_RESET] - ⚠️ Failed to clean up old reset tokens: %v", err)
	}

	return nil
}

// IsResetTokenValid checks if a reset token is valid
func (p *PasswordResetService) IsResetTokenValid(email string, token string) bool {
	_, err := p.VerifyResetToken(email, token)
	return err == nil
}

// DeleteResetToken deletes a reset token
func (p *PasswordResetService) DeleteResetToken(email string) error {
	if err := p.passwordResetRepo.DeleteByEmail(email); err != nil {
		log.Printf("[PASSWORD_RESET] - ⚠️ Failed to delete reset tokens for %s: %v", email, err)
		return err
	}
	return nil
}

// CleanupExpiredTokens deletes all expired reset tokens
func (p *PasswordResetService) CleanupExpiredTokens() error {
	if err := p.passwordResetRepo.DeleteExpiredTokens(); err != nil {
		log.Printf("[PASSWORD_RESET] - ⚠️ Failed to cleanup expired tokens: %v", err)
		return err
	}
	return nil
}
