package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

type MessageComponentHandler struct {
	CustomID     string
	Handler      func(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error
	NeedsAccount bool // Indicates if the dispatcher should always fill the account parameter
}

var handlers []*MessageComponentHandler

// RegisterMessageComponentHandler registers a message component handler
func RegisterMessageComponentHandler(handler *MessageComponentHandler) {
	handlers = append(handlers, handler)
}

func HandleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	customID := i.MessageComponentData().CustomID
	
	// First, check standalone message component handlers
	for _, handler := range handlers {
		// Check for exact match or prefix match (for dynamic IDs)
		if handler.CustomID == customID || strings.HasPrefix(customID, handler.CustomID) {
			var account *models.Account
			var err error
			
			if handler.NeedsAccount {
				var user *discordgo.User
				if i.Member != nil && i.Member.User != nil {
					user = i.Member.User
				} else if i.User != nil {
					user = i.User
				} else {
					return nil // No user information available
				}
				
				account, err = services.EnsureDiscordUser(user)
				if err != nil {
					return err
				}
			}
			
			return handler.Handler(s, i, account)
		}
	}
	
	return nil
}