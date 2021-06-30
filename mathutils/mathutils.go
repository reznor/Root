package mathutils

import (
	"math"
	"time"
)

// RoundFloat64ToInt64 rounds a `float64` in the manner of `math.Round()`, and returns the result as an `int64`.
func RoundFloat64ToInt64(f float64) int64 {
	return int64(math.Round(f))
}

// ComputeSpeedMph32 computes a `float32` speed in miles/hour given a `float32` mileage and a `time.Duration`.
func ComputeSpeedMph32(mileage float32, duration time.Duration) float32 {
	return mileage / (float32(duration) / float32(time.Hour))
}

// ComputeSpeedMph64 computes a `float64` speed in miles/hour given a `float64` mileage and a `time.Duration`.
func ComputeSpeedMph64(mileage float64, duration time.Duration) float64 {
	return mileage / (float64(duration) / float64(time.Hour))
}
