package engine

import (
	"context"
	"log"
	"time"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
	"github.com/google/uuid"
)

const (
	GarbageCollectionDelay = 30 * time.Minute
)

// GarbageCollector manages the delayed deletion of dispatched one-time reminders
type GarbageCollector struct {
	reminderRepo repositories.ReminderRepository
	stopChan     chan struct{}
	updateChan   chan uuid.UUID // Reminder ID that was updated/snoozed
	addChan      chan uuid.UUID // Reminder ID to add to deletion queue
	running      bool
	currentTimer *time.Timer
	ctx          context.Context
}

// NewGarbageCollector creates a new garbage collector instance
func NewGarbageCollector(reminderRepo repositories.ReminderRepository) *GarbageCollector {
	return &GarbageCollector{
		reminderRepo: reminderRepo,
		stopChan:     make(chan struct{}),
		updateChan:   make(chan uuid.UUID, 100),
		addChan:      make(chan uuid.UUID, 100),
		running:      false,
		currentTimer: nil,
	}
}

// Start begins the garbage collector's main loop
func (gc *GarbageCollector) Start(ctx context.Context) {
	if gc.running {
		log.Println("[ENGINE] - Garbage collector already running")
		return
	}

	gc.running = true
	gc.ctx = ctx

	// Start the main collection loop
	go gc.collectionLoop(ctx)

	log.Println("[ENGINE] - ✅ Garbage collector started")
}

// Stop gracefully stops the garbage collector
func (gc *GarbageCollector) Stop() {
	if !gc.running {
		return
	}

	if gc.currentTimer != nil {
		gc.currentTimer.Stop()
		gc.currentTimer = nil
	}

	close(gc.stopChan)
	gc.running = false
	log.Println("[ENGINE] - Garbage collector stopped")
}

// IsRunning returns whether the garbage collector is currently running
func (gc *GarbageCollector) IsRunning() bool {
	return gc.running
}

// NotifyReminderDispatched adds a reminder to the deletion queue
func (gc *GarbageCollector) NotifyReminderDispatched(reminderID uuid.UUID) {
	if !gc.running {
		return
	}

	select {
	case gc.addChan <- reminderID:
		if config.IsDebugMode() {
			log.Printf("[ENGINE] - Reminder %s added to deletion queue", reminderID)
		}
	default:
		log.Printf("[ENGINE] - Add channel full, skipping reminder %s", reminderID)
	}
}

// NotifyReminderUpdated notifies that a reminder was updated (e.g., snoozed)
// This cancels any pending deletion for this reminder
func (gc *GarbageCollector) NotifyReminderUpdated(reminderID uuid.UUID) {
	if !gc.running {
		return
	}

	select {
	case gc.updateChan <- reminderID:
		if config.IsDebugMode() {
			log.Printf("[ENGINE] - Reminder %s updated, checking deletion queue", reminderID)
		}
	default:
		log.Printf("[ENGINE] - Update channel full, skipping update for reminder %s", reminderID)
	}
}

// collectionLoop is the main loop that manages the deletion queue
func (gc *GarbageCollector) collectionLoop(ctx context.Context) {
	// Initial schedule setup
	gc.scheduleNext()

	for {
		select {
		case <-ctx.Done():
			log.Println("[ENGINE] - Context cancelled, stopping garbage collector")
			gc.running = false
			return
		case <-gc.stopChan:
			gc.running = false
			return
		case <-gc.addChan:
			// New reminder added to queue, reschedule
			if config.IsDebugMode() {
				log.Println("[ENGINE] - Reminder added, rescheduling...")
			}
			gc.scheduleNext()
		case <-gc.updateChan:
			// Reminder was updated, reschedule to check if it should still be deleted
			if config.IsDebugMode() {
				log.Println("[ENGINE] - Reminder updated, rescheduling...")
			}
			gc.scheduleNext()
		case <-gc.getTimerChan():
			// Timer fired, process deletions
			gc.checkAndDeleteReminders()
			// Schedule the next batch
			gc.scheduleNext()
		}
	}
}

// getTimerChan returns the timer channel or a nil channel if no timer is set
func (gc *GarbageCollector) getTimerChan() <-chan time.Time {
	if gc.currentTimer != nil {
		return gc.currentTimer.C
	}
	return make(<-chan time.Time)
}

// scheduleNext sets up the timer for the next reminder deletion
func (gc *GarbageCollector) scheduleNext() {
	// Stop existing timer
	if gc.currentTimer != nil {
		gc.currentTimer.Stop()
		gc.currentTimer = nil
	}

	// Get reminders ready to be deleted
	remindersToDelete, err := gc.reminderRepo.GetNextsRemindersToDelete()
	if err != nil {
		log.Printf("[ENGINE] - Error fetching reminders to delete: %v", err)
		return
	}

	if len(remindersToDelete) == 0 {
		if config.IsDebugMode() {
			log.Println("[ENGINE] - No reminders pending deletion")
		}
		return
	}

	// Find the next reminder eligible for deletion
	now := time.Now().UTC()
	var nextDeletionTime *time.Time

	for _, reminder := range remindersToDelete {
		if reminder.NextFireUTC != nil {
			continue
		}

		deletionTime := now.Add(GarbageCollectionDelay)
		if nextDeletionTime == nil || deletionTime.Before(*nextDeletionTime) {
			nextDeletionTime = &deletionTime
		}
	}

	if nextDeletionTime == nil {
		if config.IsDebugMode() {
			log.Println("[ENGINE] - No valid deletion times found")
		}
		return
	}

	duration := nextDeletionTime.Sub(now)
	if duration <= 0 {
		// Ready to delete now
		if config.IsDebugMode() {
			log.Println("[ENGINE] - Reminders ready for deletion, processing immediately")
		}
		go func() {
			time.Sleep(10 * time.Millisecond)
			select {
			case <-gc.stopChan:
				return
			default:
				gc.checkAndDeleteReminders()
				gc.scheduleNext()
			}
		}()
		return
	}

	if config.IsDebugMode() {
		log.Printf("[ENGINE] - Next deletion scheduled in %v", duration)
	}

	gc.currentTimer = time.NewTimer(duration)
}

// checkAndDeleteReminders processes and deletes eligible reminders
func (gc *GarbageCollector) checkAndDeleteReminders() {
	remindersToDelete, err := gc.reminderRepo.GetNextsRemindersToDelete()
	if err != nil {
		log.Printf("[ENGINE] - Error fetching reminders to delete: %v", err)
		return
	}

	if len(remindersToDelete) == 0 {
		return
	}

	now := time.Now().UTC()
	deletedCount := 0

	for _, reminder := range remindersToDelete {
		if reminder.NextFireUTC != nil {
			continue
		}

		// Check if enough time has passed since dispatch
		deletionTime := reminder.CreatedAt.Add(GarbageCollectionDelay)
		if now.After(deletionTime) || now.Equal(deletionTime) {
			// Delete the reminder (don't notify scheduler to avoid circular updates)
			err := gc.reminderRepo.Delete(reminder.ID, false)
			if err != nil {
				log.Printf("[ENGINE] - Error deleting reminder %s: %v", reminder.ID, err)
				continue
			}

			if config.IsDebugMode() {
				log.Printf("[ENGINE] - Deleted one-time reminder %s",
					reminder.ID)
			}
			deletedCount++
		}
	}

	if deletedCount > 0 && config.IsDebugMode() {
		log.Printf("[ENGINE] - Deleted %d one-time reminder(s)", deletedCount)
	}
}