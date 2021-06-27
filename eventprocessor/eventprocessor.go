package eventprocessor

import (
	"fmt"
	"strings"

	"root.challenge/eventhandler"
	_ "root.challenge/eventhandler/eventhandlers/driver"
	_ "root.challenge/eventhandler/eventhandlers/trip"
	"root.challenge/eventstore"
	"root.challenge/input"
)

type EventProcessor struct{}

func New() *EventProcessor {
	return &EventProcessor{}
}

func (ep *EventProcessor) Process(eventC <-chan *input.EventEnvelope, eventStore *eventstore.EventStore) <-chan error {
	errC := make(chan error)

	go func() {
		defer close(errC)

		for inputEventEnvelope := range eventC {
			if inputEventEnvelope.Err != nil {
				errC <- fmt.Errorf("error retrieving next event from channel: %w", inputEventEnvelope.Err)
				// EXPLAIN: Choose to stop all processing instead of using `continue` to follow the principle of fail-early-fail-hard; better to not mess up any state once the input stream has gotten into an error state, but this constraint can change as the requirements/functionality changes.
				// TODO: Consider making this a continue as well? Check to see when all the channel is closed, and what kinds of errors are emitted over it.
				break
			}

			// Now that there's no infrastructural error in retrieving the event itself, it's alright to skip over
			// malformed/unrecognized events here onwards, since each stands independently, and these types of errors
			// are minor enough to not warrant giving up altogether and terminating.
			if inputEventEnvelope.Body == nil {
				errC <- fmt.Errorf("retrieved malformed EventEnvelope with nil Body but nil Err as well")
				continue
			}

			parsedEvent := strings.Fields(string(*inputEventEnvelope.Body))
			if len(parsedEvent) == 0 {
				continue
			}
			eventType := eventhandler.EventType(parsedEvent[0])
			eventArgs := eventhandler.EventArgs(parsedEvent[1:])

			eventHandler, err := eventhandler.GetHandlerForEvent(eventType)
			if err != nil {
				errC <- fmt.Errorf("error retrieving handler for eventType %s: %w", eventType, err)
				continue
			}

			if err := eventHandler.Handle(eventArgs, eventStore); err != nil {
				errC <- fmt.Errorf("error handling eventType %s with args %v: %w",
					eventType, eventArgs, err)
				continue
			}
		}
	}()

	return errC
}
