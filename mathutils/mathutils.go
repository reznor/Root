package mathutils

import (
	"math"
	"time"
)

func RoundFloat64ToInt64(f float64) int64 {
	return int64(math.Round(f))
}

func ComputeSpeedMph32(mileage float32, duration time.Duration) float32 {
	return mileage / (float32(duration) / float32(time.Hour))
}

func ComputeSpeedMph64(mileage float64, duration time.Duration) float64 {
	return mileage / (float64(duration) / float64(time.Hour))
}
