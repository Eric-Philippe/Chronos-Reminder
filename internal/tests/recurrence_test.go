package tests

import (
	"testing"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/services"
)

func TestRecalculateNextOccurrence(t *testing.T) {
	// Test current time: November 1, 2025 at 23:55 Europe/Paris
	parisLoc, _ := time.LoadLocation("Europe/Paris")
	now := time.Date(2025, time.November, 1, 23, 55, 0, 0, parisLoc)

	tests := []struct {
		name              string
		from              time.Time
		timezone          string
		recurrenceState   int
		expectedAfter     time.Time
		expectedHour      int // Expected hour of day in local timezone
		expectedMinute    int // Expected minute in local timezone
		checkTimeOfDay    bool
	}{
		{
			name:            "Hourly - from October 1st 23:00, should be next hour after now",
			from:            time.Date(2025, time.October, 1, 23, 0, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceHourly,
			expectedAfter:   now,
			expectedHour:    0, // Should be November 2nd 00:xx
			checkTimeOfDay:  false,
		},
		{
			name:            "Hourly - from October 1st 23:55, should be next hour after now",
			from:            time.Date(2025, time.October, 1, 23, 55, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceHourly,
			expectedAfter:   now,
			checkTimeOfDay:  false,
		},
		{
			name:            "Daily - from October 1st 23:59, should preserve time across DST",
			from:            time.Date(2025, time.October, 1, 23, 59, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceDaily,
			expectedAfter:   now,
			expectedHour:    23,
			expectedMinute:  59,
			checkTimeOfDay:  true,
		},
		{
			name:            "Daily - from October 1st 14:00, should preserve time across DST",
			from:            time.Date(2025, time.October, 1, 14, 0, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceDaily,
			expectedAfter:   now,
			expectedHour:    14,
			expectedMinute:  0,
			checkTimeOfDay:  true,
		},
		{
			name:            "Weekly - from October 1st 10:30, should preserve time across DST",
			from:            time.Date(2025, time.October, 1, 10, 30, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceWeekly,
			expectedAfter:   now,
			expectedHour:    10,
			expectedMinute:  30,
			checkTimeOfDay:  true,
		},
		{
			name:            "Monthly - from October 1st 09:15, should preserve time across DST",
			from:            time.Date(2025, time.October, 1, 9, 15, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceMonthly,
			expectedAfter:   now,
			expectedHour:    9,
			expectedMinute:  15,
			checkTimeOfDay:  true,
		},
		{
			name:            "Yearly - from October 1st 2024 18:45, should preserve time",
			from:            time.Date(2024, time.October, 1, 18, 45, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceYearly,
			expectedAfter:   now,
			expectedHour:    18,
			expectedMinute:  45,
			checkTimeOfDay:  true,
		},
		{
			name:            "Daily - future date should just add one day",
			from:            time.Date(2025, time.November, 3, 10, 0, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceDaily,
			expectedAfter:   time.Date(2025, time.November, 3, 10, 0, 0, 0, parisLoc),
			expectedHour:    10,
			expectedMinute:  0,
			checkTimeOfDay:  true,
		},
		{
			name:            "Hourly - future date should just add one hour",
			from:            time.Date(2025, time.November, 3, 10, 0, 0, 0, parisLoc),
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceHourly,
			expectedAfter:   time.Date(2025, time.November, 3, 10, 0, 0, 0, parisLoc),
			checkTimeOfDay:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := services.GetNextOccurrence(test.from, test.recurrenceState, test.timezone)
			if err != nil {
				t.Errorf("Did not expect error, got: %v", err)
				return
			}

			// Load timezone for comparison
			loc, _ := time.LoadLocation(test.timezone)
			resultLocal := result.In(loc)

			// Check that result is after expected time
			if !result.After(test.expectedAfter) {
				t.Errorf("Expected result after %v, got: %v", test.expectedAfter, result)
			}

			// Check time of day preservation (for daily+ recurrences)
			if test.checkTimeOfDay {
				if resultLocal.Hour() != test.expectedHour {
					t.Errorf("Expected hour %d, got: %d (full time: %v)", test.expectedHour, resultLocal.Hour(), resultLocal)
				}
				if resultLocal.Minute() != test.expectedMinute {
					t.Errorf("Expected minute %d, got: %d (full time: %v)", test.expectedMinute, resultLocal.Minute(), resultLocal)
				}
			}

			t.Logf("✓ Result: %v (local: %v)", result, resultLocal)
		})
	}
}

func TestDSTTransition(t *testing.T) {
	// Specific test for DST transition on October 27, 2025 (clocks go back 1 hour in Europe/Paris)
	parisLoc, _ := time.LoadLocation("Europe/Paris")
	
	// Reminder set before DST transition
	beforeDST := time.Date(2025, time.October, 1, 23, 59, 0, 0, parisLoc)
	
	// Current time after DST transition
	afterDST := time.Date(2025, time.November, 1, 23, 55, 0, 0, parisLoc)
	
	t.Run("Daily recurrence should maintain local time across DST", func(t *testing.T) {
		result, err := services.GetNextOccurrence(beforeDST, services.RecurrenceDaily, "Europe/Paris")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		resultLocal := result.In(parisLoc)
		
		// Should be after current time
		if !result.After(afterDST) {
			t.Errorf("Result should be after %v, got: %v", afterDST, result)
		}
		
		// Should maintain 23:59 local time despite DST change
		if resultLocal.Hour() != 23 || resultLocal.Minute() != 59 {
			t.Errorf("Expected time 23:59, got: %02d:%02d (full: %v)", 
				resultLocal.Hour(), resultLocal.Minute(), resultLocal)
		}
		
		t.Logf("✓ Maintained local time 23:59 across DST transition. Result: %v", resultLocal)
	})
	
	t.Run("Hourly recurrence should not skip hours", func(t *testing.T) {
		hourlyFrom := time.Date(2025, time.October, 1, 23, 0, 0, 0, parisLoc)
		
		result, err := services.GetNextOccurrence(hourlyFrom, services.RecurrenceHourly, "Europe/Paris")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Should be after current time
		if !result.After(afterDST) {
			t.Errorf("Result should be after %v, got: %v", afterDST, result)
		}
		
		t.Logf("✓ Next hourly occurrence: %v", result.In(parisLoc))
	})
}

func TestNoCatchUpLoop(t *testing.T) {
	// Test that old reminders don't trigger repeatedly to catch up
	parisLoc, _ := time.LoadLocation("Europe/Paris")
	
	// Very old reminder
	oldReminder := time.Date(2025, time.January, 1, 12, 0, 0, 0, parisLoc)
	now := time.Date(2025, time.November, 1, 23, 55, 0, 0, parisLoc)
	
	t.Run("Old daily reminder should jump to future, not catch up", func(t *testing.T) {
		result, err := services.GetNextOccurrence(oldReminder, services.RecurrenceDaily, "Europe/Paris")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		resultLocal := result.In(parisLoc)
		
		// Should be in the future, not triggering for every missed day
		if !result.After(now) {
			t.Errorf("Result should be after now (%v), got: %v", now, result)
		}
		
		// Should maintain 12:00 time
		if resultLocal.Hour() != 12 || resultLocal.Minute() != 0 {
			t.Errorf("Expected time 12:00, got: %02d:%02d", resultLocal.Hour(), resultLocal.Minute())
		}
		
		// Should not be too far in the future (should be next occurrence)
		daysDiff := result.Sub(now).Hours() / 24
		if daysDiff > 2 {
			t.Errorf("Next occurrence too far in future: %v days", daysDiff)
		}
		
		t.Logf("✓ Jumped from Jan 1 to %v without catch-up", resultLocal)
	})
	
	t.Run("Old hourly reminder should jump to future, not catch up", func(t *testing.T) {
		result, err := services.GetNextOccurrence(oldReminder, services.RecurrenceHourly, "Europe/Paris")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Should be in the future
		if !result.After(now) {
			t.Errorf("Result should be after now (%v), got: %v", now, result)
		}
		
		// Should not be too far in the future (within next few hours)
		hoursDiff := result.Sub(now).Hours()
		if hoursDiff > 2 {
			t.Errorf("Next occurrence too far in future: %v hours", hoursDiff)
		}
		
		t.Logf("✓ Jumped from Jan 1 to %v without catch-up", result.In(parisLoc))
	})
}