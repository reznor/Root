package driver

import (
	"fmt"

	"root.challenge/eventhandler"
	"root.challenge/eventstore"
)

const driverEventType eventhandler.EventType = "Driver"

type driverEventHandler struct{}

func init() {
	eventhandler.RegisterEventHandler(driverEventType, &driverEventHandler{})
}

func (deh *driverEventHandler) Handle(eventArgs eventhandler.EventArgs, eventStore *eventstore.EventStore) error {
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
