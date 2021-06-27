package eventhandler

import (
	"fmt"
	"sync"

	"root.challenge/eventstore"
)

type EventType string
type EventArgs []string

type EventHandler interface {
	Handle(eventArgs EventArgs, eventStore *eventstore.EventStore) error
}

var (
	registryMu sync.RWMutex
	registry   = make(map[EventType]EventHandler)
)

func RegisterEventHandler(eventType EventType, eventHandler EventHandler) {
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

func GetHandlerForEvent(eventType EventType) (EventHandler, error) {
	registryMu.RLock()
	eventHandler, ok := registry[eventType]
	registryMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("no EventHandler registered for EventType '%s'", eventType)
	}

	return eventHandler, nil
}
