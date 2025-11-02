package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// CreateReminderRequest represents the request body for creating a reminder
type CreateReminderRequest struct {
	Date         string                       `json:"date"`         // ISO 8601 date format
	Time         string                       `json:"time"`         // HH:mm format
	Message      string                       `json:"message"`
	Recurrence   int16                        `json:"recurrence"`
	Destinations []CreateDestinationRequest   `json:"destinations"`
}

// CreateDestinationRequest represents a destination to create
type CreateDestinationRequest struct {
	Type     string                 `json:"type"` // "discord_dm", "discord_channel", "webhook"
	Metadata map[string]interface{} `json:"metadata"`
}

// CreateReminderResponse represents the response after creating a reminder
type CreateReminderResponse struct {
	ID              uuid.UUID         `json:"id"`
	Message         string            `json:"message"`
	RemindAtUTC     time.Time         `json:"remind_at_utc"`
	RecurrenceType  int               `json:"recurrence_type"`
	IsPaused        bool              `json:"is_paused"`
	Destinations    []interface{}     `json:"destinations"`
}

// ReminderResponse represents a reminder in API responses with decoded recurrence
type ReminderResponse struct {
	ID              uuid.UUID              `json:"id"`
	AccountID       uuid.UUID              `json:"account_id"`
	RemindAtUTC     time.Time              `json:"remind_at_utc"`
	SnoozedAtUTC    *time.Time             `json:"snoozed_at_utc,omitempty"`
	NextFireUTC     *time.Time             `json:"next_fire_utc,omitempty"`
	Message         string                 `json:"message"`
	CreatedAt       time.Time              `json:"created_at"`
	RecurrenceType  int                    `json:"recurrence_type"`
	IsPaused        bool                   `json:"is_paused"`
	Destinations    []models.ReminderDestination `json:"destinations,omitempty"`
}

// ToReminderResponse converts a Reminder model to ReminderResponse with decoded recurrence
func ToReminderResponse(reminder *models.Reminder) *ReminderResponse {
	return &ReminderResponse{
		ID:             reminder.ID,
		AccountID:      reminder.AccountID,
		RemindAtUTC:    reminder.RemindAtUTC,
		SnoozedAtUTC:   reminder.SnoozedAtUTC,
		NextFireUTC:    reminder.NextFireUTC,
		Message:        reminder.Message,
		CreatedAt:      reminder.CreatedAt,
		RecurrenceType: services.GetRecurrenceType(int(reminder.Recurrence)),
		IsPaused:       services.IsPaused(int(reminder.Recurrence)),
		Destinations:   reminder.Destinations,
	}
}

// UserHandler handles user-related requests
type UserHandler struct {
	reminderRepo            repositories.ReminderRepository
	reminderErrorRepo       repositories.ReminderErrorRepository
	reminderDestinationRepo repositories.ReminderDestinationRepository
	accountRepo             repositories.AccountRepository
	sessionService          *services.SessionService
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	reminderRepo repositories.ReminderRepository,
	reminderErrorRepo repositories.ReminderErrorRepository,
	accountRepo repositories.AccountRepository,
	sessionService *services.SessionService,
) *UserHandler {
	return &UserHandler{
		reminderRepo:      reminderRepo,
		reminderErrorRepo: reminderErrorRepo,
		accountRepo:       accountRepo,
		sessionService:    sessionService,
	}
}

// SetReminderDestinationRepository sets the reminder destination repository
func (h *UserHandler) SetReminderDestinationRepository(repo repositories.ReminderDestinationRepository) {
	h.reminderDestinationRepo = repo
}

// CreateReminder creates a new reminder for the authenticated user
// @Route: POST /api/reminders
func (h *UserHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract account ID from token
	accountID, err := h.extractAccountIDFromToken(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Get account with timezone to access timezone information
	account, err := h.accountRepo.GetWithTimezone(accountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve account")
		return
	}

	if account == nil {
		WriteError(w, http.StatusNotFound, "Account not found")
		return
	}

	if account.Timezone == nil {
		WriteError(w, http.StatusBadRequest, "Account timezone not set")
		return
	}

	// Parse request body
	var req CreateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if req.Message == "" {
		WriteError(w, http.StatusBadRequest, "Message is required")
		return
	}

	if req.Date == "" || req.Time == "" {
		WriteError(w, http.StatusBadRequest, "Date and time are required")
		return
	}

	// Parse the reminder date and time in user's timezone
	location, err := time.LoadLocation(account.Timezone.IANALocation)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timezone: %s", account.Timezone.IANALocation))
		return
	}

	parsedTime, err := services.ParseReminderDateTimeInTimezone(req.Date, req.Time, account.Timezone.IANALocation)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid date/time format")
		return
	}

	// Check if the reminder time is in the future
	now := time.Now().In(location)
	if parsedTime.Before(now) {
		WriteError(w, http.StatusBadRequest, "Reminder date/time must be in the future")
		return
	}

	// Create the reminder with UTC time
	reminder := &models.Reminder{
		AccountID:   accountID,
		RemindAtUTC: parsedTime.UTC(),
		Message:     req.Message,
		Recurrence:  req.Recurrence,
	}

	// Save the reminder to database
	if err := h.reminderRepo.Create(reminder, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create reminder")
		return
	}

	// Process destinations
	var destinations []interface{}
	for _, dest := range req.Destinations {
		destType := models.DestinationType(dest.Type)

		// Validate destination type
		if !destType.IsValid() {
			continue
		}

		// Handle discord_dm destination
		if destType == models.DestinationDiscordDM {
			// If user_id not provided, get it from the Discord identity
			if _, exists := dest.Metadata["user_id"]; !exists {
				account, err := h.accountRepo.GetWithIdentities(accountID)
				if err == nil && account != nil {
					for _, identity := range account.Identities {
						if identity.Provider == "discord" {
							dest.Metadata["user_id"] = identity.ExternalID
							break
						}
					}
				}
				// Skip if still no user_id
				if _, exists := dest.Metadata["user_id"]; !exists {
					continue
				}
			}
		}

		// Handle discord_channel destination
		if destType == models.DestinationDiscordChannel {
			// Validate required fields
			if _, hasGuild := dest.Metadata["guild_id"]; !hasGuild {
				continue
			}
			if _, hasChannel := dest.Metadata["channel_id"]; !hasChannel {
				continue
			}
			// mention_role_id is optional
		}

		// Handle webhook destination
		if destType == models.DestinationWebhook {
			// Validate required fields
			if _, hasURL := dest.Metadata["url"]; !hasURL {
				continue
			}
			
			// Validate optional platform field
			if platformVal, exists := dest.Metadata["platform"]; exists {
				if platformStr, ok := platformVal.(string); ok {
					platform := models.WebhookPlatform(platformStr)
					if !platform.IsValid() {
						// Invalid platform, skip this destination
						continue
					}
				}
			}
		}

		// Create the destination
		reminderDest := &models.ReminderDestination{
			ReminderID: reminder.ID,
			Type:       destType,
			Metadata:   dest.Metadata,
		}

		if err := h.reminderDestinationRepo.Create(reminderDest); err != nil {
			// Log error but continue - don't fail the entire operation
			fmt.Printf("[CREATE_REMINDER] Failed to create destination: %v\n", err)
			continue
		}

		destinations = append(destinations, reminderDest)
	}

	// If no destinations were provided or valid, create a default discord_dm destination
	// Get the Discord identity for the user
	if len(destinations) == 0 {
		account, err := h.accountRepo.GetWithIdentities(accountID)
		if err == nil && account != nil && len(account.Identities) > 0 {
			// Find Discord identity
			for _, identity := range account.Identities {
				if identity.Provider == "discord" {
					reminderDest := &models.ReminderDestination{
						ReminderID: reminder.ID,
						Type:       models.DestinationDiscordDM,
						Metadata: models.JSONB{
							"user_id": identity.ExternalID,
						},
					}

					if err := h.reminderDestinationRepo.Create(reminderDest); err == nil {
						destinations = append(destinations, reminderDest)
					}
					break
				}
			}
		}
	}

	// Build response with decoded recurrence
	recurrenceType := services.GetRecurrenceType(int(reminder.Recurrence))
	isPaused := services.IsPaused(int(reminder.Recurrence))

	response := CreateReminderResponse{
		ID:             reminder.ID,
		Message:        reminder.Message,
		RemindAtUTC:    reminder.RemindAtUTC,
		RecurrenceType: recurrenceType,
		IsPaused:       isPaused,
		Destinations:   destinations,
	}

	WriteJSON(w, http.StatusCreated, response)
}

// GetReminders retrieves all reminders for the authenticated user with their destinations
// @Route: GET /api/reminders
func (h *UserHandler) GetReminders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accountID, err := h.extractAccountIDFromToken(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	reminders, err := h.reminderRepo.GetByAccountIDWithDestinations(accountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve reminders")
		return
	}

	// Convert reminders to response format with decoded recurrence
	reminderResponses := make([]*ReminderResponse, len(reminders))
	for i, reminder := range reminders {
		reminderResponses[i] = ToReminderResponse(&reminder)
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"reminders": reminderResponses,
		"count":     len(reminderResponses),
	})
}

// GetReminder retrieves a single reminder by ID for the authenticated user
// @Route: GET /api/reminders/{id}
func (h *UserHandler) GetReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accountID, err := h.extractAccountIDFromToken(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	reminderIDStr := strings.TrimPrefix(r.URL.Path, "/api/reminders/")
	reminderID, err := uuid.Parse(reminderIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	reminder, err := h.reminderRepo.GetWithAccountAndDestinations(reminderID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve reminder")
		return
	}

	if reminder == nil {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Verify ownership
	if reminder.AccountID != accountID {
		WriteError(w, http.StatusForbidden, "You do not have permission to access this reminder")
		return
	}

	WriteJSON(w, http.StatusOK, ToReminderResponse(reminder))
}

// GetReminderErrors retrieves all reminders with errors for the authenticated user
// @Route: GET /api/reminders/errors
func (h *UserHandler) GetReminderErrors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accountID, err := h.extractAccountIDFromToken(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Get all reminders for the user
	reminders, err := h.reminderRepo.GetByAccountIDWithDestinations(accountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve reminders")
		return
	}

	// Collect all reminder IDs
	reminderIDs := make([]uuid.UUID, 0, len(reminders))
	for _, reminder := range reminders {
		reminderIDs = append(reminderIDs, reminder.ID)
	}

	// Get all errors for these reminders
	reminderErrors := make([]models.ReminderError, 0)
	for _, reminderID := range reminderIDs {
		errs, err := h.reminderErrorRepo.GetByReminderID(reminderID)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "Failed to retrieve reminder errors")
			return
		}
		// Append each error individually and handle both value and pointer element types
		for _, e := range errs {
			switch v := interface{}(e).(type) {
			case models.ReminderError:
				reminderErrors = append(reminderErrors, v)
			case *models.ReminderError:
				if v != nil {
					reminderErrors = append(reminderErrors, *v)
				}
			default:
				// ignore unexpected types
			}
		}
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"errors": reminderErrors,
		"count":  len(reminderErrors),
	})
}

// GetAccount retrieves the authenticated user's account information with identities
// @Route: GET /api/account
func (h *UserHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accountID, err := h.extractAccountIDFromToken(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	account, err := h.accountRepo.GetWithIdentities(accountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve account")
		return
	}

	if account == nil {
		WriteError(w, http.StatusNotFound, "Account not found")
		return
	}

	WriteJSON(w, http.StatusOK, account)
}

// DeleteReminder deletes a reminder for the authenticated user
// @Route: DELETE /api/reminders/{id}
func (h *UserHandler) DeleteReminder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	accountID, err := h.extractAccountIDFromToken(r)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	reminderIDStr := strings.TrimPrefix(r.URL.Path, "/api/reminders/")
	reminderID, err := uuid.Parse(reminderIDStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	// Get the reminder to verify ownership
	reminder, err := h.reminderRepo.GetByID(reminderID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve reminder")
		return
	}

	if reminder == nil {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Verify ownership
	if reminder.AccountID != accountID {
		WriteError(w, http.StatusForbidden, "You do not have permission to delete this reminder")
		return
	}

	// Delete the reminder with notification enabled
	if err := h.reminderRepo.Delete(reminderID, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete reminder")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Reminder deleted successfully",
	})
}

// extractAccountIDFromToken extracts the account ID from the JWT token in the request
func (h *UserHandler) extractAccountIDFromToken(r *http.Request) (uuid.UUID, error) {
	// Try to get token from Authorization header first
	authHeader := r.Header.Get("Authorization")
	var token string

	if authHeader != "" {
		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			token = parts[1]
		}
	}

	// If no token in header, try cookie
	if token == "" {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			return uuid.Nil, ErrUnauthorized
		}
		token = cookie.Value
	}

	if token == "" {
		return uuid.Nil, ErrUnauthorized
	}

	// Validate token and extract claims
	claims, err := h.sessionService.ValidateToken(token)
	if err != nil {
		return uuid.Nil, ErrUnauthorized
	}

	// Parse account ID from claims
	accountID, err := uuid.Parse(claims.AccountID)
	if err != nil {
		return uuid.Nil, ErrUnauthorized
	}

	return accountID, nil
}

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(sessionService *services.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get token from Authorization header first
			authHeader := r.Header.Get("Authorization")
			var token string

			if authHeader != "" {
				// Extract Bearer token
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					token = parts[1]
				}
			}

			// If no token in header, try cookie
			if token == "" {
				cookie, err := r.Cookie("auth_token")
				if err != nil {
					WriteError(w, http.StatusUnauthorized, "No authentication token found")
					return
				}
				token = cookie.Value
			}

			if token == "" {
				WriteError(w, http.StatusUnauthorized, "No authentication token found")
				return
			}

			// Validate token and extract claims
			claims, err := sessionService.ValidateToken(token)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Extract account ID from claims
			accountID, err := uuid.Parse(claims.AccountID)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			// Add account ID to request context
			ctx := context.WithValue(r.Context(), AccountIDKey, accountID)
			*r = *r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
