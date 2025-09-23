package engine

import (
	"fmt"
	"log"

	"github.com/ericp/chronos-bot-reminder/internal/config"
	"github.com/ericp/chronos-bot-reminder/internal/database/models"
)

// Dispatcher interface defines how reminders are sent to different destinations
type Dispatcher interface {
	Dispatch(reminder *models.Reminder, destination *models.ReminderDestination) error
	GetSupportedType() models.DestinationType
}

// DispatcherRegistry manages all available dispatchers
type DispatcherRegistry struct {
	dispatchers map[models.DestinationType]Dispatcher
}

// NewDispatcherRegistry creates a new dispatcher registry
func NewDispatcherRegistry() *DispatcherRegistry {
	return &DispatcherRegistry{
		dispatchers: make(map[models.DestinationType]Dispatcher),
	}
}

// RegisterDispatcher registers a new dispatcher for a specific destination type
func (dr *DispatcherRegistry) RegisterDispatcher(dispatcher Dispatcher) {
	dr.dispatchers[dispatcher.GetSupportedType()] = dispatcher
}

// DispatchReminder dispatches a reminder to all its destinations
func (dr *DispatcherRegistry) DispatchReminder(reminder *models.Reminder) error {
	if len(reminder.Destinations) == 0 {
		return fmt.Errorf("reminder %s has no destinations", reminder.ID)
	}

	var errors []error

	for _, destination := range reminder.Destinations {
		dispatcher, exists := dr.dispatchers[destination.Type]
		if !exists {
			log.Printf("[DISPATCHER] - No dispatcher found for type %s, skipping", destination.Type)
			continue
		}

		if err := dispatcher.Dispatch(reminder, &destination); err != nil {
			log.Printf("[DISPATCHER] - Error dispatching to %s: %v", destination.Type, err)
			errors = append(errors, err)
			continue
		}

		if config.IsDebugMode() {
			log.Printf("[DISPATCHER] - [DEBUG] Dispatched reminder %s to %s", reminder.ID, destination.Type)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to dispatch to %d destinations", len(errors))
	}

	return nil
}

// GetDispatcher returns a dispatcher for a specific type
func (dr *DispatcherRegistry) GetDispatcher(destinationType models.DestinationType) (Dispatcher, bool) {
	dispatcher, exists := dr.dispatchers[destinationType]
	return dispatcher, exists
}