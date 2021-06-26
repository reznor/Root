package trip

import (
	"fmt"

	"root.challenge/eventhandler"
)

const tripEventType eventhandler.EventType = "Trip"

type tripEventHandler struct{}

func init() {
	eventhandler.RegisterEventHandler(tripEventType, &tripEventHandler{})
}

func (teh *tripEventHandler) CheckPreconditions(eventArgs eventhandler.EventArgs) error {
	if len(eventArgs) == 0 {
		return fmt.Errorf("XXX")
	}

	return nil
}

func (teh *tripEventHandler) Handle(eventArgs eventhandler.EventArgs) error {
	if len(eventArgs) == 0 {
		return fmt.Errorf("XXX")
	}

	return nil
}
