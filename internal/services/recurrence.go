package services

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

// GetRecurrenceTypeName returns the string name for a recurrence type
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