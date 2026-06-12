package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DFMSendNow dispatches the account's DFM note immediately. It is wired by
// the engine at startup so bot code can trigger a send without importing the
// engine package (which would create an import cycle through the dispatchers).
var DFMSendNow func(accountID uuid.UUID) error

// ComputeDFMReminderSchedule resolves the first fire time (UTC) of a DFM note
// reminder. dateStr may be empty and defaults to today in the user's timezone.
// If the resolved time is already in the past, the next recurrence occurrence
// is used instead (one-time reminders in the past are rejected).
func ComputeDFMReminderSchedule(dateStr, timeStr string, recurrenceState int, ianaLocation string) (time.Time, error) {
	loc, err := time.LoadLocation(ianaLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone %s: %w", ianaLocation, err)
	}

	if dateStr == "" {
		dateStr = time.Now().In(loc).Format("2006-01-02")
	}

	parsed, err := ParseReminderDateTimeInTimezone(dateStr, timeStr, ianaLocation)
	if err != nil {
		return time.Time{}, err
	}

	if !parsed.After(time.Now()) {
		if GetRecurrenceType(recurrenceState) == RecurrenceOnce {
			return time.Time{}, fmt.Errorf("reminder time is in the past")
		}

		// Same convention as the reminder scheduler: the UTC value is the
		// reference passed to GetNextOccurrence
		next, err := GetNextOccurrence(parsed.UTC(), recurrenceState, ianaLocation)
		if err != nil {
			return time.Time{}, err
		}
		parsed = next
	}

	return parsed.UTC(), nil
}
