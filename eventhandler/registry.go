package eventhandler

import (
	"fmt"
	"sync"

	"root.challenge/eventstore"
)

type EventType string
type EventArgs []string

type Interface interface {
	Handle(EventArgs, *eventstore.EventStore) error
}

var (
	registryMu sync.RWMutex
	registry   = make(map[EventType]Interface)
)

func RegisterEventHandler(eventType EventType, eventHandler Interface) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if eventHandler == nil {
		panic(fmt.Errorf("nil EventHandler provided for EventType '%s'", eventType))
	}

	if _, dup := registry[eventType]; dup {
		panic(fmt.Errorf("EventHandler for EventType '%s' is already registered", eventType))
	}

	registry[eventType] = eventHandler
}

func GetHandlerForEvent(eventType EventType) (Interface, error) {
	registryMu.RLock()
	eventHandler, ok := registry[eventType]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no EventHandler registered for EventType '%s'", eventType)
	}

	return eventHandler, nil
}
