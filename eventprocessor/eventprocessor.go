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

// EventProcessor is the central entity responsible for orchestrating the processing of events that
// come into the system.
//
// It is expected that only one of these is typically created in the system, with multiple disparate
// calls to `Process`() (for example, if there are multiple input sources to process concurrently,
// each potentially using a separate `eventstore.EventStore`) made as needed.
type EventProcessor struct{}

// New creates a new `EventProcessor`.
func New() *EventProcessor {
	return &EventProcessor{}
}

// Process receives a stream of `input.EventEnvelope` objects and processes them in the background as
// appropriate (storing relevant information in the passed-in `eventstore.EventStore` as a by-product)
// while emitting any encountered `error`s on the returned channel.
//
// It is worth noting that all the state required to process an input stream is provided by callers
// here -- and not in `New`() -- to allow a single `EventProcessor` to own the processing of every input
// stream entering the system; that, in turn, makes this method stateless, and thus amenable to being
// hosted on serverless/FaaS technology stacks.
func (ep *EventProcessor) Process(eventC <-chan *input.EventEnvelope, eventStore *eventstore.EventStore) <-chan error {
	errC := make(chan error)

	go func() {
		defer close(errC)

		for eventEnvelope := range eventC {
			if eventEnvelope.Err != nil {
				errC <- fmt.Errorf("error retrieving next EventEnvelope from channel: %w", eventEnvelope.Err)
				continue
			}

			if eventEnvelope.Body == nil {
				errC <- fmt.Errorf("retrieved malformed EventEnvelope with nil Body but nil Err as well")
				continue
			}

			parsedEvent := strings.Fields(string(*eventEnvelope.Body))
			if len(parsedEvent) == 0 {
				// Skip over empty events.
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
