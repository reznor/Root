package mathutils_test

import (
	"math"
	"testing"
	"time"

	"root.challenge/mathutils"
)

// Convenience wrapper around floatsAreEqual64() to avoid clients performing this typecast
// all over their code.
func floatsAreEqual32(f1, f2 float32) bool {
	return floatsAreEqual64(float64(f1), float64(f2))
}

// See https://stackoverflow.com/questions/588004/is-floating-point-math-broken
// for more.
func floatsAreEqual64(f1, f2 float64) bool {
	const floatErrorTolerance float64 = 1e-3

	if math.IsInf(f1, +1) && math.IsInf(f2, +1) {
		return true
	}

	return math.Abs(f1-f2) < floatErrorTolerance
}
func TestRoundFloat64ToInt64(t *testing.T) {
	tests := map[string]struct {
		input          float64
		expectedOutput int64
	}{
		"Zero":              {input: 0.0, expectedOutput: 0},
		"PositiveRoundDown": {input: 9.4, expectedOutput: 9},
		"PositiveRoundUp":   {input: 9.6, expectedOutput: 10},
		"PositivePoint5":    {input: 9.5, expectedOutput: 10},
		"NegativeRoundDown": {input: -9.6, expectedOutput: -10},
		"NegativeRoundUp":   {input: -9.4, expectedOutput: -9},
		"NegativePoint5":    {input: -9.5, expectedOutput: -10},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualOutput := mathutils.RoundFloat64ToInt64(tc.input)
			if actualOutput != tc.expectedOutput {
				t.Fatalf("expected: %v, got: %v", tc.expectedOutput, actualOutput)
			}
		})
	}
}

func TestComputeSpeedMph32(t *testing.T) {
	tests := map[string]struct {
		inputMileage   float32
		inputDuration  time.Duration
		expectedOutput float32
	}{
		"SixtyMilesInOneHour":        {inputMileage: 60.0, inputDuration: 1 * time.Hour, expectedOutput: 60.0},
		"ThirtyMilesInThirtyMinutes": {inputMileage: 30.0, inputDuration: 30 * time.Minute, expectedOutput: 60.0},
		"OneMileInOneMinute":         {inputMileage: 1.0, inputDuration: 1 * time.Minute, expectedOutput: 60.0},
		"PointOneMileInOneSecond":    {inputMileage: 0.1, inputDuration: 1 * time.Second, expectedOutput: 360.0},
		"ZeroMileage":                {inputMileage: 0.0, inputDuration: 1 * time.Second, expectedOutput: 0.0},
		"ZeroDuration":               {inputMileage: 100.0, inputDuration: 0 * time.Second, expectedOutput: float32(math.Inf(+1))},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualOutput := mathutils.ComputeSpeedMph32(tc.inputMileage, tc.inputDuration)
			if !floatsAreEqual32(actualOutput, tc.expectedOutput) {
				t.Fatalf("expected: %v, got: %v", tc.expectedOutput, actualOutput)
			}
		})
	}
}

func TestComputeSpeedMph64(t *testing.T) {
	tests := map[string]struct {
		inputMileage   float64
		inputDuration  time.Duration
		expectedOutput float64
	}{
		"SixtyMilesInOneHour":        {inputMileage: 60.0, inputDuration: 1 * time.Hour, expectedOutput: 60.0},
		"ThirtyMilesInThirtyMinutes": {inputMileage: 30.0, inputDuration: 30 * time.Minute, expectedOutput: 60.0},
		"OneMileInOneMinute":         {inputMileage: 1.0, inputDuration: 1 * time.Minute, expectedOutput: 60.0},
		"PointOneMileInOneSecond":    {inputMileage: 0.1, inputDuration: 1 * time.Second, expectedOutput: 360.0},
		"ZeroMileage":                {inputMileage: 0.0, inputDuration: 1 * time.Second, expectedOutput: 0.0},
		"ZeroDuration":               {inputMileage: 100.0, inputDuration: 0 * time.Second, expectedOutput: float64(math.Inf(+1))},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualOutput := mathutils.ComputeSpeedMph64(tc.inputMileage, tc.inputDuration)
			if !floatsAreEqual64(actualOutput, tc.expectedOutput) {
				t.Fatalf("expected: %v, got: %v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
