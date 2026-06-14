package services

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// FcmService wraps the Firebase Admin SDK messaging client used to deliver
// Android push notifications. When no credentials are configured the service is
// created in a disabled state so the rest of the app keeps working.
type FcmService struct {
	client  *messaging.Client
	enabled bool
}

// NewFcmService initializes the Firebase Admin SDK from the service account
// JSON at credentialsPath (standard Google ADC). An empty path, or any
// initialization failure, yields a disabled service rather than an error: push
// is an optional delivery channel and must never block startup.
func NewFcmService(credentialsPath string) *FcmService {
	if credentialsPath == "" {
		log.Println("[FCM] - ⚠️  GOOGLE_APPLICATION_CREDENTIALS not set, push notifications disabled")
		return &FcmService{enabled: false}
	}

	ctx := context.Background()
	log.Printf("[FCM] - 🕒 Initializing Firebase Cloud Messaging with credentials from: %s", option.WithCredentialsFile(credentialsPath))
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Printf("[FCM] - ⚠️  Failed to initialize Firebase app, push disabled: %v", err)
		return &FcmService{enabled: false}
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("[FCM] - ⚠️  Failed to create messaging client, push disabled: %v", err)
		return &FcmService{enabled: false}
	}

	log.Println("[FCM] - ✅ Firebase Cloud Messaging initialized")
	return &FcmService{client: client, enabled: true}
}

// IsEnabled reports whether push delivery is available.
func (s *FcmService) IsEnabled() bool {
	return s.enabled
}

// ErrTokenUnregistered indicates the device token is no longer valid and should
// be removed from storage.
var ErrTokenUnregistered = fmt.Errorf("fcm token unregistered")

// Send delivers a single notification to one device token. When the token is no
// longer registered it returns ErrTokenUnregistered so the caller can prune it.
//
// We send a data-only message (no Notification field) so that onMessageReceived
// is always the sole delivery path on the client. This prevents the double-notification
// problem that occurs when Firebase surfaces a Notification payload as a system
// notification independently of the app's own display logic.
func (s *FcmService) Send(ctx context.Context, token, title, body string, data map[string]string) error {
	if !s.enabled {
		return fmt.Errorf("fcm service is disabled")
	}

	// Merge title/body into the data map so the client can read them in onMessageReceived.
	payload := make(map[string]string, len(data)+2)
	for k, v := range data {
		payload[k] = v
	}
	payload["title"] = title
	payload["message"] = body

	message := &messaging.Message{
		Token: token,
		Data:  payload,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	_, err := s.client.Send(ctx, message)
	if err != nil {
		if messaging.IsUnregistered(err) {
			return ErrTokenUnregistered
		}
		return err
	}
	return nil
}
