package services

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

// DrawService handles image generation with text overlays
type DrawService struct {
	assetsPath string
}

// NewDrawService creates a new instance of DrawService
func NewDrawService(assetsPath string) *DrawService {
	return &DrawService{
		assetsPath: assetsPath,
	}
}

// TextOverlay represents text to be overlaid on an image
type TextOverlay struct {
	Label string
	Date  time.Time
}

// ProfileData represents data for profile image generation
type ProfileData struct {
	Username     string
	Avatar       image.Image
	CreatedAt    time.Time
	ReminderCount int
	Badges       []string // "Discord", "App"
}

// GenerateReminderImage creates a reminder image with text overlay
func (ds *DrawService) GenerateReminderImage(overlay TextOverlay) (image.Image, error) {
	// Load the template image
	templateImage, err := gg.LoadImage(filepath.Join(ds.assetsPath, "/templates/NewReminder.png"))
	if err != nil {
		return nil, fmt.Errorf("failed to load template image: %w", err)
	}

	// Create a new context with the template image
	dc := gg.NewContextForImage(templateImage)

	// Set text properties
	labelFontSize := 87.0
	dateFontSize := 50.0
	
	// Load font for label
	if err := ds.loadFont(dc, labelFontSize); err != nil {
		return nil, err
	}

	// Set text color to white
	dc.SetRGB(1, 1, 1)

	// Get image dimensions
	width := float64(templateImage.Bounds().Dx())
	height := float64(templateImage.Bounds().Dy())

	// Draw the label text at the center
	if overlay.Label != "" {
		labelX := width / 2
		labelY := height / 2
		maxWidth := width * 0.8 // 80% of image width for text
		
		// Adjust font size based on text length and available width
		adjustedFontSize := ds.calculateOptimalFontSize(dc, overlay.Label, maxWidth, labelFontSize)
		if err := ds.loadFont(dc, adjustedFontSize); err != nil {
			return nil, err
		}
		
		ds.DrawMultilineText(dc, overlay.Label, labelX, labelY, maxWidth, adjustedFontSize*1.2)
	}

	// Draw the date at bottom left at coordinates (76, 270)
	if overlay.Date != (time.Time{}) {
		// Load smaller font for date
		if err := ds.loadFont(dc, dateFontSize); err != nil {
			return nil, err
		}
		
		// Format date as DD/MM/YYYY · HH:MM
		formattedDate := ds.formatDateString(overlay.Date)
		ds.drawLeftAlignedText(dc, formattedDate, 20, 784)
	}

	return dc.Image(), nil
}

// GenerateProfileImage creates a profile image with user data
func (ds *DrawService) GenerateProfileImage(data ProfileData) (image.Image, error) {
	// Load the template image
	templateImage, err := gg.LoadImage(filepath.Join(ds.assetsPath, "/templates/ProfileTemplate.png"))
	if err != nil {
		return nil, fmt.Errorf("failed to load template image: %w", err)
	}

	// Create a new context with the template image
	dc := gg.NewContextForImage(templateImage)


		// Set text color to white for all text
	dc.SetRGB(1, 1, 1)

	// Draw the username under the avatar
	if data.Username != "" {
		usernameFontSize := 40.0
		if err := ds.loadFont(dc, usernameFontSize); err != nil {
			return nil, err
		}
		
		usernameX := 525.0
		usernameY := 300.0 // Below the avatar
		ds.drawCenteredText(dc, data.Username, usernameX, usernameY)
	}

	// Draw the creation date
	if !data.CreatedAt.IsZero() {
		dateFontSize := 23.0
		if err := ds.loadFont(dc, dateFontSize); err != nil {
			return nil, err
		}
		
		dateText := fmt.Sprintf("Created on %s", ds.formatProfileDate(data.CreatedAt))
		dateX := 125.0 
		dateY := 255.0
		ds.drawLeftAlignedText(dc, dateText, dateX, dateY)
	}

	// Draw the reminder count
	reminderFontSize := 27.0
	if err := ds.loadFont(dc, reminderFontSize); err != nil {
		return nil, err
	}
	
	reminderText := "No Reminders yet"

	if data.ReminderCount == 1 {
		reminderText = "1 Reminder"
	} else if data.ReminderCount > 1 {
		reminderText = fmt.Sprintf("%d Reminders", data.ReminderCount)
	}
	reminderX := 210.0 // Same X as date
	reminderY := 178.0 // Below the date
	ds.drawLeftAlignedText(dc, reminderText, reminderX, reminderY)

	// Draw the badges at the bottom
	if len(data.Badges) > 0 {
		ds.drawBadges(dc, data.Badges)
	}

	// Draw the avatar (circled) in the middle right
	if data.Avatar != nil {
		avatarX := 525.0
		avatarY := 160.0
		avatarSize := 120.0
		ds.drawCircularAvatar(dc, data.Avatar, avatarX, avatarY, avatarSize)
	}

	return dc.Image(), nil
}

// WrapText wraps text to fit within a specified width
func (ds *DrawService) WrapText(text string, maxWidth float64, dc *gg.Context) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word

		width, _ := dc.MeasureString(testLine)
		if width <= maxWidth || currentLine == "" {
			currentLine = testLine
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// DrawMultilineText draws text that can span multiple lines
func (ds *DrawService) DrawMultilineText(dc *gg.Context, text string, x, y, maxWidth, lineHeight float64) {
	lines := ds.WrapText(text, maxWidth, dc)

	startY := y - (float64(len(lines)-1) * lineHeight / 2)

	for i, line := range lines {
		lineY := startY + (float64(i) * lineHeight)
		ds.drawCenteredText(dc, line, x, lineY)
	}
}

// loadFont loads a font at the specified size
func (ds *DrawService) loadFont(dc *gg.Context, fontSize float64) error {
	// Try to load the target font
	if err := dc.LoadFontFace("./assets/fonts/nasa21.ttf", fontSize); err != nil {
		return nil
	}
	return nil
}

// calculateOptimalFontSize calculates the best font size for text to fit within maxWidth
func (ds *DrawService) calculateOptimalFontSize(dc *gg.Context, text string, maxWidth, startingSize float64) float64 {
	fontSize := startingSize
	
	for fontSize > 22.4 {
		if err := ds.loadFont(dc, fontSize); err != nil {
			return fontSize
		}
		
		width, _ := dc.MeasureString(text)
		if width <= maxWidth {
			return fontSize
		}
		
		fontSize -= 2.8
	}
	
	return 22.4
}

// formatDateString formats a date string to DD/MM/YYYY . HH:MM format
func (ds *DrawService) formatDateString(date time.Time) string {
	// Replace space with · if it exists
	parts := strings.SplitN(date.String(), " ", 2)
	if len(parts) == 2 {
		// Remove anything from parts[1] after the minutes
		timeParts := strings.SplitN(parts[1], ":", 3)
		return strings.ReplaceAll(parts[0], "-", "/") + " at " + timeParts[0] + ":" + timeParts[1]
	}

	return date.String()
}

// drawCenteredText draws text centered at the given coordinates
func (ds *DrawService) drawCenteredText(dc *gg.Context, text string, x, y float64) {
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)
}

// drawLeftAlignedText draws text left-aligned at the given coordinates
func (ds *DrawService) drawLeftAlignedText(dc *gg.Context, text string, x, y float64) {
	dc.DrawStringAnchored(text, x, y, 0, 0.5)
}

// drawCircularAvatar draws a circular avatar image at the specified position
func (ds *DrawService) drawCircularAvatar(dc *gg.Context, avatar image.Image, x, y, size float64) {
	// Save the current state
	dc.Push()
	
	// Create a circular clipping path
	dc.DrawCircle(x, y, size/2)
	dc.Clip()
	
	// Calculate position to center the avatar in the circle
	bounds := avatar.Bounds()
	avatarWidth := float64(bounds.Dx())
	avatarHeight := float64(bounds.Dy())
	
	// Scale to fit the circle
	scale := size / avatarWidth
	if avatarHeight > avatarWidth {
		scale = size / avatarHeight
	}
	
	// Draw the avatar centered in the circle
	avatarX := x - (avatarWidth*scale)/2
	avatarY := y - (avatarHeight*scale)/2
	
	dc.Scale(scale, scale)
	dc.DrawImage(avatar, int(avatarX/scale), int(avatarY/scale))
	
	// Restore the state
	dc.Pop()
}

// formatProfileDate formats a date for profile display (DD/MM/YYYY)
func (ds *DrawService) formatProfileDate(date time.Time) string {
	return date.Format("02/01/2006")
}

// drawBadges draws badges at the bottom of the profile image
func (ds *DrawService) drawBadges(dc *gg.Context, badges []string) {
	badgeY := 350.0 // Y position for badges based on template
	badgeStartX := 89.0 // Starting X position for badges
	badgeSpacing := 60.0 // Spacing between badges
	badgeSize := 40.0 // Size of badge images
	
	for i, badge := range badges {
		var badgeImagePath string
		switch badge {
		case "Discord":
			badgeImagePath = filepath.Join(ds.assetsPath, "/badges/BadgeDiscord.png")
		case "App":
			badgeImagePath = filepath.Join(ds.assetsPath, "/badges/BadgeApp.png")
		default:
			continue // Skip unknown badges
		}
		
		// Load badge image
		badgeImage, err := gg.LoadImage(badgeImagePath)
		if err != nil {
			continue // Skip if badge image can't be loaded
		}
		
		// Calculate position for this badge
		badgeX := badgeStartX + float64(i)*badgeSpacing
		
		// Draw the badge image
		ds.drawBadgeImage(dc, badgeImage, badgeX, badgeY, badgeSize)
	}
}

// drawBadgeImage draws a badge image at the specified position and size
func (ds *DrawService) drawBadgeImage(dc *gg.Context, badgeImage image.Image, x, y, size float64) {
	bounds := badgeImage.Bounds()
	badgeWidth := float64(bounds.Dx())
	badgeHeight := float64(bounds.Dy())
	
	// Calculate scale to fit the desired size
	scale := size / badgeWidth
	if badgeHeight > badgeWidth {
		scale = size / badgeHeight
	}
	
	// Save current state
	dc.Push()
	
	// Position and scale for the badge
	dc.Translate(x, y)
	dc.Scale(scale, scale)
	
	// Draw the badge centered
	dc.DrawImage(badgeImage, int(-badgeWidth/2), int(-badgeHeight/2))
	
	// Restore state
	dc.Pop()
}