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
		input eventhandler.EventArgs
		// For when `Handle`() returns an error.
		expectError bool
		// For when the event is discarded.
		expectNoOp bool
		// `expectedOutput` is mutually exclusive with `expectError` and `expectNoOp`.
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
		"WronglyFormattedStartTime": {
			input:       eventhandler.EventArgs{"DriverA", "01:00:45", "02:00", "25.5"},
			expectError: true,
		},
		"WronglyFormattedStopTime": {
			input:       eventhandler.EventArgs{"DriverA", "01:00", "02:00:45", "25.5"},
			expectError: true,
		},
		"StartTimeNotBeforeStopTime": {
			input:       eventhandler.EventArgs{"DriverA", "02:00", "01:00", "25.5"},
			expectError: true,
		},
		"StartTimeSameAsStopTime": {
			input:       eventhandler.EventArgs{"DriverA", "01:00", "01:00", "25.5"},
			expectError: true,
		},
		"WronglyFormattedMileage": {
			input:       eventhandler.EventArgs{"DriverA", "01:00", "02:00", "Ten"},
			expectError: true,
		},
		"MileageFormattedAsInt": {
			input: eventhandler.EventArgs{"DriverA", "01:00", "02:00", "25"},
			expectedOutput: eventstore.VisitableEntity{
				DriverFirstName:     "DriverA",
				TotalDurationDriven: 1 * time.Hour,
				TotalMilesDriven:    25.0,
			},
		},
		"CorrectArgs": {
			input: eventhandler.EventArgs{"DriverA", "01:00", "02:00", "25.5"},
			expectedOutput: eventstore.VisitableEntity{
				DriverFirstName:     "DriverA",
				TotalDurationDriven: 1 * time.Hour,
				TotalMilesDriven:    25.5,
			},
		},
		"DiscardTripWithSpeedLessThan5Mph": {
			input:      eventhandler.EventArgs{"DriverA", "01:00", "02:00", "4.9"},
			expectNoOp: true,
		},
		"DiscardTripWithSpeedGreaterThan100Mph": {
			input:      eventhandler.EventArgs{"DriverA", "01:00", "02:00", "100.1"},
			expectNoOp: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			es := eventstore.New()
			teh := &trip.EventHandler{}

			err := teh.Handle(tc.input, es)
			switch {
			case tc.expectError && err != nil:
				return
			case !tc.expectError && err != nil:
				t.Fatalf("expected: no error, got: %v", err)
			case tc.expectError && err == nil:
				t.Fatalf("expected: error, got: no error")
			}

			r := eventstore.NewRecorder()
			es.Visit(r)

			if tc.expectNoOp {
				if len(r.Entities) == 0 {
					return
				} else {
					t.Fatalf("expected: 0 VisitableEntities in EventStore, got %#v", r.Entities)
				}
			}

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
