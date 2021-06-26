package driver

import (
	"fmt"

	"root.challenge/eventhandler"
)

const driverEventType eventhandler.EventType = "Driver"

type driverEventHandler struct{}

func init() {
	eventhandler.RegisterEventHandler(driverEventType, &driverEventHandler{})
}

func (deh *driverEventHandler) CheckPreconditions(eventArgs eventhandler.EventArgs) error {
	if len(eventArgs) == 0 {
		return fmt.Errorf("XXX")
	}

	return nil
}

func (deh *driverEventHandler) Handle(eventArgs eventhandler.EventArgs) error {
	if len(eventArgs) == 0 {
		return fmt.Errorf("XXX")
	}

	return nil
}
