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