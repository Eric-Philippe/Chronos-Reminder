package tests

import (
	"testing"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/services"
)

func TestRecalculateNextOccurrence(t *testing.T) {
	// Test current time: November 1, 2025 at 23:55 Europe/Paris
	parisLoc, _ := time.LoadLocation("Europe/Paris")
	LosAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
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
		{
			name:            "Workdays - future date should jump to next monday (in timezone) if triggered on weekend",
			from:            time.Date(2025, time.November, 14, 10, 0, 0, 0, parisLoc), // Nov 14, 2025 is Friday
			timezone:        "Europe/Paris",
			recurrenceState: services.RecurrenceWorkdays,
			expectedAfter:   time.Date(2025, time.November, 17, 10, 0, 0, 0, parisLoc), // Next Should be Monday Nov 17, 2025
			expectedHour:    10,
			expectedMinute:  0,
			checkTimeOfDay:  true,
		},
		{
			name:			"Workdays - from weekend should jump to next monday",
			from:			time.Date(2025, time.November, 14, 20, 51, 0, 0, LosAngelesLoc), // Nov 14, 2025 is Friday
			recurrenceState:	services.RecurrenceWorkdays,
			expectedAfter:		time.Date(2025, time.November, 17, 20, 51, 0, 0, LosAngelesLoc), // Next Should be Monday Nov 17, 2025
			expectedHour:		20,
			expectedMinute:		51,
			checkTimeOfDay:		true,
			timezone:		"America/Los_Angeles",
		},
		{
			name:			"Weekend - from weekday should jump to next saturday",
			from:			time.Date(2025, time.November, 14, 20, 51, 0, 0, parisLoc), // Nov 14, 2025 is Friday
			timezone:		"Europe/Paris",
			recurrenceState:	services.RecurrenceWeekend,
			expectedAfter:		time.Date(2025, time.November, 15, 20, 51, 0, 0, parisLoc), // Next Should be Saturday Nov 15, 2025
			expectedHour:		20,
			expectedMinute:		51,
			checkTimeOfDay:		true,
		
		},
		{
			name:			"Weekend - Saturday should jump to Sunday",
			from:			time.Date(2025, time.November, 15, 20, 51, 0, 0, parisLoc), // Nov 15, 2025 is Saturday
			timezone:		"Europe/Paris",
			recurrenceState:	services.RecurrenceWeekend,
			expectedAfter:		time.Date(2025, time.November, 16, 20, 51, 0, 0, parisLoc), // Next Should be Sunday Nov 16, 2025
			expectedHour:		20,
			expectedMinute:		51,
			checkTimeOfDay:		true,
		
		},
		{
			name:			"Weekend - Sunday should jump to next Saturday",
			from:			time.Date(2025, time.November, 16, 20, 51, 0, 0, parisLoc), // Nov 16, 2025 is Sunday
			timezone:		"Europe/Paris",
			recurrenceState:	services.RecurrenceWeekend,
			expectedAfter:		time.Date(2025, time.November, 22, 20, 51, 0, 0, parisLoc), // Next Should be Saturday Nov 22, 2025
			expectedHour:		20,
			expectedMinute:		51,
			checkTimeOfDay:		true,
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

	// Very old reminder — far enough in the past that many intervals have been missed
	oldReminder := time.Date(2020, time.January, 1, 12, 0, 0, 0, parisLoc)
	now := time.Now().In(parisLoc)

	t.Run("Old daily reminder should jump to future, not catch up", func(t *testing.T) {
		result, err := services.GetNextOccurrence(oldReminder, services.RecurrenceDaily, "Europe/Paris")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		resultLocal := result.In(parisLoc)

		// Should be in the future
		if !result.After(now) {
			t.Errorf("Result should be after now (%v), got: %v", now, result)
		}

		// Should maintain 12:00 time
		if resultLocal.Hour() != 12 || resultLocal.Minute() != 0 {
			t.Errorf("Expected time 12:00, got: %02d:%02d", resultLocal.Hour(), resultLocal.Minute())
		}

		// Should not be too far in the future — next occurrence, not a catch-up
		daysDiff := result.Sub(now).Hours() / 24
		if daysDiff > 2 {
			t.Errorf("Next occurrence too far in future: %v days", daysDiff)
		}

		t.Logf("✓ Jumped from Jan 1 2020 to %v without catch-up", resultLocal)
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

		t.Logf("✓ Jumped from Jan 1 2020 to %v without catch-up", result.In(parisLoc))
	})
}

// TestRecurrenceEdgeCases tests exact date transitions for quirky boundary conditions.
// All "from" dates are in 2027-2028 (future relative to real time.Now()) so GetNextOccurrence
// always takes the simple "advance by one interval" path — no catch-up logic involved.
func TestRecurrenceEdgeCases(t *testing.T) {
	parisLoc, _ := time.LoadLocation("Europe/Paris")

	// Jan 1, 2027 = Friday; Jan 2 = Sat, Jan 3 = Sun, Jan 4 = Mon, Jan 6 = Wed, Jan 7 = Thu, Jan 9 = Sat
	tests := []struct {
		name     string
		from     time.Time
		timezone string
		state    int
		expected time.Time
	}{
		{
			name:     "Workdays - Saturday skips to Monday",
			from:     time.Date(2027, time.January, 2, 10, 0, 0, 0, parisLoc), // Saturday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWorkdays,
			expected: time.Date(2027, time.January, 4, 10, 0, 0, 0, parisLoc), // Monday
		},
		{
			name:     "Workdays - Sunday skips to Monday",
			from:     time.Date(2027, time.January, 3, 10, 0, 0, 0, parisLoc), // Sunday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWorkdays,
			expected: time.Date(2027, time.January, 4, 10, 0, 0, 0, parisLoc), // Monday
		},
		{
			name:     "Workdays - Friday skips over weekend to Monday",
			from:     time.Date(2027, time.January, 1, 10, 0, 0, 0, parisLoc), // Friday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWorkdays,
			expected: time.Date(2027, time.January, 4, 10, 0, 0, 0, parisLoc), // Monday
		},
		{
			name:     "Workdays - Wednesday advances to Thursday",
			from:     time.Date(2027, time.January, 6, 10, 0, 0, 0, parisLoc), // Wednesday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWorkdays,
			expected: time.Date(2027, time.January, 7, 10, 0, 0, 0, parisLoc), // Thursday
		},
		{
			name:     "Weekend - Monday skips 5 days to Saturday",
			from:     time.Date(2027, time.January, 4, 15, 30, 0, 0, parisLoc), // Monday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWeekend,
			expected: time.Date(2027, time.January, 9, 15, 30, 0, 0, parisLoc), // Saturday
		},
		{
			name:     "Weekend - Thursday skips to Saturday",
			from:     time.Date(2027, time.January, 7, 15, 30, 0, 0, parisLoc), // Thursday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWeekend,
			expected: time.Date(2027, time.January, 9, 15, 30, 0, 0, parisLoc), // Saturday
		},
		{
			name:     "Weekend - Saturday advances to Sunday",
			from:     time.Date(2027, time.January, 9, 15, 30, 0, 0, parisLoc), // Saturday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWeekend,
			expected: time.Date(2027, time.January, 10, 15, 30, 0, 0, parisLoc), // Sunday
		},
		{
			name:     "Weekend - Sunday wraps to next Saturday (6 days)",
			from:     time.Date(2027, time.January, 10, 15, 30, 0, 0, parisLoc), // Sunday
			timezone: "Europe/Paris",
			state:    services.RecurrenceWeekend,
			expected: time.Date(2027, time.January, 16, 15, 30, 0, 0, parisLoc), // next Saturday
		},
		{
			name:     "Monthly - Jan 31 overflows to March 3 (Feb has 28 days in 2027)",
			from:     time.Date(2027, time.January, 31, 10, 0, 0, 0, parisLoc),
			timezone: "Europe/Paris",
			state:    services.RecurrenceMonthly,
			expected: time.Date(2027, time.March, 3, 10, 0, 0, 0, parisLoc),
		},
		{
			name:     "Yearly - Feb 29 2028 (leap) overflows to March 1 2029 (non-leap)",
			from:     time.Date(2028, time.February, 29, 10, 0, 0, 0, parisLoc),
			timezone: "Europe/Paris",
			state:    services.RecurrenceYearly,
			expected: time.Date(2029, time.March, 1, 10, 0, 0, 0, parisLoc),
		},
		{
			name:     "Weekly - Dec 26 2027 (Sunday) crosses year boundary to Jan 2 2028",
			from:     time.Date(2027, time.December, 26, 9, 0, 0, 0, parisLoc),
			timezone: "Europe/Paris",
			state:    services.RecurrenceWeekly,
			expected: time.Date(2028, time.January, 2, 9, 0, 0, 0, parisLoc),
		},
		{
			name:     "Daily - Dec 31 crosses year boundary to Jan 1",
			from:     time.Date(2027, time.December, 31, 22, 0, 0, 0, parisLoc),
			timezone: "Europe/Paris",
			state:    services.RecurrenceDaily,
			expected: time.Date(2028, time.January, 1, 22, 0, 0, 0, parisLoc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := services.GetNextOccurrence(tt.from, tt.state, tt.timezone)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			loc, _ := time.LoadLocation(tt.timezone)
			resultLocal := result.In(loc)
			expectedLocal := tt.expected.In(loc)

			if !resultLocal.Equal(expectedLocal) {
				t.Errorf("Expected %v, got %v", expectedLocal, resultLocal)
			}

			t.Logf("✓ %v → %v", tt.from.In(loc).Format("Mon Jan 2 2006 15:04"), resultLocal.Format("Mon Jan 2 2006 15:04"))
		})
	}
}

func TestPausedRecurrence(t *testing.T) {
	parisLoc, _ := time.LoadLocation("Europe/Paris")
	from := time.Date(2027, time.June, 15, 10, 0, 0, 0, parisLoc)

	recurrences := []struct {
		name  string
		state int
	}{
		{"Daily paused", services.BuildRecurrenceState(services.RecurrenceDaily, true)},
		{"Hourly paused", services.BuildRecurrenceState(services.RecurrenceHourly, true)},
		{"Weekly paused", services.BuildRecurrenceState(services.RecurrenceWeekly, true)},
		{"Monthly paused", services.BuildRecurrenceState(services.RecurrenceMonthly, true)},
		{"Workdays paused", services.BuildRecurrenceState(services.RecurrenceWorkdays, true)},
		{"Weekend paused", services.BuildRecurrenceState(services.RecurrenceWeekend, true)},
	}

	for _, r := range recurrences {
		t.Run(r.name, func(t *testing.T) {
			result, err := services.GetNextOccurrence(from, r.state, "Europe/Paris")
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if !result.Equal(from) {
				t.Errorf("Expected paused recurrence to return from=%v unchanged, got %v", from, result)
			}
			t.Logf("✓ Paused: returned same time %v", result.In(parisLoc))
		})
	}
}