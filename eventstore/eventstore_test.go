package eventstore_test

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"root.challenge/eventstore"
)

type eventStoreMethodInvoker func(eventStore *eventstore.EventStore, methodParams interface{}) error

func registerDriverInvoker(eventStore *eventstore.EventStore, driverInfo interface{}) error {
	return eventStore.RegisterDriver(driverInfo.(*eventstore.DriverInfo))
}

func recordTripInvoker(eventStore *eventstore.EventStore, tripInfo interface{}) error {
	return eventStore.RecordTrip(tripInfo.(*eventstore.TripInfo))
}

type eventStoreMethodInvocation struct {
	invoker eventStoreMethodInvoker
	params  interface{}
}

func TestEventStore(t *testing.T) {
	tests := map[string]struct {
		input          []eventStoreMethodInvocation
		expectedOutput []eventstore.VisitableEntity
	}{
		"OneDriverNoTrip": {
			input: []eventStoreMethodInvocation{
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverA",
					},
				},
			},
			expectedOutput: []eventstore.VisitableEntity{
				{
					DriverFirstName:     "DriverA",
					TotalDurationDriven: 0 * time.Second,
					TotalMilesDriven:    0.0,
				},
			},
		},
		"NoDriverOneTrip": {
			input: []eventStoreMethodInvocation{
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverA",
						TripDuration:    1 * time.Hour,
						TripMileage:     20.0,
					},
				},
			},
			expectedOutput: []eventstore.VisitableEntity{
				{
					DriverFirstName:     "DriverA",
					TotalDurationDriven: 1 * time.Hour,
					TotalMilesDriven:    20.0,
				},
			},
		},
		"OneDriverOneTrip": {
			input: []eventStoreMethodInvocation{
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverA",
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverA",
						TripDuration:    1 * time.Hour,
						TripMileage:     20.0,
					},
				},
			},
			expectedOutput: []eventstore.VisitableEntity{
				{
					DriverFirstName:     "DriverA",
					TotalDurationDriven: 1 * time.Hour,
					TotalMilesDriven:    20.0,
				},
			},
		},
		"OneTripOneDriver": {
			input: []eventStoreMethodInvocation{
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverA",
						TripDuration:    1 * time.Hour,
						TripMileage:     20.0,
					},
				},
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverA",
					},
				},
			},
			expectedOutput: []eventstore.VisitableEntity{
				{
					DriverFirstName:     "DriverA",
					TotalDurationDriven: 1 * time.Hour,
					TotalMilesDriven:    20.0,
				},
			},
		},
		"OneDriverMultipleTrips": {
			input: []eventStoreMethodInvocation{
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverA",
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverA",
						TripDuration:    1 * time.Hour,
						TripMileage:     20.0,
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverA",
						TripDuration:    10 * time.Minute,
						TripMileage:     5.0,
					},
				},
			},
			expectedOutput: []eventstore.VisitableEntity{
				{
					DriverFirstName:     "DriverA",
					TotalDurationDriven: 1*time.Hour + 10*time.Minute,
					TotalMilesDriven:    25.0,
				},
			},
		},
		"MultipleDriversMultipleTrips": {
			input: []eventStoreMethodInvocation{
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverA",
					},
				},
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverB",
					},
				},
				{
					invoker: registerDriverInvoker,
					params: &eventstore.DriverInfo{
						FirstName: "DriverC",
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverB",
						TripDuration:    1 * time.Hour,
						TripMileage:     20.0,
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverC",
						TripDuration:    10 * time.Minute,
						TripMileage:     0.5,
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverB",
						TripDuration:    45 * time.Minute,
						TripMileage:     50.5,
					},
				},
				{
					invoker: recordTripInvoker,
					params: &eventstore.TripInfo{
						DriverFirstName: "DriverA",
						TripDuration:    25 * time.Minute,
						TripMileage:     35.0,
					},
				},
			},
			expectedOutput: []eventstore.VisitableEntity{
				{
					DriverFirstName:     "DriverA",
					TotalDurationDriven: 25 * time.Minute,
					TotalMilesDriven:    35.0,
				},
				{
					DriverFirstName:     "DriverB",
					TotalDurationDriven: 1*time.Hour + 45*time.Minute,
					TotalMilesDriven:    70.5,
				},
				{
					DriverFirstName:     "DriverC",
					TotalDurationDriven: 10 * time.Minute,
					TotalMilesDriven:    0.5,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			es := eventstore.New()
			for _, invocation := range tc.input {
				invocation.invoker(es, invocation.params)
			}

			r := eventstore.NewRecorder()
			es.Visit(r)

			actualOutput := r.Entities
			// The output of Visit() above is not guaranteed to be in any order, so sort by DriverFirstName to
			// be able to work with something predictable (and comparable to tc.expectedOutput).
			sort.Slice(actualOutput, func(i, j int) bool {
				return actualOutput[i].DriverFirstName < actualOutput[j].DriverFirstName
			})
			if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
				t.Fatalf("expected: %#v, got: %#v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
