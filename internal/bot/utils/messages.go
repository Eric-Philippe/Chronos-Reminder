package utils

import "github.com/bwmarrin/discordgo"

func BuildEmbed(session *discordgo.Session, title, description string, footerText *string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       ColorInfo,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ClockLogo,
		},
	}

	if footerText != nil {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: *footerText,
		}
	}

	return embed
}

func SendEmbed(session *discordgo.Session, interaction *discordgo.InteractionCreate, title, description string, footerText *string) error {
	embed := BuildEmbed(session, title, description, footerText)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func BuildInfoEmbed(session *discordgo.Session, title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: "ℹ️ - " + description,
		Color:       ColorInfo,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: session.State.User.AvatarURL(""),
			Text:    "Chronos Bot Reminder",
		},
	}
}

// SendInfo sends an info message as an interaction response
func SendInfo(session *discordgo.Session, interaction *discordgo.InteractionCreate, title, description string) error {
	embed := BuildInfoEmbed(session, title, description)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// BuildSuccessEmbed builds a success embed message
func BuildSuccessEmbed(session *discordgo.Session, title, description string, footerText *string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: "✅ - " + description,
		Color:       ColorSuccess,
		Author: &discordgo.MessageEmbedAuthor{
			IconURL: session.State.User.AvatarURL(""),
			Name:    "Chronos Bot Reminder",
		},
	}

	if footerText != nil {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text:    *footerText,
		}
	}

	return embed
}

// SendSuccess sends a success message as an interaction response
func SendSuccess(session *discordgo.Session, interaction *discordgo.InteractionCreate, title, description string, footerText *string) error {
	embed := BuildSuccessEmbed(session, title, description, footerText)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
		
// BuildErrorEmbed builds an error embed message
func BuildErrorEmbed(session *discordgo.Session, title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: "❌ - " + description,
		Color:       ColorError,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: session.State.User.AvatarURL(""),
			Text:    "Chronos Bot Reminder",
		},
	}
}

// SendError sends an error message as an interaction response
func SendError(session *discordgo.Session, interaction *discordgo.InteractionCreate, title, description string) error {
	embed := BuildErrorEmbed(session, title, description)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// BuildWarningEmbed builds a warning embed message
func BuildWarningEmbed(session *discordgo.Session, title, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: "⚠️ - " + description,
		Color:       ColorWarning,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: session.State.User.AvatarURL(""),
			Text:    "Chronos Bot Reminder",
		},
	}
}

// SendWarning sends a warning message as an interaction response
func SendWarning(session *discordgo.Session, interaction *discordgo.InteractionCreate, title, description string) error {
	embed := BuildWarningEmbed(session, title, description)

	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}