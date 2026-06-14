package dispatchers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/google/uuid"
)

// AndroidPushDispatcher delivers reminders to a user's Android devices via FCM.
type AndroidPushDispatcher struct {
	fcmService   *services.FcmService
	fcmTokenRepo repositories.FcmTokenRepository
}

// NewAndroidPushDispatcher creates a new Android push dispatcher
func NewAndroidPushDispatcher(fcmService *services.FcmService, fcmTokenRepo repositories.FcmTokenRepository) *AndroidPushDispatcher {
	return &AndroidPushDispatcher{
		fcmService:   fcmService,
		fcmTokenRepo: fcmTokenRepo,
	}
}

// GetSupportedType returns the destination type this dispatcher supports
func (d *AndroidPushDispatcher) GetSupportedType() models.DestinationType {
	return models.DestinationAndroidPush
}

// Dispatch sends the reminder as a push notification to every registered token
// for the account in the destination metadata. Tokens that FCM reports as
// unregistered are pruned from the database.
func (d *AndroidPushDispatcher) Dispatch(reminder *models.Reminder, destination *models.ReminderDestination, account *models.Account) error {
	if destination.Type != models.DestinationAndroidPush {
		return fmt.Errorf("invalid destination type for android push dispatcher: %s", destination.Type)
	}

	if d.fcmService == nil || !d.fcmService.IsEnabled() {
		return fmt.Errorf("push notifications are not configured")
	}

	accountIDVal, exists := destination.Metadata["account_id"]
	if !exists {
		return fmt.Errorf("account_id not found in destination metadata")
	}
	accountIDStr, ok := accountIDVal.(string)
	if !ok || accountIDStr == "" {
		return fmt.Errorf("account_id in destination metadata is not a valid string")
	}
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return fmt.Errorf("account_id in destination metadata is not a valid UUID: %w", err)
	}

	tokens, err := d.fcmTokenRepo.GetByAccountID(accountID)
	if err != nil {
		return fmt.Errorf("failed to load FCM tokens: %w", err)
	}
	if len(tokens) == 0 {
		return fmt.Errorf("no FCM tokens registered for account %s", accountID)
	}

	ctx := context.Background()
	data := map[string]string{"reminder_id": reminder.ID.String()}

	var sendErrors []error
	for _, token := range tokens {
		err := d.fcmService.Send(ctx, token.Token, "Chronos", reminder.Message, data)
		if err == nil {
			continue
		}

		if errors.Is(err, services.ErrTokenUnregistered) {
			// Stale token — remove it so we stop trying to reach a dead device.
			if delErr := d.fcmTokenRepo.DeleteByToken(token.Token); delErr != nil {
				log.Printf("[ANDROID_PUSH] - Failed to delete unregistered token: %v", delErr)
			}
			continue
		}

		sendErrors = append(sendErrors, err)
	}

	if len(sendErrors) > 0 {
		return fmt.Errorf("failed to deliver push to %d/%d devices: %v", len(sendErrors), len(tokens), sendErrors[0])
	}

	return nil
}
