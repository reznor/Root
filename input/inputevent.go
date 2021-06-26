package input

type InputEvent string

type InputEventEnvelope struct {
	Err error

	// EXPLAIN: Why this is a pointer (for future compatibility, as the type can change and become meatier)
	Body *InputEvent
}

func NewInputEventEnvelopeWithError(err error) *InputEventEnvelope {
	return &InputEventEnvelope{
		Err: err,
	}
}

func NewInputEventEnvelopeWithBody(inputEvent *InputEvent) *InputEventEnvelope {
	return &InputEventEnvelope{
		Body: inputEvent,
	}
}
