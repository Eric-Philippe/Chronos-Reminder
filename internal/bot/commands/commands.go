package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// Description represents the description object for a command
type Description struct {
	Name             string `json:"name"`
	Emoji            string `json:"emoji"`
	CategoryName     string `json:"categoryName"`
	ShortDescription string `json:"shortDescription"`
	FullDescription  string `json:"fullDescription"`
	Usage            string `json:"usage"`
	Example          string `json:"example"`
}

// AutocompleteFunc represents the autocomplete function signature
type AutocompleteFunc func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error

// RunFunc represents the run function signature
type RunFunc func(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error

type MessageComponentHandler struct {
	CustomID string
	Handler  func(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error
}


// Command represents a bot command with its description, data, and handlers
type Command struct {
	Description  Description                           `json:"description"`
	Data         *discordgo.ApplicationCommand         `json:"data"`
	Autocomplete *AutocompleteFunc                     `json:"autocomplete,omitempty"`
	NeedsAccount bool                                  `json:"needsAccount"`
	Run          RunFunc                               `json:"run"`
	MessageComponentHandlers []MessageComponentHandler `json:"messageComponentHandlers,omitempty"`
}

// commandRegistry holds registered commands by name
var commandRegistry = make(map[string]*Command)
var commands []*Command

// RegisterCommand registers a single command
func RegisterCommand(command *Command) {
	commands = append(commands, command)
	commandRegistry[command.Data.Name] = command
}

// RegisterCommands registers all commands with Discord and builds the command registry
func RegisterCommands(session *discordgo.Session) (int, error) {
	var applicationCommands []*discordgo.ApplicationCommand
	
	for _, cmd := range commands {
		applicationCommands = append(applicationCommands, cmd.Data)
	}
	
	_, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", applicationCommands)
	if err != nil {
		return 0, err
	}

	return len(applicationCommands), nil
}

// HandleCommand routes an interaction to the appropriate command handler
func HandleCommand(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return nil
	}
	
	commandName := interaction.ApplicationCommandData().Name
	command, exists := commandRegistry[commandName]
	if !exists || command.Run == nil {
		return nil
	}

	var account *models.Account
	var err error
	
	if command.NeedsAccount {
		var user *discordgo.User
		if interaction.Member != nil && interaction.Member.User != nil {
			user = interaction.Member.User
		} else if interaction.User != nil {
			user = interaction.User
		} else {
			return nil // No user information available
		}
		
		account, err = services.EnsureDiscordUser(user)
		if err != nil {
			return err
		}
	}

	return command.Run(session, interaction, account)
}

// HandleAutocomplete routes an autocomplete interaction to the appropriate command handler
func HandleAutocomplete(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
	if interaction.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return nil
	}
	
	commandName := interaction.ApplicationCommandData().Name
	command, exists := commandRegistry[commandName]
	if !exists || command.Autocomplete == nil {
		return nil
	}

	return (*command.Autocomplete)(session, interaction)
}

func HandleMessageComponent(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	customID := i.MessageComponentData().CustomID
	
	for _, command := range commandRegistry {
		for _, handler := range command.MessageComponentHandlers {
			// Check for exact match or prefix match (for dynamic IDs)
			if handler.CustomID == customID || strings.HasPrefix(customID, handler.CustomID) {
				var account *models.Account
				var err error
				
				if command.NeedsAccount {
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
	}
	
	return nil
}
