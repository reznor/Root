package eventhandler

import (
	"fmt"
	"sync"
)

// Registry stores the concrete `eventhandler.Interface` implementations registered with the system.
type Registry struct {
	mutex    sync.RWMutex
	registry map[EventType]Interface
}

// NewRegistry creates a new `Registry` object.
//
// This should typically only be used for testing or for niche scenarios that require isolation in the
// face of multi-tenancy -- the expected production scenario is to just use `GlobalRegistry`().
func NewRegistry() *Registry {
	return &Registry{
		registry: make(map[EventType]Interface),
	}
}

var globalRegistry = NewRegistry()

// GlobalRegistry provides access to the one shared global `Registry` for the entire system to use.
func GlobalRegistry() *Registry {
	return globalRegistry
}

// RegisterEventHandler registers an `eventhandler.Interface` implementation to be invoked at runtime
// for all events of a particular `EventType`.
//
// It returns `error` if a nil implementation is provided, or if this method has already previously been
// called for the same `EventType`.
//
// The call to this method is typically expected to be made in each handler package's `init`() function.
func (r *Registry) RegisterEventHandler(eventType EventType, eventHandler Interface) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if eventHandler == nil {
		return fmt.Errorf("nil EventHandler provided for EventType '%s'", eventType)
	}

	if _, dup := r.registry[eventType]; dup {
		return fmt.Errorf("EventHandler for EventType '%s' is already registered", eventType)
	}

	r.registry[eventType] = eventHandler
	return nil
}

// GetHandlerForEvent returns the previously-registered `eventhandler.Interface` implementation for
// a particular `EventType`.
//
// It returns `error` if no prior call to `RegisterEventHandler`() was made for the `EventType`.
func (r *Registry) GetHandlerForEvent(eventType EventType) (Interface, error) {
	r.mutex.RLock()
	eventHandler, ok := r.registry[eventType]
	r.mutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no EventHandler registered for EventType '%s'", eventType)
	}

	return eventHandler, nil
}
