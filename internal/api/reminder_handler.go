package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
)

// ReminderHandler handles reminder-related HTTP requests
type ReminderHandler struct {
	reminderRepo    repositories.ReminderRepository
	destinationRepo repositories.ReminderDestinationRepository
	reminderErrorRepo repositories.ReminderErrorRepository
}

// NewReminderHandler creates a new reminder handler
func NewReminderHandler(
	reminderRepo repositories.ReminderRepository,
	destinationRepo repositories.ReminderDestinationRepository,
	reminderErrorRepo repositories.ReminderErrorRepository,
) *ReminderHandler {
	return &ReminderHandler{
		reminderRepo:      reminderRepo,
		destinationRepo:   destinationRepo,
		reminderErrorRepo: reminderErrorRepo,
	}
}

// GetReminder retrieves a single reminder by ID
func (h *ReminderHandler) GetReminder(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(AccountIDKey).(uuid.UUID)
	reminderID := r.PathValue("id")

	id, err := uuid.Parse(reminderID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	reminder, err := h.reminderRepo.GetWithAccountAndDestinations(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch reminder")
		return
	}

	if reminder == nil || reminder.AccountID != accountID {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	WriteJSON(w, http.StatusOK, ToReminderResponse(reminder))
}

// UpdateReminder updates an existing reminder
func (h *ReminderHandler) UpdateReminder(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(AccountIDKey).(uuid.UUID)
	reminderID := r.PathValue("id")

	id, err := uuid.Parse(reminderID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	// Fetch existing reminder
	reminder, err := h.reminderRepo.GetWithAccountAndDestinations(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch reminder")
		return
	}

	if reminder == nil || reminder.AccountID != accountID {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Parse request body
	var updateData struct {
		Message      string `json:"message"`
		Date         string `json:"date"`
		Time         string `json:"time"`
		Recurrence   int    `json:"recurrence"`
		Destinations []struct {
			Type     string                 `json:"type"`
			Metadata map[string]interface{} `json:"metadata"`
		} `json:"destinations"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update reminder fields
	if updateData.Message != "" {
		reminder.Message = updateData.Message
	}

	if updateData.Date != "" && updateData.Time != "" {
		remindAtUTC, err := time.Parse(time.RFC3339, updateData.Date+"T"+updateData.Time+":00Z")
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Invalid date/time format")
			return
		}
		reminder.RemindAtUTC = remindAtUTC
		reminder.NextFireUTC = &remindAtUTC
	}

	if updateData.Recurrence >= 0 {
		reminder.Recurrence = int16(updateData.Recurrence)
	}

	// Update destinations if provided
	if len(updateData.Destinations) > 0 {
		// Delete old destinations
		if err := h.destinationRepo.DeleteByReminderID(id); err != nil {
			WriteError(w, http.StatusInternalServerError, "Failed to update destinations")
			return
		}

		// Create new destinations
		newDestinations := make([]models.ReminderDestination, len(updateData.Destinations))
		for i, dest := range updateData.Destinations {
			destType := models.DestinationType(dest.Type)
			if !destType.IsValid() {
				WriteError(w, http.StatusBadRequest, "Invalid destination type")
				return
			}

			newDestinations[i] = models.ReminderDestination{
				ID:         uuid.New(),
				ReminderID: id,
				Type:       destType,
				Metadata:   models.JSONB(dest.Metadata),
			}
		}

		if err := h.destinationRepo.CreateMultiple(newDestinations); err != nil {
			WriteError(w, http.StatusInternalServerError, "Failed to create destinations")
			return
		}

		reminder.Destinations = newDestinations
	}

	// Save reminder
	if err := h.reminderRepo.Update(reminder, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update reminder")
		return
	}

	WriteJSON(w, http.StatusOK, ToReminderResponse(reminder))
}

// DeleteReminder deletes a reminder
func (h *ReminderHandler) DeleteReminder(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(AccountIDKey).(uuid.UUID)
	reminderID := r.PathValue("id")

	id, err := uuid.Parse(reminderID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	// Verify ownership
	reminder, err := h.reminderRepo.GetByID(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch reminder")
		return
	}

	if reminder == nil || reminder.AccountID != accountID {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Delete destinations first
	if err := h.destinationRepo.DeleteByReminderID(id); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete reminder")
		return
	}

	// Delete reminder
	if err := h.reminderRepo.Delete(id, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete reminder")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Reminder deleted successfully"})
}

// PauseReminder pauses a reminder
func (h *ReminderHandler) PauseReminder(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(AccountIDKey).(uuid.UUID)
	reminderID := r.PathValue("id")

	id, err := uuid.Parse(reminderID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	reminder, err := h.reminderRepo.GetByID(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch reminder")
		return
	}

	if reminder == nil || reminder.AccountID != accountID {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Set pause bit
	const pauseBit = 128
	reminder.Recurrence |= pauseBit

	if err := h.reminderRepo.Update(reminder, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to pause reminder")
		return
	}

	WriteJSON(w, http.StatusOK, ToReminderResponse(reminder))
}

// ResumeReminder resumes a paused reminder
func (h *ReminderHandler) ResumeReminder(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(AccountIDKey).(uuid.UUID)
	reminderID := r.PathValue("id")

	id, err := uuid.Parse(reminderID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	reminder, err := h.reminderRepo.GetByID(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch reminder")
		return
	}

	if reminder == nil || reminder.AccountID != accountID {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Clear pause bit
	const pauseBit = 128
	reminder.Recurrence &= ^pauseBit

	if err := h.reminderRepo.Update(reminder, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to resume reminder")
		return
	}

	WriteJSON(w, http.StatusOK, ToReminderResponse(reminder))
}

// DuplicateReminder duplicates a reminder
func (h *ReminderHandler) DuplicateReminder(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(AccountIDKey).(uuid.UUID)
	reminderID := r.PathValue("id")

	id, err := uuid.Parse(reminderID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	// Fetch original reminder with destinations
	original, err := h.reminderRepo.GetWithAccountAndDestinations(id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch reminder")
		return
	}

	if original == nil || original.AccountID != accountID {
		WriteError(w, http.StatusNotFound, "Reminder not found")
		return
	}

	// Create new reminder
	newReminder := &models.Reminder{
		ID:             uuid.New(),
		AccountID:      original.AccountID,
		RemindAtUTC:    original.RemindAtUTC,
		Message:        original.Message,
		Recurrence:     original.Recurrence,
		CreatedAt:      time.Now().UTC(),
		NextFireUTC:    original.NextFireUTC,
		SnoozedAtUTC:   original.SnoozedAtUTC,
	}

	// Create reminder
	if err := h.reminderRepo.Create(newReminder, true); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to duplicate reminder")
		return
	}

	// Duplicate destinations
	if len(original.Destinations) > 0 {
		newDestinations := make([]models.ReminderDestination, len(original.Destinations))
		for i, dest := range original.Destinations {
			newDestinations[i] = models.ReminderDestination{
				ID:        uuid.New(),
				ReminderID: newReminder.ID,
				Type:      dest.Type,
				Metadata:  dest.Metadata,
			}
		}

		if err := h.destinationRepo.CreateMultiple(newDestinations); err != nil {
			WriteError(w, http.StatusInternalServerError, "Failed to duplicate destinations")
			return
		}

		newReminder.Destinations = newDestinations
	}

	WriteJSON(w, http.StatusCreated, ToReminderResponse(newReminder))
}
