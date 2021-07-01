package driver

import (
	"fmt"

	"root.challenge/eventhandler"
	"root.challenge/eventstore"
)

const eventType eventhandler.EventType = "Driver"

// EventHandler is an implementation of `eventhandler.Interface` for the "Driver" EventType.
type EventHandler struct{}

func init() {
	if err := eventhandler.GlobalRegistry().RegisterEventHandler(eventType, &EventHandler{}); err != nil {
		panic(err)
	}
}

// Conforms to `eventhandler.Interface`.
func (eh *EventHandler) Handle(eventArgs eventhandler.EventArgs, eventStore *eventstore.EventStore) error {
	if len(eventArgs) != 1 {
		return fmt.Errorf("expecting exactly 1 arg (first name) to Driver event %v; got %d",
			eventArgs, len(eventArgs))
	}

	if err := eventStore.RegisterDriver(&eventstore.DriverInfo{
		FirstName: eventArgs[0],
	}); err != nil {
		return fmt.Errorf("failed to register Driver event %v with EventStore: %w", eventArgs, err)
	}

	return nil
}
