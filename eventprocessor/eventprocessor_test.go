package eventprocessor_test

import (
	"fmt"
	"reflect"
	"testing"

	"root.challenge/eventhandler"
	"root.challenge/eventprocessor"
	"root.challenge/eventstore"
	"root.challenge/input"
)

const testEventType eventhandler.EventType = "TestEvent"

// `testEventHandler` is an implementation of `eventhandler.Interface` useful for tests because its
// behavior can be controlled, and its invocations can be inspected.
type testEventHandler struct {
	handleShouldReturnError bool
	recordedEventArgs       []eventhandler.EventArgs
}

func (teh *testEventHandler) Handle(eventArgs eventhandler.EventArgs, eventStore *eventstore.EventStore) error {
	if teh.handleShouldReturnError {
		return fmt.Errorf("testEventHandler Handle error")
	}

	teh.recordedEventArgs = append(teh.recordedEventArgs, eventArgs)

	return nil
}

// Configure and initialize `testEventHandler` before another round of tests commences.
func (teh *testEventHandler) setup(handleShouldReturnError bool) {
	teh.handleShouldReturnError = handleShouldReturnError
	teh.recordedEventArgs = make([]eventhandler.EventArgs, 0)
}

// Return `testEventHandler` to a pristine state after a round of tests completes.
func (teh *testEventHandler) teardown() {
	teh.recordedEventArgs = nil
	teh.handleShouldReturnError = false
}

func TestTripEventHandler(t *testing.T) {
	tests := map[string]struct {
		input                   []*input.EventEnvelope
		handleShouldReturnError bool
		// For when `Process`() returns errors.
		numExpectedErrors int
		expectedOutput    []eventhandler.EventArgs
	}{
		"EmptyInput": {
			input:          []*input.EventEnvelope{},
			expectedOutput: []eventhandler.EventArgs{},
		},
		"InputWithEmptyLines": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("                  ")),
				input.NewEventEnvelopeForBody(input.NewEventFromString("    ")),
				input.NewEventEnvelopeForBody(input.NewEventFromString("            ")),
			},
			expectedOutput: []eventhandler.EventArgs{},
		},
		"SingleEventWithBody": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent1Arg1")),
			},
			expectedOutput: []eventhandler.EventArgs{
				{"TestEvent1Arg1"},
			},
		},
		"MultipleEventsWithBody": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent1Arg1")),
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent2Arg1")),
			},
			expectedOutput: []eventhandler.EventArgs{
				{"TestEvent1Arg1"},
				{"TestEvent2Arg1"},
			},
		},
		"EventWithMultipleArgs": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent1Arg1 TestEvent1Arg2")),
			},
			expectedOutput: []eventhandler.EventArgs{
				{"TestEvent1Arg1", "TestEvent1Arg2"},
			},
		},
		"UnrecognizedEvent": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("UnrecognizedEvent Arg1")),
			},
			numExpectedErrors: 1,
			expectedOutput:    []eventhandler.EventArgs{},
		},
		"HandleReturnsError": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent1Arg1")),
			},
			handleShouldReturnError: true,
			numExpectedErrors:       1,
			expectedOutput:          []eventhandler.EventArgs{},
		},
		"SingleEventWithError": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForError(fmt.Errorf("Some input error.")),
			},
			numExpectedErrors: 1,
			expectedOutput:    []eventhandler.EventArgs{},
		},
		"MalformedEmptyEventWithoutBodyNorError": {
			input: []*input.EventEnvelope{
				{},
			},
			numExpectedErrors: 1,
			expectedOutput:    []eventhandler.EventArgs{},
		},
		"MixOfAllEvents": {
			input: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent1Arg1 TestEvent1Arg2")),
				{},
				input.NewEventEnvelopeForError(fmt.Errorf("Some input error.")),
				input.NewEventEnvelopeForBody(input.NewEventFromString("                            ")),
				input.NewEventEnvelopeForBody(input.NewEventFromString("TestEvent TestEvent2Arg1 TestEvent2Arg2")),
				input.NewEventEnvelopeForBody(input.NewEventFromString("UnrecognizedEvent Arg1")),
			},
			numExpectedErrors: 3,
			expectedOutput: []eventhandler.EventArgs{
				{"TestEvent1Arg1", "TestEvent1Arg2"},
				{"TestEvent2Arg1", "TestEvent2Arg2"},
			},
		},
	}

	// Register `testEventHandler` one time for the entire test case, akin to how packages are initialized exactly
	// once at load time.
	teh := &testEventHandler{}
	eventhandler.GlobalRegistry().RegisterEventHandler(testEventType, teh)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			teh.setup(tc.handleShouldReturnError)
			defer teh.teardown()

			// Load `eventC` up with all the intended `input.EventEnvelope`s up-front to avoid any unnecessarily-distracting
			// concurrency primitives in test code.
			eventC := make(chan *input.EventEnvelope, len(tc.input))
			for _, eventEnvelope := range tc.input {
				eventC <- eventEnvelope
			}
			close(eventC)

			// Record the `error`s emitted by `Process`() to make the test easier to debug.
			actualErrors := make([]error, 0)
			for err := range eventprocessor.New().Process(eventC, eventstore.New()) {
				actualErrors = append(actualErrors, err)
			}

			if len(actualErrors) != tc.numExpectedErrors {
				t.Fatalf("expected: %d errors, got %d errors (%#v)",
					tc.numExpectedErrors, len(actualErrors), actualErrors)
			}

			// Compare the invocations of `testEventHandler` to avoid having to peer deeper into the workings of
			// `eventhandler.Interface` (by inspecting what it stored in `eventstore.EventStore`) -- that level of
			// invasive testing is better-suited for the unit tests of the `eventhandler` package itself (where
			// such an inspection would not be considered invasive).
			actualOutput := teh.recordedEventArgs
			if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
				t.Fatalf("expected: %#v, got: %#v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
