package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
)

// FcmHandler handles Firebase Cloud Messaging token registration
type FcmHandler struct {
	fcmTokenRepo repositories.FcmTokenRepository
}

// NewFcmHandler creates a new FCM handler
func NewFcmHandler(fcmTokenRepo repositories.FcmTokenRepository) *FcmHandler {
	return &FcmHandler{fcmTokenRepo: fcmTokenRepo}
}

// fcmTokenRequest is the body for both register and unregister requests
type fcmTokenRequest struct {
	Token    string `json:"token"`
	DeviceID string `json:"device_id"`
}

func (h *FcmHandler) accountID(r *http.Request) (uuid.UUID, bool) {
	val := r.Context().Value(AccountIDKey)
	if val == nil {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// HasTokens returns whether the calling account has at least one registered FCM token.
// @Route: GET /api/fcm/status
func (h *FcmHandler) HasTokens(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.accountID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Account ID not found in context")
		return
	}

	tokens, err := h.fcmTokenRepo.GetByAccountID(accountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to check FCM tokens")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]bool{"has_push": len(tokens) > 0})
}

// RegisterToken registers or updates the FCM token for the calling account.
// @Route: POST /api/fcm/token
func (h *FcmHandler) RegisterToken(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.accountID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Account ID not found in context")
		return
	}

	var req fcmTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	req.DeviceID = strings.TrimSpace(req.DeviceID)
	if req.Token == "" {
		WriteError(w, http.StatusBadRequest, "Token is required")
		return
	}
	if req.DeviceID == "" {
		WriteError(w, http.StatusBadRequest, "Device ID is required")
		return
	}

	token := &models.FcmToken{
		AccountID: accountID,
		Token:     req.Token,
		DeviceID:  req.DeviceID,
	}
	if err := h.fcmTokenRepo.Upsert(token); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to register FCM token")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "FCM token registered"})
}

// UnregisterToken removes the FCM token on logout.
// @Route: DELETE /api/fcm/token
func (h *FcmHandler) UnregisterToken(w http.ResponseWriter, r *http.Request) {
	accountID, ok := h.accountID(r)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Account ID not found in context")
		return
	}

	var req fcmTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	req.DeviceID = strings.TrimSpace(req.DeviceID)

	// Prefer the exact token; fall back to (account, device) when only the
	// device id is known (e.g. the FCM token could not be read on the client).
	var err error
	switch {
	case req.Token != "":
		err = h.fcmTokenRepo.DeleteByToken(req.Token)
	case req.DeviceID != "":
		err = h.fcmTokenRepo.DeleteByAccountAndDevice(accountID, req.DeviceID)
	default:
		WriteError(w, http.StatusBadRequest, "Token or device ID is required")
		return
	}

	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to unregister FCM token")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "FCM token unregistered"})
}
