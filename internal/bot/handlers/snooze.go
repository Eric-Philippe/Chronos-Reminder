package handlers

import "github.com/ericp/chronos-bot-reminder/internal/bot/logic"

func init() {
	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID:     "reminder_request_snooze_",
		Handler:      logic.HandleSnooze,
		NeedsAccount: true,
	})

	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID:     "reminder_snooze_duration_",
		Handler:      logic.HandleSnoozeDuration,
		NeedsAccount: true,
	})

	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID:     "reminder_snooze_cancel",
		Handler:      logic.HandleSnoozeCancel,
		NeedsAccount: true,
	})
}