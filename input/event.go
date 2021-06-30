package input

type Event string

type EventEnvelope struct {
	Err error

	// EXPLAIN: Why this is a pointer (for future compatibility, as the type can change and become meatier)
	Body *Event
}

func NewEventEnvelopeForError(err error) *EventEnvelope {
	return &EventEnvelope{
		Err: err,
	}
}

func NewEventEnvelopeForBody(inputEvent *Event) *EventEnvelope {
	return &EventEnvelope{
		Body: inputEvent,
	}
}
