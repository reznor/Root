package output_test

import (
	"reflect"
	"testing"
	"time"

	"root.challenge/eventstore"
	"root.challenge/output"
)

func TestReportGenerator(t *testing.T) {
	tests := map[string]struct {
		input          []*eventstore.VisitableEntity
		expectedOutput output.GeneratedReport
	}{
		"EmptyInput": {
			input:          []*eventstore.VisitableEntity{},
			expectedOutput: output.GeneratedReport{},
		},
		"SimpleInput": {
			input: []*eventstore.VisitableEntity{
				{DriverFirstName: "DriverA", TotalDurationDriven: 1 * time.Hour, TotalMilesDriven: 30.0},
				{DriverFirstName: "DriverB", TotalDurationDriven: 3 * time.Hour, TotalMilesDriven: 45.0},
				{DriverFirstName: "DriverC", TotalDurationDriven: 20 * time.Minute, TotalMilesDriven: 40.0},
				{DriverFirstName: "DriverD", TotalDurationDriven: 15 * time.Minute, TotalMilesDriven: 15.0},
			},
			expectedOutput: output.GeneratedReport{
				"DriverB: 45 miles @ 15 mph",
				"DriverC: 40 miles @ 120 mph",
				"DriverA: 30 miles @ 30 mph",
				"DriverD: 15 miles @ 60 mph",
			},
		},
		"AlreadySortedInput": {
			input: []*eventstore.VisitableEntity{
				{DriverFirstName: "DriverA", TotalDurationDriven: 3 * time.Hour, TotalMilesDriven: 45.0},
				{DriverFirstName: "DriverB", TotalDurationDriven: 20 * time.Minute, TotalMilesDriven: 40.0},
				{DriverFirstName: "DriverC", TotalDurationDriven: 1 * time.Hour, TotalMilesDriven: 30.0},
				{DriverFirstName: "DriverD", TotalDurationDriven: 15 * time.Minute, TotalMilesDriven: 15.0},
			},
			expectedOutput: output.GeneratedReport{
				"DriverA: 45 miles @ 15 mph",
				"DriverB: 40 miles @ 120 mph",
				"DriverC: 30 miles @ 30 mph",
				"DriverD: 15 miles @ 60 mph",
			},
		},
		"ReverseSortedInput": {
			input: []*eventstore.VisitableEntity{
				{DriverFirstName: "DriverA", TotalDurationDriven: 15 * time.Minute, TotalMilesDriven: 15.0},
				{DriverFirstName: "DriverB", TotalDurationDriven: 1 * time.Hour, TotalMilesDriven: 30.0},
				{DriverFirstName: "DriverC", TotalDurationDriven: 20 * time.Minute, TotalMilesDriven: 40.0},
				{DriverFirstName: "DriverD", TotalDurationDriven: 3 * time.Hour, TotalMilesDriven: 45.0},
			},
			expectedOutput: output.GeneratedReport{
				"DriverD: 45 miles @ 15 mph",
				"DriverC: 40 miles @ 120 mph",
				"DriverB: 30 miles @ 30 mph",
				"DriverA: 15 miles @ 60 mph",
			},
		},
		"DriverWithNoTripEvents": {
			input: []*eventstore.VisitableEntity{
				{DriverFirstName: "DriverA", TotalDurationDriven: 0 * time.Second, TotalMilesDriven: 0.0},
				{DriverFirstName: "DriverB", TotalDurationDriven: 1 * time.Hour, TotalMilesDriven: 30.0},
			},
			expectedOutput: output.GeneratedReport{
				"DriverB: 30 miles @ 30 mph",
				"DriverA: 0 miles",
			},
		},
		"RoundingFractionalMileageAndSpeed": {
			input: []*eventstore.VisitableEntity{
				{DriverFirstName: "DriverA", TotalDurationDriven: 11 * time.Minute, TotalMilesDriven: 67.4},
				{DriverFirstName: "DriverB", TotalDurationDriven: 1*time.Hour + 39*time.Minute, TotalMilesDriven: 72.9},
			},
			expectedOutput: output.GeneratedReport{
				"DriverB: 73 miles @ 44 mph",
				"DriverA: 67 miles @ 368 mph",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rg := output.NewReportGenerator()
			for _, ve := range tc.input {
				rg.Visit(ve)
			}

			actualOutput := rg.Generate()
			if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
				t.Fatalf("expected: %#v, got: %#v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
