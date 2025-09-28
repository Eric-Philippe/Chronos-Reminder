package engine

import (
	"context"
	"log"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/ericp/chronos-bot-reminder/internal/services"
)

// SchedulerEvent represents different types of events that can trigger a reschedule
type SchedulerEvent struct {
	Type string // "created", "updated", "deleted"
}

// Scheduler manages the timing and dispatching of reminders
type Scheduler struct {
	reminderRepo       repositories.ReminderRepository
	dispatcherRegistry *DispatcherRegistry
	stopChan           chan struct{}
	updateChan         chan SchedulerEvent
	running            bool
	currentTimer       *time.Timer
}

// NewScheduler creates a new scheduler instance
func NewScheduler(reminderRepo repositories.ReminderRepository, dispatcherRegistry *DispatcherRegistry) *Scheduler {
	return &Scheduler{
		reminderRepo:       reminderRepo,
		dispatcherRegistry: dispatcherRegistry,
		stopChan:           make(chan struct{}),
		updateChan:         make(chan SchedulerEvent, 100), // Buffered channel for updates
		running:            false,
		currentTimer:       nil,
	}
}

// Start begins the scheduler's main loop
func (s *Scheduler) Start(ctx context.Context) {
	if s.running {
		log.Println("[ENGINE] - Scheduler already running")
		return
	}

	s.running = true

	// Start the main scheduling loop
	go s.scheduleLoop(ctx)
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	if !s.running {
		log.Println("[ENGINE] - Scheduler already stopped")
		return
	}
	
	// Stop the current timer if it exists
	if s.currentTimer != nil {
		s.currentTimer.Stop()
		s.currentTimer = nil
	}
	
	close(s.stopChan)
	s.running = false
}

// IsRunning returns whether the scheduler is currently running
func (s *Scheduler) IsRunning() bool {
	return s.running
}

// NotifyReminderCreated notifies the scheduler that a new reminder was created
func (s *Scheduler) NotifyReminderCreated() {
	if !s.running {
		return
	}
	
	select {
	case s.updateChan <- SchedulerEvent{Type: "created"}:
		if config.IsDebugMode() {
			log.Println("[ENGINE] - Notified of reminder creation")
		}
	default:
		log.Println("[ENGINE] - Update channel full, skipping creation notification")
	}
}

// NotifyReminderUpdated notifies the scheduler that a reminder was updated
func (s *Scheduler) NotifyReminderUpdated() {
	if !s.running {
		return
	}
	
	select {
	case s.updateChan <- SchedulerEvent{Type: "updated"}:
		if config.IsDebugMode() {
			log.Println("[ENGINE] - Notified of reminder update")
		}
	default:
		log.Println("[ENGINE] - Update channel full, skipping update notification")
	}
}

// NotifyReminderDeleted notifies the scheduler that a reminder was deleted
func (s *Scheduler) NotifyReminderDeleted() {
	if !s.running {
		return
	}
	
	select {
	case s.updateChan <- SchedulerEvent{Type: "deleted"}:
		if config.IsDebugMode() {
			log.Println("[ENGINE] - Notified of reminder deletion")
		}
	default:
		log.Println("[ENGINE] - Update channel full, skipping deletion notification")
	}
}

// scheduleLoop is the main loop that waits for the next reminder or updates
func (s *Scheduler) scheduleLoop(ctx context.Context) {
	// Initial schedule setup
	s.scheduleNext()

	for {
		select {
		case <-ctx.Done():
			log.Println("[ENGINE] - Context cancelled, stopping scheduler")
			s.running = false
			return
		case <-s.stopChan:
			s.running = false
			return
		case <-s.updateChan:
			if config.IsDebugMode() {
				log.Println("[ENGINE] - Received update event, rescheduling...")
			}
			s.scheduleNext()
		case <-s.getTimerChan():
			// Timer fired, process due reminders
			s.checkAndProcessReminders()
			// Schedule the next batch
			s.scheduleNext()
		}
	}
}

// getTimerChan returns the timer channel or a nil channel if no timer is set
func (s *Scheduler) getTimerChan() <-chan time.Time {
	if s.currentTimer != nil {
		return s.currentTimer.C
	}
	// Return a channel that will never fire
	return make(<-chan time.Time)
}

// scheduleNext sets up the timer for the next reminder(s)
func (s *Scheduler) scheduleNext() {
	// Stop existing timer
	if s.currentTimer != nil {
		s.currentTimer.Stop()
		s.currentTimer = nil
	}

	// Get the next reminders
	nextReminders, err := s.reminderRepo.GetNextReminders()
	if err != nil {
		log.Printf("[ENGINE] - Error fetching next reminders: %v", err)
		return
	}

	if len(nextReminders) == 0 {
		log.Println("[ENGINE] - No upcoming reminders, waiting for updates...")
		return
	}

	// Calculate time until next reminder
	nextTime := nextReminders[0].RemindAtUTC
	now := time.Now().UTC()
	duration := nextTime.Sub(now)

	if duration <= 0 {
		// Reminder is already due, process immediately
		log.Printf("[ENGINE] - Reminder is already due, processing immediately")
		go func() {
			time.Sleep(10 * time.Millisecond) // Small delay to avoid tight loop
			select {
			case <-s.stopChan:
				return
			default:
				s.checkAndProcessReminders()
				s.scheduleNext()
			}
		}()
		return
	}

	if config.IsDebugMode() {
		log.Printf("[ENGINE] - Next reminder at %v (in %v)", nextTime, duration)
	}
	
	s.currentTimer = time.NewTimer(duration)
}

// checkAndProcessReminders fetches and processes all due reminders
func (s *Scheduler) checkAndProcessReminders() {
	// Get the next reminders (ones with the closest time)
	nextReminders, err := s.reminderRepo.GetNextReminders()
	if err != nil {
		log.Printf("[ENGINE] - Error fetching next reminders: %v", err)
		return
	}

	if len(nextReminders) == 0 {
		log.Println("[ENGINE] - No upcoming reminders found")
		return
	}

	// Check if any of the next reminders are due (they should be, since timer fired)
	now := time.Now().UTC()
	tolerance := time.Minute // Allow 1 minute tolerance

	var dueReminders []models.Reminder
	for _, reminder := range nextReminders {
		if reminder.RemindAtUTC.Before(now.Add(tolerance)) {
			dueReminders = append(dueReminders, reminder)
		}
	}

	if len(dueReminders) == 0 {
		log.Println("[ENGINE] - Timer fired but no reminders are actually due")
		return
	}

	if (config.IsDebugMode()) {
		log.Printf("[ENGINE] - Found %d due reminders", len(dueReminders))
	}

	// Process each due reminder
	for _, reminder := range dueReminders {
		s.processReminder(&reminder)
	}
}

// processReminder handles the dispatching of a single reminder
func (s *Scheduler) processReminder(reminder *models.Reminder) {
	// Dispatch the reminder to all its destinations
	err := s.dispatcherRegistry.DispatchReminder(reminder)
	if err != nil {
		log.Printf("[ENGINE] - Error dispatching reminder %s: %v", reminder.ID, err)
		return
	}

	// If it's a one-time reminder, delete it
	if reminder.Recurrence == 0 {
		err := s.reminderRepo.Delete(reminder.ID, false)
		if err != nil {
			log.Printf("[ENGINE] - Error deleting completed reminder %s: %v", reminder.ID, err)
		}
	} else {
		s.handleRecurrence(reminder)
	}
}

// handleRecurrence manages recurring reminders
func (s *Scheduler) handleRecurrence(reminder *models.Reminder) {
	if reminder.Recurrence == 0 {
		return // No recurrence
	}

	newTime, err := services.GetNextOccurrence(reminder.RemindAtUTC, int(reminder.Recurrence))
	if err != nil {
		log.Printf("[ENGINE] - Error getting next occurrence for reminder %s: %v", reminder.ID, err)
		return
	}

	err = s.reminderRepo.Reschedule(reminder.ID, newTime, false)
	if err != nil {
		log.Printf("[ENGINE] - Error rescheduling recurring reminder %s: %v", reminder.ID, err)
	}
}

