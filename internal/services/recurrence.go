package services

import (
	"fmt"
	"log"
	"time"
)

// Recurrence type constants
const (
	RecurrenceOnce     = 0
	RecurrenceYearly   = 1
	RecurrenceMonthly  = 2
	RecurrenceWeekly   = 3
	RecurrenceDaily    = 4
	RecurrenceHourly   = 5
	RecurrenceWorkdays = 6
	RecurrenceWeekend  = 7
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
	// 1 day = 86400 seconds, 5 workdays = 432000 seconds
	daysToAdd := interval
	weeksToAdd := daysToAdd / 5
	extraDays := daysToAdd % 5

	// Calculate total seconds to add
	totalSeconds := int64(weeksToAdd*7*86400 + extraDays*86400)

	// Adjust for weekends
	dayOfWeek := (from / 86400) % 7 // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	if dayOfWeek+int64(extraDays) >= 6 { // If it goes into the weekend
		totalSeconds += 2 * 86400 // Skip Saturday and Sunday
	}

	return from + totalSeconds
}

// WeekendRecurrence struct
type WeekendRecurrence struct{}

// NextOccurrence returns the next occurrence timestamp for weekend recurrence
func (r WeekendRecurrence) NextOccurrence(from int64, interval int) int64 {
	// 1 day = 86400 seconds, 2 weekend days = 172800 seconds
	daysToAdd := interval
	weeksToAdd := daysToAdd / 2
	extraDays := daysToAdd % 2

	// Calculate total seconds to add
	totalSeconds := int64(weeksToAdd*7*86400 + extraDays*86400)

	// Adjust for weekdays
	dayOfWeek := (from / 86400) % 7 // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	if dayOfWeek+int64(extraDays) < 6 { // If it goes into the weekday
		totalSeconds += (6 - dayOfWeek) * 86400 // Skip to Saturday
	}

	return from + totalSeconds
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

	nextTimestamp := recurrence.NextOccurrence(fromLocal.Unix(), 1)
	nextTimeLocal := time.Unix(nextTimestamp, 0).In(loc)
	return nextTimeLocal, nil
}

// RecalculateNextOccurrence is a helper to recalculate the next occurrence for a reminder that has been paused and then unpaused
// ianaLocation is the IANA timezone identifier for the user (e.g., "Europe/Paris")
func RecalculateNextOccurrence(from time.Time, recurrenceState int, ianaLocation string) (time.Time, error) {
	// Extract the actual recurrence type from the bit-encoded state
	recurrenceType := GetRecurrenceType(recurrenceState)
	isPaused := IsPaused(recurrenceState)

	// If the recurrence is still paused, return the same time
	if isPaused {
		return from, nil
	}

	// Load the user's timezone to properly handle DST transitions
	loc, err := time.LoadLocation(ianaLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load timezone %s: %w", ianaLocation, err)
	}

	recurrence := Recurrences[GetRecurrenceTypeName(recurrenceType)]
	if recurrence == nil {
		return time.Time{}, fmt.Errorf("invalid recurrence type: %d (extracted from state: %d)", recurrenceType, recurrenceState)
	}

	// Get current time in the user's timezone
	now := time.Now().In(loc)
	
	// For special recurrence types (workdays, weekend), we need to ensure we land on the correct day type
	switch recurrenceType {
	case RecurrenceWorkdays:
		// Find the next workday (Monday-Friday)
		for {
			weekday := now.Weekday()
			if weekday >= time.Monday && weekday <= time.Friday {
				break
			}
			now = now.AddDate(0, 0, 1)
		}
	case RecurrenceWeekend:
		// Find the next weekend day (Saturday-Sunday)
		for {
			weekday := now.Weekday()
			if weekday == time.Saturday || weekday == time.Sunday {
				break
			}
			now = now.AddDate(0, 0, 1)
		}
	}
	
	// Use the original reminder time for hour/minute/second precision
	nextTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		from.Hour(), from.Minute(), from.Second(), from.Nanosecond(),
		loc,
	)
	
	// If the calculated time is in the past (same day but earlier time), move to next occurrence
	if nextTime.Before(now) {
		// Create a local time version of 'from' for the calculation
		fromLocal := time.Date(
			from.Year(), from.Month(), from.Day(),
			from.Hour(), from.Minute(), from.Second(), from.Nanosecond(),
			loc,
		)
		nextTimestamp := recurrence.NextOccurrence(fromLocal.Unix(), 1)
		nextTime = time.Unix(nextTimestamp, 0).In(loc)
	}
	
	log.Printf("[RECURRENCE] - Recalculated next occurrence from %v to %v for recurrence type %s", 
		from, nextTime, GetRecurrenceTypeName(recurrenceType))
	
	return nextTime, nil
}