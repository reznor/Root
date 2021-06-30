package input

// Event represents an input event entering the system.
//
// While it's a simple `string` for now, having a type definition allows it to grow to be something
// more complex and structured later with minimal change across the codebase.
type Event string

// NewEventFromString creates a pointer to an `Event` from the passed-in `string`.
func NewEventFromString(s string) *Event {
	e := Event(s)
	return &e
}

// EventEnvelope provides an encapsulation for every `Event` that is provided as input to the system.
//
// Specifically, it isolates the processing of input `Event`s by the system from errors in the input
// stream (that may contain errors as it reads from its source).
type EventEnvelope struct {
	// `Err` and `Body` are mutually exclusive.
	Err error

	// Body is a pointer for future compatibility -- the type definition of `Event` can change and become
	// meatier, and this allows that change to happen in the codebase in a minimally disruptive manner.
	Body *Event
}

// NewEventEnvelopeForError is a helper to generate an `EventEnvelope` that contains an `error`.
func NewEventEnvelopeForError(err error) *EventEnvelope {
	return &EventEnvelope{
		Err: err,
	}
}

// NewEventEnvelopeForis a helper to generate an `EventEnvelope` that contains an `Event`.
func NewEventEnvelopeForBody(event *Event) *EventEnvelope {
	return &EventEnvelope{
		Body: event,
	}
}
