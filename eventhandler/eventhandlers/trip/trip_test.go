package trip_test

import (
	"testing"
	"time"

	"root.challenge/eventhandler"
	"root.challenge/eventhandler/eventhandlers/trip"
	"root.challenge/eventstore"
)

func TestTripEventHandler(t *testing.T) {
	tests := map[string]struct {
		input       eventhandler.EventArgs
		expectError bool
		// expectedOutput is mutually exclusive with expectedError.
		expectedOutput eventstore.VisitableEntity
	}{
		"TooFewArgs": {
			input:       eventhandler.EventArgs{},
			expectError: true,
		},
		"TooManyArgs": {
			input:       eventhandler.EventArgs{"DriverA", "01:00", "02:00", "25.5", "SomeUnwantedInfo"},
			expectError: true,
		},
		"CorrectArgs": {
			input: eventhandler.EventArgs{"DriverA", "01:00", "02:00", "25.5"},
			expectedOutput: eventstore.VisitableEntity{
				DriverFirstName:     "DriverA",
				TotalDurationDriven: 1 * time.Hour,
				TotalMilesDriven:    25.5,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			es := eventstore.New()
			teh := &trip.EventHandler{}

			err := teh.Handle(tc.input, es)
			switch {
			case err != nil && tc.expectError:
				return
			case err != nil && !tc.expectError:
				t.Fatalf("expected: no error, got: %v", err)
			case err == nil && tc.expectError:
				t.Fatalf("expected: error, got: no error")
			}

			r := eventstore.NewRecorder()
			es.Visit(r)

			if len(r.Entities) != 1 {
				t.Fatalf("expected exactly 1 VisitableEntity in EventStore, got %v", len(r.Entities))
			}

			actualOutput := r.Entities[0]
			if actualOutput != tc.expectedOutput {
				t.Fatalf("expected: %v, got: %v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
