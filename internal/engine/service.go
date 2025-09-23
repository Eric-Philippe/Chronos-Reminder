package engine

import (
	"context"
	"log"
	"sync"

	"github.com/ericp/chronos-bot-reminder/internal/database"
	"github.com/ericp/chronos-bot-reminder/internal/database/repositories"
)

// SchedulerService manages the complete scheduling system
type SchedulerService struct {
	Scheduler          *Scheduler
	DispatcherRegistry *DispatcherRegistry
	NotifyingRepo      *NotifyingReminderRepository
}

var (
	schedulerService *SchedulerService
	schedulerCtx     context.Context
	schedulerCancel  context.CancelFunc
	schedulerMutex   sync.Mutex
)

// GetSchedulerService returns a singleton scheduler service
func GetSchedulerService() *SchedulerService {
	schedulerMutex.Lock()
	defer schedulerMutex.Unlock()

	if schedulerService == nil {
		repos := database.GetRepositories()
		if repos == nil {
			log.Fatalf("[ENGINE] - ❌ Cannot get repositories")
		}

		schedulerService = NewSchedulerService(repos.Reminder)
	}

	return schedulerService
}

// StartSchedulerService initializes and starts the scheduler service
func StartSchedulerService() {
	schedulerMutex.Lock()
	defer schedulerMutex.Unlock()

	if schedulerService == nil {
		repos := database.GetRepositories()
		if repos == nil {
			log.Fatalf("[ENGINE] - ❌ Cannot get repositories")
		}

		schedulerService = NewSchedulerService(repos.Reminder)
	}

	if schedulerService.Scheduler.IsRunning() {
		log.Println("[ENGINE] - Already running")
		return
	}

	// Create context for the scheduler
	schedulerCtx, schedulerCancel = context.WithCancel(context.Background())

	// Initialize the repository notifier
	initializeRepositoryNotifier(schedulerService)

	// Start the scheduler
	schedulerService.Scheduler.Start(schedulerCtx)
	log.Println("[ENGINE] - ✅ Scheduler started")
}

// StopSchedulerService gracefully stops the scheduler service
func StopSchedulerService() {
	schedulerMutex.Lock()
	defer schedulerMutex.Unlock()

	if schedulerService != nil && schedulerService.Scheduler.IsRunning() {
		schedulerService.Scheduler.Stop()
		log.Println("[ENGINE] - ✅ Scheduler stopped")
	}

	if schedulerCancel != nil {
		schedulerCancel()
		schedulerCancel = nil
	}

	// Reset the service
	schedulerService = nil
}

// GetNotifyingRepository returns the notifying repository for use by other parts of the application
func GetNotifyingRepository() *NotifyingReminderRepository {
	service := GetSchedulerService()
	if service == nil {
		return nil
	}
	return service.NotifyingRepo
}

// InitializeRepositoryNotifier sets up the scheduler notifier in the base repository
// This should be called after the scheduler service is created
func InitializeRepositoryNotifier() {
	service := GetSchedulerService()
	if service == nil {
		log.Printf("[ENGINE] - ⚠️ Cannot initialize repository notifier: scheduler service not available")
		return
	}

	initializeRepositoryNotifier(service)
}

// initializeRepositoryNotifier is the internal implementation that doesn't call GetSchedulerService
// to avoid deadlocks when called from within StartSchedulerService
func initializeRepositoryNotifier(service *SchedulerService) {
	// Get the base repositories
	repos := database.GetRepositories()
	if repos == nil {
		log.Printf("[ENGINE] - ⚠️ Cannot initialize repository notifier: repositories not available")
		return
	}

	// Set the scheduler notifier in the base reminder repository
	if reminderRepo, ok := repos.Reminder.(interface{ SetScheduler(repositories.SchedulerNotifier) }); ok {
		reminderRepo.SetScheduler(service.Scheduler)
		log.Printf("[ENGINE] - ✅ Repository notifier initialized")
	} else {
		log.Printf("[ENGINE] - ⚠️ Reminder repository does not support scheduler notification")
	}
}

// GetSchedulerAwareReminderRepository returns either the notifying repository if scheduler is running,
// or the base repository if scheduler is not available
func GetSchedulerAwareReminderRepository() repositories.ReminderRepository {
	// Try to get the notifying repository first
	notifyingRepo := GetNotifyingRepository()
	if notifyingRepo != nil {
		return notifyingRepo
	}

	// Fallback to base repository
	repos := database.GetRepositories()
	if repos != nil {
		return repos.Reminder
	}

	return nil
}

// IsSchedulerNotificationEnabled returns true if the scheduler is running and can receive notifications
func IsSchedulerNotificationEnabled() bool {
	service := GetSchedulerService()
	return service != nil && service.Scheduler.IsRunning()
}

// NewSchedulerService creates a new complete scheduler service with all dispatchers registered
func NewSchedulerService(reminderRepo repositories.ReminderRepository) *SchedulerService {
	// Create dispatcher registry
	dispatcherRegistry := NewDispatcherRegistry()
	
	// Register all dispatchers
	dispatcherRegistry.RegisterDispatcher(NewDiscordDMDispatcher())
	dispatcherRegistry.RegisterDispatcher(NewWebhookDispatcher())
	dispatcherRegistry.RegisterDispatcher(NewDiscordChannelDispatcher())
	
	// Create scheduler
	scheduler := NewScheduler(reminderRepo, dispatcherRegistry)
	
	// Create notifying repository that will update the scheduler
	notifyingRepo := NewNotifyingReminderRepository(reminderRepo, scheduler)
	
	return &SchedulerService{
		Scheduler:          scheduler,
		DispatcherRegistry: dispatcherRegistry,
		NotifyingRepo:      notifyingRepo,
	}
}
