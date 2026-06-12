package dispatchers

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// DFMDispatcher sends "Don't Forget Me" notes to their owner. The note is
// always private: it goes to the user's Discord DM when a Discord identity
// exists, otherwise it falls back to the app identity email.
type DFMDispatcher struct {
	session   *discordgo.Session
	mailer    *services.MailerService
	webAppURL string
}

// NewDFMDispatcher creates a new DFM dispatcher
func NewDFMDispatcher(mailer *services.MailerService, webAppURL string) *DFMDispatcher {
	return &DFMDispatcher{
		session:   bot.GetDiscordSession(),
		mailer:    mailer,
		webAppURL: webAppURL,
	}
}

// NoteWebURL returns the web application URL of the DFM page
func (d *DFMDispatcher) NoteWebURL() string {
	return strings.TrimSuffix(d.webAppURL, "/") + "/dont-forget-me"
}

// Dispatch sends the note to its owner on every enabled private channel
// (Discord DM and/or email). A failure on one channel does not prevent the
// other from being attempted.
func (d *DFMDispatcher) Dispatch(note *models.DFMNote, identities []models.Identity) error {
	var discordIdentity, appIdentity *models.Identity
	for i := range identities {
		switch identities[i].Provider {
		case models.ProviderDiscord:
			discordIdentity = &identities[i]
		case models.ProviderApp:
			appIdentity = &identities[i]
		}
	}

	if !note.SendDiscordDM && !note.SendEmail {
		return fmt.Errorf("no destination enabled for DFM note %s", note.ID)
	}

	var errs []error

	if note.SendDiscordDM {
		if discordIdentity == nil || d.session == nil {
			errs = append(errs, fmt.Errorf("no Discord identity linked for DFM note %s", note.ID))
		} else if err := d.dispatchDiscordDM(note, discordIdentity.ExternalID); err != nil {
			errs = append(errs, err)
		}
	}

	if note.SendEmail {
		if appIdentity == nil {
			errs = append(errs, fmt.Errorf("no email identity linked for DFM note %s", note.ID))
		} else if err := d.dispatchEmail(note, appIdentity.ExternalID); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// RenderDFMNoteText renders the note items as a plain text checklist
func RenderDFMNoteText(note *models.DFMNote) string {
	if len(note.Items) == 0 {
		return "Your note is empty."
	}

	var builder strings.Builder
	for _, item := range note.Items {
		if item.Checked {
			builder.WriteString(fmt.Sprintf("[x] %s\n", item.Content))
		} else {
			builder.WriteString(fmt.Sprintf("[ ] %s\n", item.Content))
		}
	}
	return strings.TrimRight(builder.String(), "\n")
}

// dispatchDiscordDM sends the note content as a private Discord message
func (d *DFMDispatcher) dispatchDiscordDM(note *models.DFMNote, discordUserID string) error {
	dmChannel, err := d.session.UserChannelCreate(discordUserID)
	if err != nil {
		return fmt.Errorf("failed to create DM channel with user %s: %w", discordUserID, err)
	}

	var description strings.Builder
	if len(note.Items) == 0 {
		description.WriteString("Your note is empty.")
	} else {
		for _, item := range note.Items {
			if item.Checked {
				description.WriteString(fmt.Sprintf("✅ ~~%s~~\n", item.Content))
			} else {
				description.WriteString(fmt.Sprintf("⬜ %s\n", item.Content))
			}
		}
	}
	description.WriteString("\nYou can edit your note, check items and manage the reminder from the web application.")

	embed := &discordgo.MessageEmbed{
		Title:       "💭 Don't Forget Me - Your note",
		Description: description.String(),
		Color:       0xCEA04D,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Chronos Bot Reminder",
		},
	}

	msg := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Open in the web app",
						Style: discordgo.LinkButton,
						URL:   d.NoteWebURL(),
					},
				},
			},
		},
	}

	if _, err := d.session.ChannelMessageSendComplex(dmChannel.ID, msg); err != nil {
		return fmt.Errorf("failed to send DFM note DM: %w", err)
	}
	return nil
}

// dispatchEmail sends the note content to the user's email address
func (d *DFMDispatcher) dispatchEmail(note *models.DFMNote, email string) error {
	var itemsHTML strings.Builder
	if len(note.Items) == 0 {
		itemsHTML.WriteString("<p>Your note is empty.</p>")
	} else {
		itemsHTML.WriteString("<ul style=\"list-style: none; padding-left: 0;\">")
		for _, item := range note.Items {
			content := html.EscapeString(item.Content)
			if item.Checked {
				itemsHTML.WriteString(fmt.Sprintf("<li style=\"margin: 6px 0;\">[x] <span style=\"text-decoration: line-through; color: #999;\">%s</span></li>", content))
			} else {
				itemsHTML.WriteString(fmt.Sprintf("<li style=\"margin: 6px 0;\">[ ] %s</li>", content))
			}
		}
		itemsHTML.WriteString("</ul>")
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Don't Forget Me</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #CEA04D;">Don't Forget Me - Your note</h2>
		%s
		<p style="margin: 30px 0;">
			<a href="%s" style="background-color: #CEA04D; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">
				Open in the web app
			</a>
		</p>
		<p style="margin-top: 30px; color: #999; font-size: 12px;">This is an automated reminder from Chronos Reminder</p>
	</div>
</body>
</html>
	`, itemsHTML.String(), d.NoteWebURL())

	textBody := fmt.Sprintf("Don't Forget Me - Your note\n\n%s\n\nEdit your note: %s", RenderDFMNoteText(note), d.NoteWebURL())

	_, err := d.mailer.SendEmail(&services.EmailRequest{
		To:       email,
		Subject:  "Don't Forget Me: your note reminder",
		HtmlBody: htmlBody,
		TextBody: textBody,
	})
	return err
}
