package api

import (
	"encoding/json"
	"net/http"

	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// APIKeyHandler handles API key operations
type APIKeyHandler struct {
	apiKeyService *services.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(apiKeyService *services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// CreateAPIKeyRequest represents the request to create a new API key
type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

// ListAPIKeysResponse represents the response with all API keys
type ListAPIKeysResponse struct {
	Keys []APIKeyResponse `json:"keys"`
}

// APIKeyResponse represents an API key in responses
type APIKeyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Scopes    string `json:"scopes"`
	CreatedAt string `json:"created_at"`
	Key       string `json:"key,omitempty"` // Only populated on creation
}

// CreateAPIKey creates a new API key
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from context
	accountIDVal := r.Context().Value(AccountIDKey)
	if accountIDVal == nil {
		WriteError(w, http.StatusUnauthorized, "Account ID not found in context")
		return
	}

	accountID, ok := accountIDVal.(uuid.UUID)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Invalid account ID in context")
		return
	}

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		WriteError(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Create the API key
	metadata, err := h.apiKeyService.CreateAPIKey(accountID, req.Name)
	if err != nil {
		if err.Error() == "maximum of 5 API keys per account" {
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to create API key")
		return
	}

	resp := APIKeyResponse{
		ID:        metadata.ID,
		Name:      metadata.Name,
		Scopes:    metadata.Scopes,
		CreatedAt: metadata.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Key:       metadata.Key,
	}

	WriteJSON(w, http.StatusCreated, resp)
}

// GetAPIKeys retrieves all API keys for the account
func (h *APIKeyHandler) GetAPIKeys(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from context
	accountIDVal := r.Context().Value(AccountIDKey)
	if accountIDVal == nil {
		WriteError(w, http.StatusUnauthorized, "Account ID not found in context")
		return
	}

	accountID, ok := accountIDVal.(uuid.UUID)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Invalid account ID in context")
		return
	}

	keys, err := h.apiKeyService.GetAPIKeys(accountID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve API keys")
		return
	}

	responses := make([]APIKeyResponse, len(keys))
	for i, key := range keys {
		responses[i] = APIKeyResponse{
			ID:        key.ID,
			Name:      key.Name,
			Scopes:    key.Scopes,
			CreatedAt: key.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	WriteJSON(w, http.StatusOK, ListAPIKeysResponse{Keys: responses})
}

// RevokeAPIKey revokes an API key
func (h *APIKeyHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	// Extract account ID from context
	accountIDVal := r.Context().Value(AccountIDKey)
	if accountIDVal == nil {
		WriteError(w, http.StatusUnauthorized, "Account ID not found in context")
		return
	}

	accountID, ok := accountIDVal.(uuid.UUID)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "Invalid account ID in context")
		return
	}

	keyID := r.PathValue("id")

	if err := h.apiKeyService.RevokeAPIKey(accountID, keyID); err != nil {
		if err.Error() == "unauthorized" || err.Error() == "API key not found" {
			WriteError(w, http.StatusNotFound, "API key not found")
			return
		}
		if err.Error() == "not an API key" {
			WriteError(w, http.StatusBadRequest, "Not an API key")
			return
		}
		WriteError(w, http.StatusInternalServerError, "Failed to revoke API key")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "API key revoked successfully",
	})
}
