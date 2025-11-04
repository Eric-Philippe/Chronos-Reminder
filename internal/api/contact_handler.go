package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// ContactRequest represents a contact form submission
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// ContactHandler handles contact form submissions
type ContactHandler struct {
	mailerService *services.MailerService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(mailerService *services.MailerService) *ContactHandler {
	return &ContactHandler{
		mailerService: mailerService,
	}
}

// SubmitContact handles POST /api/contact requests
// @Summary Submit a contact form
// @Description Submit a contact form with name, email, type, subject, and message
// @Tags Contact
// @Accept json
// @Produce json
// @Param request body ContactRequest true "Contact request"
// @Success 200 {object} map[string]interface{} "Message sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/contact [post]
func (h *ContactHandler) SubmitContact(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if err := validateContactRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Sanitize inputs
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Type = strings.TrimSpace(req.Type)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Message = strings.TrimSpace(req.Message)

	// Build email HTML body
	htmlBody := buildContactEmailHTML(&req)

	// Build email text body
	textBody := fmt.Sprintf(`Contact Form Submission

From: %s <%s>
Type: %s
Subject: %s

Message:
%s

---
This message was sent through the Chronos Reminder contact form.
`, req.Name, req.Email, req.Type, req.Subject, req.Message)

	// Send email to admin/support
	_, err := h.mailerService.SendEmail(&services.EmailRequest{
		To:       "ericphlpp@proton.me",
		Subject:  fmt.Sprintf("[%s] %s - from %s", strings.ToUpper(req.Type), req.Subject, req.Name),
		HtmlBody: htmlBody,
		TextBody: textBody,
	})

	if err != nil {
		log.Printf("[CONTACT] - ❌ Failed to send contact email: %v", err)
		WriteError(w, http.StatusInternalServerError, "Failed to send message. Please try again later.")
		return
	}

	// Optionally send confirmation email to the user
	_, err = h.mailerService.SendEmailWithTemplate(
		req.Email,
		"We received your message",
		buildConfirmationEmailHTML(req.Name),
	)

	if err != nil {
		// Log the error but don't fail the request
		log.Printf("[CONTACT] - ⚠️  Failed to send confirmation email to %s: %v", req.Email, err)
	}

	log.Printf("[CONTACT] - ✅ Contact form submitted successfully from %s (%s)", req.Email, req.Type)

	// Return success response
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Message sent successfully",
		"status":  "success",
	})
}

// validateContactRequest validates the contact request
func validateContactRequest(req *ContactRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email validation
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email address")
	}

	if strings.TrimSpace(req.Subject) == "" {
		return fmt.Errorf("subject is required")
	}

	if strings.TrimSpace(req.Message) == "" {
		return fmt.Errorf("message is required")
	}

	// Validate type
	validTypes := []string{"general", "feedback", "bug", "feature"}
	isValidType := false
	for _, t := range validTypes {
		if req.Type == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid message type")
	}

	return nil
}

// buildContactEmailHTML builds the HTML email body for the contact form submission
func buildContactEmailHTML(req *ContactRequest) string {
	typeEmoji := map[string]string{
		"general":  "📝",
		"feedback": "💬",
		"bug":      "🐛",
		"feature":  "✨",
	}

	emoji := typeEmoji[req.Type]
	if emoji == "" {
		emoji = "📧"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>New Contact Form Submission</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h1 style="color: #4CAF50;">%s New Contact Form Submission</h1>
		
		<div style="background-color: #f5f5f5; padding: 20px; border-radius: 5px; margin: 20px 0;">
			<h2 style="margin-top: 0;">Contact Details</h2>
			<p><strong>Name:</strong> %s</p>
			<p><strong>Email:</strong> <a href="mailto:%s">%s</a></p>
			<p><strong>Type:</strong> <span style="background-color: #e3f2fd; padding: 2px 8px; border-radius: 3px;">%s</span></p>
			<p><strong>Subject:</strong> %s</p>
		</div>

		<div style="background-color: #fafafa; padding: 20px; border-left: 4px solid #4CAF50; margin: 20px 0;">
			<h3 style="margin-top: 0;">Message</h3>
			<p style="white-space: pre-wrap; word-wrap: break-word;">%s</p>
		</div>

		<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
		<p style="color: #999; font-size: 12px;">
			This message was submitted through the Chronos Reminder contact form.
		</p>
	</div>
</body>
</html>
	`, emoji, req.Name, req.Email, req.Email, req.Type, req.Subject, req.Message)
}

// buildConfirmationEmailHTML builds the HTML confirmation email for the user
func buildConfirmationEmailHTML(name string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>We Received Your Message</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h1 style="color: #4CAF50;">✅ We Received Your Message</h1>
		
		<p>Hi <strong>%s</strong>,</p>
		
		<p>Thank you for contacting us! We've received your message and appreciate you taking the time to reach out.</p>
		
		<p>Our team will review your submission and get back to you as soon as possible. We typically respond within 24-48 hours.</p>
		
		<div style="background-color: #f5f5f5; padding: 20px; border-radius: 5px; margin: 20px 0;">
			<h3 style="margin-top: 0;">What to expect next</h3>
			<ul>
				<li>We'll review your message carefully</li>
				<li>Our team will respond via email</li>
				<li>For urgent matters, you can reach us on Discord</li>
			</ul>
		</div>

		<p style="color: #666; font-size: 12px; margin-top: 30px;">
			Best regards,<br>
			The Chronos Reminder Team
		</p>
	</div>
</body>
</html>
	`, name)
}
