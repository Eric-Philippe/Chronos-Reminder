package logic

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ericp/chronos-bot-reminder/internal/bot/utils"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// calculTimeHandler handles the time calculation command
func CalculTimeHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate, account *models.Account) error {
	options := interaction.ApplicationCommandData().Options

	var timeOne string
	var operation string
	var timeTwo string

	// Parse command options
	for _, option := range options {
		switch option.Name {
		case "time1":
			timeOne = option.StringValue()
		case "operation":
			operation = option.StringValue()
		case "time2":
			timeTwo = option.StringValue()
		}
	}

	// Validate inputs
	if timeOne == "" || operation == "" || timeTwo == "" {
		return utils.SendError(session, interaction, "Missing Parameters", 
			"Please provide all required parameters: time1, operation, and time2/factor.")
	}

	// Parse the first time
	parsedTime1, err := parseTimeInput(timeOne)
	if err != nil {
		return utils.SendError(session, interaction, "Invalid Time Format", 
			fmt.Sprintf("Could not parse time1 '%s'. Use formats like '2h 30m', '14:30', or '2.5h'.", timeOne))
	}

	var result time.Duration
	var resultStr string

	switch strings.ToUpper(operation) {
	case "ADD", "+":
		parsedTime2, err := parseTimeInput(timeTwo)
		if err != nil {
			return utils.SendError(session, interaction, "Invalid Time Format", 
				fmt.Sprintf("Could not parse time2 '%s'. Use formats like '2h 30m', '14:30', or '2.5h'.", timeTwo))
		}
		result = parsedTime1 + parsedTime2
		resultStr = fmt.Sprintf("%s + %s = %s", 
			formatDuration(parsedTime1), formatDuration(parsedTime2), formatDuration(result))

	case "SUBTRACT", "SUB", "-":
		parsedTime2, err := parseTimeInput(timeTwo)
		if err != nil {
			return utils.SendError(session, interaction, "Invalid Time Format", 
				fmt.Sprintf("Could not parse time2 '%s'. Use formats like '2h 30m', '14:30', or '2.5h'.", timeTwo))
		}
		result = parsedTime1 - parsedTime2
		if result < 0 {
			result = -result
			resultStr = fmt.Sprintf("%s - %s = -%s", 
				formatDuration(parsedTime1), formatDuration(parsedTime2), formatDuration(result))
		} else {
			resultStr = fmt.Sprintf("%s - %s = %s", 
				formatDuration(parsedTime1), formatDuration(parsedTime2), formatDuration(result))
		}

	case "MULTIPLY", "MUL", "*", "×":
		factor, err := parseFactorInput(timeTwo)
		if err != nil {
			return utils.SendError(session, interaction, "Invalid Factor", 
				fmt.Sprintf("Could not parse factor '%s'. Use a number like '2', '1.5', or '0.5'.", timeTwo))
		}
		result = time.Duration(float64(parsedTime1) * factor)
		resultStr = fmt.Sprintf("%s × %.2f = %s", 
			formatDuration(parsedTime1), factor, formatDuration(result))

	case "DIVIDE", "DIV", "/", "÷":
		factor, err := parseFactorInput(timeTwo)
		if err != nil {
			return utils.SendError(session, interaction, "Invalid Factor", 
				fmt.Sprintf("Could not parse factor '%s'. Use a number like '2', '1.5', or '0.5'.", timeTwo))
		}
		if factor == 0 {
			return utils.SendError(session, interaction, "Division by Zero", 
				"Cannot divide by zero.")
		}
		result = time.Duration(float64(parsedTime1) / factor)
		resultStr = fmt.Sprintf("%s ÷ %.2f = %s", 
			formatDuration(parsedTime1), factor, formatDuration(result))

	default:
		return utils.SendError(session, interaction, "Invalid Operation", 
			fmt.Sprintf("Unknown operation '%s'. Use: add (+), subtract (-), multiply (×), or divide (÷).", operation))
	}

	// Format additional information
	additionalInfo := fmt.Sprintf("**Result in different formats:**\n• Hours: %.2f h\n• Minutes: %.0f min\n• Seconds: %.0f sec", 
		result.Hours(), result.Minutes(), result.Seconds())

	description := fmt.Sprintf("**Calculation:** %s\n\n%s", resultStr, additionalInfo)

	return utils.SendEmbed(session, interaction, "Time Calculation Result 🧮", description, nil)
}

// parseTimeInput parses various time input formats and returns a duration
func parseTimeInput(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	
	// Try parsing as duration first (e.g., "2h30m", "1.5h", "90m")
	if duration, err := parseDurationInput(input); err == nil {
		return duration, nil
	}

	// Try parsing as time-of-day (e.g., "14:30", "2:30 PM") and convert to duration since midnight
	if timeOfDay, err := services.ParseReminderTime(input); err == nil {
		// Convert time-of-day to duration since midnight
		midnight := time.Date(timeOfDay.Year(), timeOfDay.Month(), timeOfDay.Day(), 0, 0, 0, 0, timeOfDay.Location())
		return timeOfDay.Sub(midnight), nil
	}

	return 0, fmt.Errorf("unable to parse time input: %s", input)
}

// parseDurationInput parses duration strings with various formats
func parseDurationInput(input string) (time.Duration, error) {
	var totalDuration time.Duration

	// Replace common words and formats
	input = strings.ReplaceAll(input, " hours", "h")
	input = strings.ReplaceAll(input, " hour", "h")
	input = strings.ReplaceAll(input, " minutes", "m")
	input = strings.ReplaceAll(input, " minute", "m")
	input = strings.ReplaceAll(input, " seconds", "s")
	input = strings.ReplaceAll(input, " second", "s")
	input = strings.ReplaceAll(input, " hrs", "h")
	input = strings.ReplaceAll(input, " hr", "h")
	input = strings.ReplaceAll(input, " mins", "m")
	input = strings.ReplaceAll(input, " min", "m")
	input = strings.ReplaceAll(input, " secs", "s")
	input = strings.ReplaceAll(input, " sec", "s")

	// Try parsing as standard Go duration
	if d, err := time.ParseDuration(input); err == nil {
		return d, nil
	}

	// Try parsing decimal hours (e.g., "2.5h", "1.75")
	if strings.HasSuffix(input, "h") {
		if hours, err := strconv.ParseFloat(strings.TrimSuffix(input, "h"), 64); err == nil {
			return time.Duration(hours * float64(time.Hour)), nil
		}
	}

	// Try parsing as plain number (assume hours)
	if hours, err := strconv.ParseFloat(input, 64); err == nil {
		return time.Duration(hours * float64(time.Hour)), nil
	}

	// Try parsing space-separated components (e.g., "2 hours 30 minutes")
	parts := strings.Fields(input)
	for i := 0; i < len(parts); i++ {
		if i+1 < len(parts) {
			// Check if next part is a unit
			numStr := parts[i]
			unitStr := parts[i+1]

			if num, err := strconv.ParseFloat(numStr, 64); err == nil {
				switch strings.ToLower(unitStr) {
				case "h", "hour", "hours", "hr", "hrs":
					totalDuration += time.Duration(num * float64(time.Hour))
					i++ // Skip the unit part
				case "m", "min", "mins", "minute", "minutes":
					totalDuration += time.Duration(num * float64(time.Minute))
					i++ // Skip the unit part
				case "s", "sec", "secs", "second", "seconds":
					totalDuration += time.Duration(num * float64(time.Second))
					i++ // Skip the unit part
				}
			}
		}
	}

	if totalDuration > 0 {
		return totalDuration, nil
	}

	return 0, fmt.Errorf("unable to parse duration: %s", input)
}

// parseFactorInput parses factor input for multiplication and division
func parseFactorInput(input string) (float64, error) {
	input = strings.TrimSpace(input)
	return strconv.ParseFloat(input, 64)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "-" + formatDuration(-d)
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}