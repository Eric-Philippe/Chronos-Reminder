package tests

import (
	"testing"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/services"
	"github.com/stretchr/testify/assert"
)


func TestParseDateTime(t *testing.T) {
	// Load timezones for proper offset handling
	nyLoc, _ := time.LoadLocation("America/New_York")
	parisLoc, _ := time.LoadLocation("Europe/Paris")
	
	tests := []struct {
		input_date  string
		input_time  string
		timezone    string
		expected    time.Time
		shouldError bool
	}{
		{
			input_date: "11/10/2025",
			input_time: "15:00",
			timezone: "America/New_York",
			expected: time.Date(time.Now().Year(), 11, 10, 15, 0, 0, 0, nyLoc),
		},
		{
			input_date: "03/11",
			input_time: "15:00",
			timezone: "America/New_York",
			expected: time.Date(time.Now().Year(), 3, 11, 15, 0, 0, 0, nyLoc),
		},
		{
			// 3 November current year in Paris timezone
			input_date: "03/11",
			input_time: "15:00",
			timezone: "Europe/Paris",
			expected: time.Date(time.Now().Year(), 11, 3, 15, 0, 0, 0, parisLoc),
		},
	}

	for _, test := range tests {
		result, err := services.ParseReminderDateTimeInTimezone(test.input_date, test.input_time, test.timezone)
		if test.shouldError {
			assert.Error(t, err, "Expected error for input: %s", test.input_date+" "+test.input_time)
		} else {
			assert.NoError(t, err, "Did not expect error for input: %s", test.input_date+" "+test.input_time)
			assert.Equal(t, test.expected, result, "Unexpected result for input: %s", test.input_date+" "+test.input_time)
		}
	}
}

