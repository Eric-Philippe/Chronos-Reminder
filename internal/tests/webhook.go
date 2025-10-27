package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	webhookURL := "https://webhook.site/3a50d443-7530-4828-96df-c0660c8d50b6"

	payload := map[string]interface{}{
		"content": "Hello from Go webhook test!", // optional plain text
		"embeds": []map[string]interface{}{
			{
				"title":       "User Update",
				"description": "A user has just updated their profile.",
				"color":       0x00FF00, // green
				"fields": []map[string]string{
					{"name": "User", "value": "Eric", "inline": "true"},
					{"name": "Status", "value": "Active", "inline": "true"},
				},
				"footer": map[string]string{
					"text": "Webhook Test • 2025-10-27",
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status 200 OK")
}