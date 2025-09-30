package logic

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// ProfileHandler handles the profile command
func ProfileHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options

	var targetUser *discordgo.User
	var targetAccount *models.Account

	// Check if a user was specified
	if len(options) > 0 && options[0].Name == "user" {
		targetUser = options[0].UserValue(session)
		// Get or create account for target user
		var err error
		targetAccount, err = services.GetAccountFromDiscordUser(targetUser)
		// If error or no account, display a message saying that user has no account
		if err != nil || targetAccount == nil {
			msg := "They can create one by setting a reminder or calling that command !"
			return utils.SendErrorDetailed(session, interaction, "No Account", fmt.Sprintf("User %s does not have an account yet !", targetUser.Username), &msg)
		}
	} else {
		// Use the command invoker's account
		if interaction.Member != nil && interaction.Member.User != nil {
			targetUser = interaction.Member.User
		} else if interaction.User != nil {
			targetUser = interaction.User
		}
		targetAccount = account
	}

	if targetUser == nil || targetAccount == nil {
		return utils.SendError(session, interaction, "Error", "Unable to determine target user.")
	}

	// Fetch user's reminder count
	reminderCount := len(targetAccount.Reminders)

	// Download user's avatar
	var avatarImage image.Image
	if targetUser.Avatar != "" {
		avatarURL := targetUser.AvatarURL("256")
		var err error
		avatarImage, err = downloadAvatar(avatarURL)
		if err != nil {
			// If avatar download fails, continue without it
			avatarImage = nil
		}
	}

	// Determine badges based on user's platform usage
	var badges []string
	// Add Discord badge (since they're using the Discord bot)
	badges = append(badges, "Discord")

	// Check if user has used the app
	if services.DiscordUserUsesApp(targetAccount) {
		badges = append(badges, "App")
	}

	// Create profile data
	profileData := services.ProfileData{
		Username:      targetUser.Username,
		Avatar:        avatarImage,
		CreatedAt:     targetAccount.CreatedAt,
		ReminderCount: reminderCount,
		Badges:        badges,
	}

	// Generate profile image
	drawService := services.NewDrawService("./assets")
	profileImage, err := drawService.GenerateProfileImage(profileData)
	if err != nil {
		log.Println("Error generating profile image:", err)
		return utils.SendError(session, interaction, "Error", "Failed to generate profile image.")
	}

	// Convert image to bytes
	var buf bytes.Buffer
	if err := png.Encode(&buf, profileImage); err != nil {
		return utils.SendError(session, interaction, "Error", "Failed to encode profile image.")
	}

	// Send response with image attachment
	return session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{
				{
					Name:   "profile.png",
					Reader: bytes.NewReader(buf.Bytes()),
				},
			},
		},
	})
}

// downloadAvatar downloads a user's avatar from Discord
func downloadAvatar(avatarURL string) (image.Image, error) {
	resp, err := http.Get(avatarURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download avatar: status %d", resp.StatusCode)
	}

	// If gif, decode as gif
	if strings.HasSuffix(avatarURL, ".gif") {
		gifImg, err := gif.Decode(resp.Body)
		if err != nil {
			return nil, err
		}
		return gifImg, nil
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
