package tests

import (
	"testing"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/services"
)

func TestRecalculateNextOccurrence(t *testing.T) {
	// Use dynamic dates relative to current time
	now := time.Date(2025, time.May, 13, 15, 0, 0, 0, time.UTC)

	tests := []struct {
		from          time.Time
		timezone	  string
		recurrenceState int
		expectedAfter time.Time // Expected result should be after this time
	}{
		{
			from: time.Date(2025, time.May, 13, 14, 0, 0, 0, time.UTC),
			timezone: "Europe/Paris",
			recurrenceState: services.RecurrenceDaily,
			expectedAfter: now,
		},
		{
			from: time.Date(2025, time.May, 12, 16, 0, 0, 0, time.UTC),
			timezone: "Europe/Paris",
			recurrenceState: services.RecurrenceWeekly,
			expectedAfter: now,
		},
		{
			from: time.Date(2025, time.May, 1, 10, 0, 0, 0, time.UTC),
			timezone: "Europe/Paris",
			recurrenceState: services.RecurrenceMonthly,
			expectedAfter: now,
		},
	}

	for _, test := range tests {
		result, err := services.GetNextOccurrence(test.from, test.recurrenceState, test.timezone)
		if err != nil {
			t.Errorf("Did not expect error for from: %v, timezone: %s, recurrenceState: %d, got: %v", test.from, test.timezone, test.recurrenceState, err)
		} else if !result.After(test.expectedAfter) {
			t.Errorf("Unexpected result for from: %v, timezone: %s, recurrenceState: %d, expected after: %v, got: %v", test.from, test.timezone, test.recurrenceState, test.expectedAfter, result)
		}
	}
}