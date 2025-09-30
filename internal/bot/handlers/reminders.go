package handlers

import "github.com/ericp/chronos-bot-reminder/internal/bot/logic"

func init() {
	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID: "confirm_delete_",
		Handler:  logic.HandleConfirmDelete,
		NeedsAccount: true,
	})
	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID: "cancel_delete_",
		Handler:  logic.HandleCancelDelete,
	})
	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID: "show_reminder_",
		Handler:  logic.HandleShowReminderFromList,
		NeedsAccount: true,
	})
	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID: "back_to_list",
		Handler:  logic.HandleBackToList,
		NeedsAccount: true,
	})
}