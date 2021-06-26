package eventprocessor

import (
	"log"

	_ "root.challenge/eventhandler/eventhandlers/driver"
	_ "root.challenge/eventhandler/eventhandlers/trip"
	"root.challenge/input"
)

type EventProcessor struct {
	eventC <-chan *input.EventEnvelope
}

func New(eventC <-chan *input.EventEnvelope) *EventProcessor {
	return &EventProcessor{
		eventC: eventC,
	}
}

func (ep *EventProcessor) Process() {
	for inputEventEnvelope := range ep.eventC {
		if inputEventEnvelope.Err != nil {
			log.Fatalf("Error event --> %s", inputEventEnvelope.Err)
			// EXPLAIN: Choose to stop all processing instead of using `continue` to follow the principle of fail-early-fail-hard; better to not mess up any state once the input stream has gotten into an error state, but this constraint can change as the requirements/functionality changes.
			break
		}

		log.Printf("Read event ==> %s", *inputEventEnvelope.Body)
	}
}
