package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseReminderTime parses various time formats and returns a time.Time
func ParseReminderTime(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	now := time.Now()

	// Try parsing as time only first (for today)
	if timeOnlyResult, err := parseTimeOnlyForToday(timeStr, now); err == nil {
		return timeOnlyResult, nil
	}

	// Try parsing as absolute datetime formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"01/02/2006 15:04",
		"01/02/2006 3:04 PM",
		"2006-01-02 3:04 PM",
		"January 2, 2006 3:04 PM",
		"Jan 2, 2006 3:04 PM",
		"2 Jan 2006 15:04",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			// If no year specified and parsed time is in the past, assume next year
			if t.Year() == 0 {
				t = t.AddDate(now.Year(), 0, 0)
			}
			if t.Before(now) && !strings.Contains(timeStr, strconv.Itoa(now.Year())) {
				t = t.AddDate(1, 0, 0)
			}
			return t, nil
		}
	}

	// Try parsing relative time (e.g., "1h 30m", "2d", "30 minutes")
	if relativeTime, err := parseRelativeTime(timeStr); err == nil {
		return now.Add(relativeTime), nil
	}

	// Try parsing simple words (tomorrow, today, etc.)
	if simpleTime, err := parseSimpleTime(timeStr); err == nil {
		return simpleTime, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse time format: %s", timeStr)
}

// ParseReminderDateTime combines separate date and time strings into a time.Time
func ParseReminderDateTime(dateStr, timeStr string) (time.Time, error) {
	now := time.Now()
	
	// Parse the date component
	parsedDate, err := parseDateOnly(dateStr, now)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %w", err)
	}
	
	// Parse the time component
	parsedTime, err := parseTimeOfDay(timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time: %w", err)
	}
	
	// Combine date and time
	result := time.Date(
		parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
		0, now.Location(),
	)
	
	return result, nil
}

// parseTimeOnlyForToday parses time-only formats (like "10:30", "10h30", "3pm") and applies them to today's date
func parseTimeOnlyForToday(timeStr string, now time.Time) (time.Time, error) {
	// Normalize the input - replace 'h' with ':'
	normalizedTime := strings.ReplaceAll(timeStr, "h", ":")
	
	// Time formats to try for time-only input
	timeOnlyFormats := []string{
		"15:04",     // 24-hour format like "14:30"
		"3:04 PM",   // 12-hour format with space like "2:30 PM"
		"3:04PM",    // 12-hour format without space like "2:30PM"
		"3PM",       // Hour only with PM like "3PM"
		"3pm",       // Hour only with pm like "3pm"
		"15",        // 24-hour format hour only like "14"
	}

	for _, format := range timeOnlyFormats {
		if t, err := time.Parse(format, normalizedTime); err == nil {
			// Apply the parsed time to today's date
			todayWithTime := time.Date(now.Year(), now.Month(), now.Day(), 
				t.Hour(), t.Minute(), t.Second(), 0, now.Location())
			
			// If the time is in the past today, schedule it for tomorrow
			if todayWithTime.Before(now) {
				todayWithTime = todayWithTime.AddDate(0, 0, 1)
			}
			
			return todayWithTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time-only format: %s", timeStr)
}

// parseRelativeTime parses relative time expressions like "1h 30m", "2 days", etc.
func parseRelativeTime(timeStr string) (time.Duration, error) {
	var totalDuration time.Duration

	// Replace common words with abbreviations
	timeStr = strings.ReplaceAll(timeStr, " minutes", "m")
	timeStr = strings.ReplaceAll(timeStr, " minute", "m")
	timeStr = strings.ReplaceAll(timeStr, " hours", "h")
	timeStr = strings.ReplaceAll(timeStr, " hour", "h")
	timeStr = strings.ReplaceAll(timeStr, " days", "d")
	timeStr = strings.ReplaceAll(timeStr, " day", "d")
	timeStr = strings.ReplaceAll(timeStr, " weeks", "w")
	timeStr = strings.ReplaceAll(timeStr, " week", "w")

	// Split by spaces and try to parse each part
	parts := strings.Fields(timeStr)
	for _, part := range parts {
		// Try parsing as standard duration first
		if d, err := time.ParseDuration(part); err == nil {
			totalDuration += d
			continue
		}

		// Try parsing custom formats like "2d", "1w"
		if len(part) >= 2 {
			numStr := part[:len(part)-1]
			unit := part[len(part)-1:]

			if num, err := strconv.Atoi(numStr); err == nil {
				switch unit {
				case "d":
					totalDuration += time.Duration(num) * 24 * time.Hour
				case "w":
					totalDuration += time.Duration(num) * 7 * 24 * time.Hour
				}
			}
		}
	}

	if totalDuration == 0 {
		return 0, fmt.Errorf("no valid duration found")
	}

	return totalDuration, nil
}

// parseSimpleTime parses simple time expressions like "tomorrow", "today", etc.
func parseSimpleTime(timeStr string) (time.Time, error) {
	now := time.Now()
	lower := strings.ToLower(strings.TrimSpace(timeStr))

	switch lower {
	case "now":
		return now, nil
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()), nil
	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 12, 0, 0, 0, now.Location()), nil
	case "next week":
		nextWeek := now.AddDate(0, 0, 7)
		return time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 12, 0, 0, 0, now.Location()), nil
	case "next month":
		nextMonth := now.AddDate(0, 1, 0)
		return time.Date(nextMonth.Year(), nextMonth.Month(), nextMonth.Day(), 12, 0, 0, 0, now.Location()), nil
	}

	return time.Time{}, fmt.Errorf("unable to parse simple time: %s", timeStr)
}

// parseTimeOfDay parses time expressions like "3pm", "15:30", "9:30am"
func parseTimeOfDay(timeStr string) (time.Time, error) {
	// Normalize the input - replace 'h' with ':'
	normalizedTime := strings.ReplaceAll(timeStr, "h", ":")
	
	timeFormats := []string{
		"15:04",     // 24-hour format like "14:30"
		"3:04 PM",   // 12-hour format with space like "2:30 PM"
		"3:04PM",    // 12-hour format without space like "2:30PM"
		"3:04 pm",   // 12-hour format with lowercase pm
		"3:04am",    // 12-hour format with lowercase am
		"3:04 am",   // 12-hour format with lowercase am and space
		"3PM",       // Hour only with PM like "3PM"
		"3pm",       // Hour only with pm like "3pm"
		"3AM",       // Hour only with AM like "3AM"
		"3am",       // Hour only with am like "3am"
		"15",        // 24-hour format hour only like "14"
	}

	for _, format := range timeFormats {
		if t, err := time.Parse(format, normalizedTime); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time of day: %s", timeStr)
}

// parseDateOnly parses date strings and returns a time.Time with date only
func parseDateOnly(dateStr string, now time.Time) (time.Time, error) {
	dateStr = strings.TrimSpace(strings.ToLower(dateStr))
	
	// Handle simple date expressions
	switch dateStr {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location()), nil
	case "next week":
		nextWeek := now.AddDate(0, 0, 7)
		return time.Date(nextWeek.Year(), nextWeek.Month(), nextWeek.Day(), 0, 0, 0, 0, now.Location()), nil
	case "next month":
		nextMonth := now.AddDate(0, 1, 0)
		return time.Date(nextMonth.Year(), nextMonth.Month(), nextMonth.Day(), 0, 0, 0, 0, now.Location()), nil
	}
	
	// Date formats to try
	dateFormats := []string{
		"2006-01-02",      // YYYY-MM-DD
		"2006/01/02",      // YYYY/MM/DD
		"01/02/2006",      // MM/DD/YYYY (US format)
		"02/01/2006",      // DD/MM/YYYY (European format)
		"01-02-2006",      // MM-DD-YYYY
		"02-01-2006",      // DD-MM-YYYY
		"01/02",           // MM/DD (current year)
		"02/01",           // DD/MM (current year) 
		"01-02",           // MM-DD (current year)
		"02-01",           // DD-MM (current year)
		"02",              // DD (current month/year)
		"January 2, 2006", // Full month name
		"Jan 2, 2006",     // Short month name
		"2 Jan 2006",      // Day month year
		"January 2",       // Month day (current year)
		"Jan 2",           // Short month day (current year)
		"2",               // Day only (current month/year)
	}
	
	for _, format := range dateFormats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			// Handle cases where year is missing
			if parsedTime.Year() == 0 || parsedTime.Year() == 1 {
				parsedTime = parsedTime.AddDate(now.Year()-1, 0, 0)
			}
			
			// Handle cases where month is missing (day only)
			if parsedTime.Month() == 1 && parsedTime.Day() != 1 && !strings.Contains(dateStr, "jan") {
				parsedTime = time.Date(now.Year(), now.Month(), parsedTime.Day(), 0, 0, 0, 0, now.Location())
			}
			
			// Special handling for ambiguous DD/MM vs MM/DD formats
			// If we have something like "25/12", it's clearly DD/MM (day > 12)
			// If we have "12/25", it's clearly MM/DD (month > 12)
			if strings.Contains(dateStr, "/") || strings.Contains(dateStr, "-") {
				parts := strings.FieldsFunc(dateStr, func(c rune) bool { return c == '/' || c == '-' })
				if len(parts) >= 2 {
					if first, err1 := strconv.Atoi(parts[0]); err1 == nil {
						if second, err2 := strconv.Atoi(parts[1]); err2 == nil {
							// If first number > 12, it must be DD/MM format
							if first > 12 && format == "01/02/2006" {
								// Try DD/MM format instead
								if ddmmTime, err := time.Parse("02/01/2006", dateStr); err == nil {
									parsedTime = ddmmTime
								}
							}
							// If second number > 12, it must be MM/DD format
							if second > 12 && format == "02/01/2006" {
								// Try MM/DD format instead
								if mmddTime, err := time.Parse("01/02/2006", dateStr); err == nil {
									parsedTime = mmddTime
								}
							}
						}
					}
				}
			}
			
			// If date is in the past and no year was specified, assume next year
			if parsedTime.Before(now) && !strings.Contains(dateStr, strconv.Itoa(now.Year())) {
				if strings.Contains(format, "2006") && !strings.Contains(dateStr, strconv.Itoa(now.Year())) {
					// Full date format but year in past, move to next year
					parsedTime = parsedTime.AddDate(1, 0, 0)
				}
			}
			
			return time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, now.Location()), nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date format: %s", dateStr)
}