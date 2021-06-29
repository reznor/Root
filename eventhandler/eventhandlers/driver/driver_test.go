package driver_test

import (
	"testing"
	"time"

	"root.challenge/eventhandler"
	"root.challenge/eventhandler/eventhandlers/driver"
	"root.challenge/eventstore"
)

func TestDriverEventHandler(t *testing.T) {
	tests := map[string]struct {
		input eventhandler.EventArgs
		// For when Handle() returns an error.
		expectError bool
		// expectedOutput is mutually exclusive with expectError.
		expectedOutput eventstore.VisitableEntity
	}{
		"TooFewArgs": {
			input:       eventhandler.EventArgs{},
			expectError: true,
		},
		"TooManyArgs": {
			input:       eventhandler.EventArgs{"DriverA", "SomeLastName"},
			expectError: true,
		},
		"CorrectArgs": {
			input: eventhandler.EventArgs{"DriverA"},
			expectedOutput: eventstore.VisitableEntity{
				DriverFirstName:     "DriverA",
				TotalDurationDriven: 0 * time.Second,
				TotalMilesDriven:    0.0,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			es := eventstore.New()
			deh := &driver.EventHandler{}

			err := deh.Handle(tc.input, es)
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
