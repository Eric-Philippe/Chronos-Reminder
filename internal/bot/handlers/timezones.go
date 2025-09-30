package handlers

import "github.com/ericp/chronos-bot-reminder/internal/bot/logic"

func init() {
	RegisterMessageComponentHandler(&MessageComponentHandler{
		CustomID: "timezone_change_select",
		Handler:  logic.HandleTimezoneSelectMenu,
		NeedsAccount: true,
	})
}

