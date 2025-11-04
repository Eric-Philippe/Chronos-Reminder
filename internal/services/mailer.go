package services

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

// MailerService handles email sending operations using Resend
type MailerService struct {
	client    *resend.Client
	fromEmail string
}

// EmailRequest represents the data needed to send an email
type EmailRequest struct {
	To       string
	Subject  string
	HtmlBody string
	TextBody string
}

// NewMailerService creates a new mailer service instance
func NewMailerService(apiKey string, fromEmail string) *MailerService {
	return &MailerService{
		client:    resend.NewClient(apiKey),
		fromEmail: fromEmail,
	}
}

// SendEmail sends an email using the Resend API
func (m *MailerService) SendEmail(req *EmailRequest) (string, error) {
	if req == nil {
		return "", fmt.Errorf("email request is nil")
	}

	if req.To == "" {
		return "", fmt.Errorf("recipient email is required")
	}

	if req.Subject == "" {
		return "", fmt.Errorf("email subject is required")
	}

	// Build the email request
	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{req.To},
		Subject: req.Subject,
	}

	// Set the body - prefer HTML, fall back to text
	if req.HtmlBody != "" {
		params.Html = req.HtmlBody
	} else if req.TextBody != "" {
		params.Text = req.TextBody
	} else {
		return "", fmt.Errorf("either HTML or text body is required")
	}

	// If both are provided, set text as well
	if req.TextBody != "" {
		params.Text = req.TextBody
	}

	// Send the email
	sent, err := m.client.Emails.Send(params)
	if err != nil {
		log.Printf("[MAILER] - ❌ Failed to send email to %s: %v", req.To, err)
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return sent.Id, nil
}

// SendEmailWithTemplate sends an email with custom template support
func (m *MailerService) SendEmailWithTemplate(to string, subject string, htmlBody string) (string, error) {
	req := &EmailRequest{
		To:       to,
		Subject:  subject,
		HtmlBody: htmlBody,
	}
	return m.SendEmail(req)
}

// SendWelcomeEmail sends a welcome email to a new user
func (m *MailerService) SendWelcomeEmail(email string, username string) (string, error) {
	subject := "Welcome to Chronos Reminder!"
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Welcome to Chronos Reminder</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h1 style="color: #4CAF50;">Welcome to Chronos Reminder! ⏰</h1>
		<p>Hi <strong>%s</strong>,</p>
		<p>Thank you for signing up! We're thrilled to have you on board.</p>
		<p>With Chronos Reminder, you can:</p>
		<ul>
			<li>Create and manage reminders effortlessly</li>
			<li>Get timely notifications via Discord</li>
			<li>Organize your schedule across multiple timezones</li>
		</ul>
		<p>Get started now and never miss an important moment!</p>
		<br>
		<p style="color: #666; font-size: 12px;">Best regards,<br>The Chronos Reminder Team</p>
	</div>
</body>
</html>
	`, username)

	textBody := fmt.Sprintf(`
Welcome to Chronos Reminder!

Hi %s,

Thank you for signing up! We're thrilled to have you on board.

With Chronos Reminder, you can:
- Create and manage reminders effortlessly
- Get timely notifications via Discord
- Organize your schedule across multiple timezones

Get started now and never miss an important moment!

Best regards,
The Chronos Reminder Team
	`, username)

	return m.SendEmail(&EmailRequest{
		To:       email,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	})
}

// SendPasswordResetEmail sends a password reset email
func (m *MailerService) SendPasswordResetEmail(email string, resetLink string) (string, error) {
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
		<h1 style="color: #FF9800;">Password Reset Request</h1>
		<p>We received a request to reset your password. Click the button below to proceed:</p>
		<p style="margin: 30px 0;">
			<a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">
				Reset Password
			</a>
		</p>
		<p style="color: #666; font-size: 12px;">If you didn't request this, please ignore this email. This link will expire in 1 hour.</p>
	</div>
</body>
</html>
	`, resetLink)

	textBody := fmt.Sprintf(`
Password Reset Request

We received a request to reset your password. Use this link to proceed:
%s

This link will expire in 1 hour.

If you didn't request this, please ignore this email.
	`, resetLink)

	return m.SendEmail(&EmailRequest{
		To:       email,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	})
}

// SendReminderNotificationEmail sends a reminder notification email
func (m *MailerService) SendReminderNotificationEmail(email string, reminderTitle string, reminderTime string) (string, error) {
	subject := fmt.Sprintf("Reminder: %s", reminderTitle)
	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Reminder Notification</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2196F3;">⏰ Reminder Notification</h2>
		<p style="font-size: 18px; margin: 20px 0;">
			<strong>%s</strong>
		</p>
		<p style="color: #666;">Scheduled for: <strong>%s</strong></p>
		<p style="margin-top: 30px; color: #999; font-size: 12px;">This is an automated reminder from Chronos Reminder</p>
	</div>
</body>
</html>
	`, reminderTitle, reminderTime)

	textBody := fmt.Sprintf(`
Reminder Notification

%s

Scheduled for: %s

This is an automated reminder from Chronos Reminder
	`, reminderTitle, reminderTime)

	return m.SendEmail(&EmailRequest{
		To:       email,
		Subject:  subject,
		HtmlBody: htmlBody,
		TextBody: textBody,
	})
}
