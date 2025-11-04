package services

// Example usage of the MailerService
// This file demonstrates how to use the mailer service in your API handlers

/*
BASIC USAGE EXAMPLES:

1. Send a simple email:
   mailer := services.NewMailerService(apiKey, "noreply@chronosrmd.com")

   req := &services.EmailRequest{
       To: "user@example.com",
       Subject: "Test Email",
       HtmlBody: "<h1>Hello!</h1><p>This is a test email.</p>",
   }

   emailID, err := mailer.SendEmail(req)
   if err != nil {
       log.Printf("Error sending email: %v", err)
   } else {
       log.Printf("Email sent with ID: %s", emailID)
   }

2. Send a welcome email:
   emailID, err := mailer.SendWelcomeEmail("user@example.com", "John")
   if err != nil {
       log.Printf("Error sending welcome email: %v", err)
   }

3. Send a password reset email:
   emailID, err := mailer.SendPasswordResetEmail(
       "user@example.com",
       "https://yourapp.com/reset?token=abc123",
   )
   if err != nil {
       log.Printf("Error sending reset email: %v", err)
   }

4. Send a reminder notification:
   emailID, err := mailer.SendReminderNotificationEmail(
       "user@example.com",
       "Team Meeting",
       "2024-11-04 14:30 UTC",
   )
   if err != nil {
       log.Printf("Error sending reminder: %v", err)
   }

5. Send a custom template email:
   customHTML := `
   <h1>Custom Email</h1>
   <p>Some custom content here</p>
   `
   emailID, err := mailer.SendEmailWithTemplate(
       "user@example.com",
       "Custom Subject",
       customHTML,
   )
   if err != nil {
       log.Printf("Error sending custom email: %v", err)
   }

USAGE IN API HANDLERS:

In your auth handler:
   func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
       // ... registration logic ...

       // Send welcome email
       mailerService := h.server.GetMailerService() // Get from server instance
       _, err := mailerService.SendWelcomeEmail(account.Email, account.Username)
       if err != nil {
           log.Printf("Failed to send welcome email: %v", err)
           // Don't fail the registration, just log the error
       }

       // ... return response ...
   }

In your user handler:
   func (h *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
       // ... validation logic ...

       // Generate reset token and send email
       resetToken := generateToken() // Your token generation logic
       resetLink := fmt.Sprintf("https://yourapp.com/reset?token=%s", resetToken)

       _, err := h.mailerService.SendPasswordResetEmail(email, resetLink)
       if err != nil {
           http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
           return
       }

       w.WriteHeader(http.StatusOK)
       json.NewEncoder(w).Encode(map[string]string{"message": "Reset email sent"})
   }

ENVIRONMENT CONFIGURATION:

Make sure to add to your .env file:
   RESEND_API_KEY="your_resend_api_key_here"

Get your API key from: https://resend.com/api-keys

FEATURES:

- ✅ Send simple text and HTML emails
- ✅ Pre-built templates for common emails (welcome, password reset, reminders)
- ✅ Error handling and logging
- ✅ Support for both HTML and plain text bodies
- ✅ Resend API integration with email ID tracking

RESEND API DOCUMENTATION:

For more details on Resend API:
- GitHub: github.com/resend/resend-go/v2
- Documentation: resend.com/docs
- API Reference: resend.com/docs/api-reference

ERROR HANDLING:

Always handle errors gracefully:
   emailID, err := mailerService.SendEmail(req)
   if err != nil {
       // Log the error but don't fail the user's request
       log.Printf("[MAILER] - ❌ Failed to send email: %v", err)
       // Optionally, you could store the email in a queue for retry
       return
   }
*/
