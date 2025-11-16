package services

import (
	"fmt"
	"time"
)

// Recurrence type constants
// Int representation:
const (
	RecurrenceOnce     = 0 // 0/128
	RecurrenceYearly   = 1 // 1/129
	RecurrenceMonthly  = 2 // 2/130
	RecurrenceWeekly   = 3 // 3/131
	RecurrenceDaily    = 4 // 4/132
	RecurrenceHourly   = 5 // 5/133
	RecurrenceWorkdays = 6 // 6/134
	RecurrenceWeekend  = 7 // 7/135
)

// Pause bit flag (8th bit)
const PauseBit = 1 << 7 // 128

// Map of all the recurrences and their implementation (WORKDAYS, WEEKEND, etc)
var Recurrences = map[string]Recurrence{
	"YEARLY":   YearlyRecurrence{},
	"MONTHLY":  MonthlyRecurrence{},
	"WEEKLY":   WeeklyRecurrence{},
	"DAILY":    DailyRecurrence{},
	"HOURLY":   HourlyRecurrence{},
	"WORKDAYS": WorkdaysRecurrence{},
	"WEEKEND":  WeekendRecurrence{},
}

// RecurrenceTypeMap maps string names to type constants
var RecurrenceTypeMap = map[string]int{
	"ONCE":     RecurrenceOnce,
	"YEARLY":   RecurrenceYearly,
	"MONTHLY":  RecurrenceMonthly,
	"WEEKLY":   RecurrenceWeekly,
	"DAILY":    RecurrenceDaily,
	"HOURLY":   RecurrenceHourly,
	"WORKDAYS": RecurrenceWorkdays,
	"WEEKEND":  RecurrenceWeekend,
}

// GetRecurrenceTypeName
func GetRecurrenceTypeName(recurrenceType int) string {
	typeNames := map[int]string{
		RecurrenceOnce:     "ONCE",
		RecurrenceYearly:   "YEARLY",
		RecurrenceMonthly:  "MONTHLY",
		RecurrenceWeekly:   "WEEKLY",
		RecurrenceDaily:    "DAILY",
		RecurrenceHourly:   "HOURLY",
		RecurrenceWorkdays: "WORKDAYS",
		RecurrenceWeekend:  "WEEKEND",
	}
	if name, exists := typeNames[recurrenceType]; exists {
		return name
	}
	return "UNKNOWN"
}

// BuildRecurrenceState builds a state value from recurrence type and pause status
func BuildRecurrenceState(recurrenceType int, isPaused bool) int {
	state := recurrenceType
	if isPaused {
		state |= PauseBit
	}
	return state
}

// GetRecurrenceType extracts the recurrence type from state value
func GetRecurrenceType(state int) int {
	return state & 0x7F // Mask out the pause bit
}

// IsPaused checks if the recurrence is paused from state value
func IsPaused(state int) bool {
	return (state & PauseBit) != 0
}

// SetPauseState updates the pause status in state value
func SetPauseState(state int, isPaused bool) int {
	if isPaused {
		return state | PauseBit
	}
	return state &^ PauseBit
}

// GetRecurrenceTypeLabel returns the string name for a recurrence type
func GetRecurrenceTypeLabel(recurrenceType int) string {
	typeNames := map[int]string{
		RecurrenceOnce:     "Once",
		RecurrenceYearly:   "Yearly",
		RecurrenceMonthly:  "Monthly",
		RecurrenceWeekly:   "Weekly",
		RecurrenceDaily:    "Daily",
		RecurrenceHourly:   "Hourly",
		RecurrenceWorkdays: "Workdays",
		RecurrenceWeekend:  "Weekend",
	}
	if name, exists := typeNames[recurrenceType]; exists {
		return name
	}
	return "UNKNOWN"
}

// Recurrence interface for different recurrence types
type Recurrence interface {
	NextOccurrence(from int64, interval int) int64
}

// YearlyRecurrence struct
type YearlyRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for yearly recurrence
func (r YearlyRecurrence) NextOccurrence(from int64, interval int) int64 {
	return from + int64(interval*31536000) // 365 days
}

// MonthlyRecurrence struct
type MonthlyRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for monthly recurrence
func (r MonthlyRecurrence) NextOccurrence(from int64, interval int) int64 {
	return from + int64(interval*2592000) // 30 days
}

// WeeklyRecurrence struct
type WeeklyRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for weekly recurrence
func (r WeeklyRecurrence) NextOccurrence(from int64, interval int) int64 {
	return from + int64(interval*604800) // 7 days
}

// DailyRecurrence struct
type DailyRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for daily recurrence
func (r DailyRecurrence) NextOccurrence(from int64, interval int) int64 {
	return from + int64(interval*86400) // 1 day
}

// HourlyRecurrence struct
type HourlyRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for hourly recurrence
func (r HourlyRecurrence) NextOccurrence(from int64, interval int) int64 {
	return from + int64(interval*3600) // 1 hour
}

// WorkdaysRecurrence struct
type WorkdaysRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for workdays recurrence
func (r WorkdaysRecurrence) NextOccurrence(from int64, interval int) int64 {
	// This is a placeholder - workdays/weekend need timezone context
	// The actual logic is handled in addInterval with timezone awareness
	return from + int64(interval*86400)
}

// WeekendRecurrence struct
type WeekendRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for weekend recurrence
func (r WeekendRecurrence) NextOccurrence(from int64, interval int) int64 {
	// This is a placeholder - workdays/weekend need timezone context
	// The actual logic is handled in addInterval with timezone awareness
	return from + int64(interval*86400)
}

// findNextFutureOccurrence calculates the next occurrence that is in the future
func findNextFutureOccurrence(from time.Time, recurrence Recurrence, loc *time.Location) time.Time {
	now := time.Now().In(loc)

	// Determine if we should preserve time-of-day (not for hourly recurrence)
	preserveTimeOfDay := true
	if _, isHourly := recurrence.(HourlyRecurrence); isHourly {
		preserveTimeOfDay = false
	}

	// If 'from' is already in the future, calculate next occurrence from it
	if from.After(now) {
		var nextTime time.Time
		if preserveTimeOfDay {
			// For daily+ recurrences, add days while preserving time-of-day
			nextTime = addInterval(from, recurrence, 1, loc)
		} else {
			// For hourly, just add the fixed seconds
			nextTimestamp := recurrence.NextOccurrence(from.Unix(), 1)
			nextTime = time.Unix(nextTimestamp, 0).In(loc)
		}
		return nextTime
	}

	// If 'from' is in the past or equal to now, calculate how many intervals to skip
	var nextTime time.Time

	if preserveTimeOfDay {
		// For daily+ recurrences, calculate approximate intervals to skip
		var intervalsToSkip int
		switch recurrence.(type) {
		case DailyRecurrence:
			days := int(now.Sub(from).Hours() / 24)
			intervalsToSkip = days + 1
		case WeeklyRecurrence:
			days := int(now.Sub(from).Hours() / 24)
			intervalsToSkip = (days / 7) + 1
		case MonthlyRecurrence:
			months := (now.Year()-from.Year())*12 + int(now.Month()-from.Month())
			intervalsToSkip = months + 1
		case YearlyRecurrence:
			years := now.Year() - from.Year()
			intervalsToSkip = years + 1
		default:
			// For workdays/weekend, approximate with days
			days := int(now.Sub(from).Hours() / 24)
			intervalsToSkip = days + 1
		}
		
		// Jump ahead by the calculated intervals
		nextTime = addInterval(from, recurrence, intervalsToSkip, loc)
		
		// Fine-tune: if we're still in the past, keep adding intervals
		maxAttempts := 10
		for attempt := 0; attempt < maxAttempts && !nextTime.After(now); attempt++ {
			nextTime = addInterval(nextTime, recurrence, 1, loc)
		}
	} else {
		// For hourly recurrence, use fixed seconds calculation
		secondsDiff := now.Unix() - from.Unix()
		intervalSeconds := int64(3600)
		intervalsToSkip := int(secondsDiff/intervalSeconds) + 2

		nextTimestamp := recurrence.NextOccurrence(from.Unix(), intervalsToSkip)
		nextTime = time.Unix(nextTimestamp, 0).In(loc)

		// Safety check
		maxAttempts := 10
		for attempt := 0; !nextTime.After(now) && attempt < maxAttempts; attempt++ {
			nextTimestamp = recurrence.NextOccurrence(nextTime.Unix(), 1)
			nextTime = time.Unix(nextTimestamp, 0).In(loc)
		}
	}

	return nextTime
}

// addInterval adds one interval to the given time while preserving time-of-day and handling DST
func addInterval(t time.Time, recurrence Recurrence, intervals int, loc *time.Location) time.Time {
	// Preserve the time-of-day components
	hour, minute, second := t.Hour(), t.Minute(), t.Second()
	nanosecond := t.Nanosecond()

	var nextTime time.Time
	switch recurrence.(type) {
	case DailyRecurrence:
		nextTime = t.AddDate(0, 0, intervals)
	case WeeklyRecurrence:
		nextTime = t.AddDate(0, 0, 7*intervals)
	case MonthlyRecurrence:
		nextTime = t.AddDate(0, intervals, 0)
	case YearlyRecurrence:
		nextTime = t.AddDate(intervals, 0, 0)
	case WorkdaysRecurrence:
		// Add days while skipping weekends
		daysAdded := 0
		currentTime := t
		for daysAdded < intervals {
			currentTime = currentTime.AddDate(0, 0, 1)
			// Check if the day is a weekday (Monday-Friday)
			// time.Weekday: Sunday=0, Monday=1, ..., Friday=5, Saturday=6
			if currentTime.Weekday() > 0 && currentTime.Weekday() < 6 {
				daysAdded++
			}
		}
		nextTime = currentTime
	case WeekendRecurrence:
		// Add days while skipping weekdays
		daysAdded := 0
		currentTime := t
		for daysAdded < intervals {
			currentTime = currentTime.AddDate(0, 0, 1)
			// Check if the day is a weekend day (Saturday or Sunday)
			// time.Weekday: Sunday=0, Saturday=6
			if currentTime.Weekday() == 0 || currentTime.Weekday() == 6 {
				daysAdded++
			}
		}
		nextTime = currentTime
	default:
		// Fallback
		nextTime = t.AddDate(0, 0, intervals)
	}

	// Restore the original time-of-day to handle DST properly
	nextTime = time.Date(
		nextTime.Year(), nextTime.Month(), nextTime.Day(),
		hour, minute, second, nanosecond,
		loc,
	)

	return nextTime
}

// GetNextOccurrence calculates the next occurrence timestamp based on recurrence state (with bits) and interval
// ianaLocation is the IANA timezone identifier for the user (e.g., "Europe/Paris")
func GetNextOccurrence(from time.Time, recurrenceState int, ianaLocation string) (time.Time, error) {
	// Extract the actual recurrence type from the bit-encoded state
	recurrenceType := GetRecurrenceType(recurrenceState)
	isPaused := IsPaused(recurrenceState)

	// If the recurrence is paused, return the same time
	if isPaused {
		return from, nil
	}

	// Load the user's timezone to properly handle DST transitions
	loc, err := time.LoadLocation(ianaLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone %s: %w", ianaLocation, err)
	}

	// Interpret 'from' as a local time in the user's timezone
	// from is stored with UTC location, but represents local time
	fromLocal := time.Date(
		from.Year(), from.Month(), from.Day(),
		from.Hour(), from.Minute(), from.Second(), from.Nanosecond(),
		loc,
	)

	recurrence := Recurrences[GetRecurrenceTypeName(recurrenceType)]
	if recurrence == nil {
		return time.Time{}, fmt.Errorf("invalid recurrence type: %d (extracted from state: %d)", recurrenceType, recurrenceState)
	}

	// Find the next future occurrence by iterating through past ones if needed
	// maxIterations prevents infinite loops for edge cases (set to 1000 as safety limit)
	nextTimeLocal := findNextFutureOccurrence(fromLocal, recurrence, loc)
	return nextTimeLocal, nil
}