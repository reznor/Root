package trip

import (
	"fmt"
	"strconv"
	"time"

	"root.challenge/eventhandler"
	"root.challenge/eventstore"
	"root.challenge/mathutils"
)

const tripEventType eventhandler.EventType = "Trip"

type tripEventHandler struct{}

func init() {
	eventhandler.RegisterEventHandler(tripEventType, &tripEventHandler{})
}

func (teh *tripEventHandler) Handle(eventArgs eventhandler.EventArgs, eventStore *eventstore.EventStore) error {
	if len(eventArgs) != 4 {
		return fmt.Errorf(
			"expecting exactly 4 args (first name, start time, stop time, miles driven) to Trip event %v; got %d",
			eventArgs, len(eventArgs))
	}

	driverFirstName, startTimeStr, stopTimeStr, tripMileageStr :=
		eventArgs[0], eventArgs[1], eventArgs[2], eventArgs[3]

	tripDuration, err := computeTripDuration(startTimeStr, stopTimeStr)
	if err != nil {
		return fmt.Errorf("failed to compute trip duration for Trip event %v: %w", eventArgs, err)
	}

	tripMileage, err := computeTripMileage(tripMileageStr)
	if err != nil {
		return fmt.Errorf("failed to compute trip mileage for Trip event %v: %w", eventArgs, err)
	}

	if isUsableTripSample(tripMileage, tripDuration) {
		if err := eventStore.RecordTrip(&eventstore.TripInfo{
			DriverFirstName: driverFirstName,
			TripDuration:    tripDuration,
			TripMileage:     tripMileage,
		}); err != nil {
			return fmt.Errorf("failed to record Trip event %v with EventStore: %w", eventArgs, err)
		}
	}

	return nil
}

func computeTripDuration(startTimeStr, stopTimeStr string) (time.Duration, error) {
	const expectedTimeFormatStr string = "15:04"

	startTime, err := time.Parse(expectedTimeFormatStr, startTimeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse start time %s (likely doesn't conform to %s format): %w",
			startTimeStr, expectedTimeFormatStr, err)
	}

	stopTime, err := time.Parse(expectedTimeFormatStr, stopTimeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse stop time %s (likely doesn't conform to %s format): %w",
			stopTimeStr, expectedTimeFormatStr, err)
	}

	if !startTime.Before(stopTime) {
		return 0, fmt.Errorf("start time %s doesn't come before stop time %s",
			startTimeStr, stopTimeStr)
	}

	return stopTime.Sub(startTime), nil
}

func computeTripMileage(tripMileageStr string) (float32, error) {
	tripMileage64, err := strconv.ParseFloat(tripMileageStr, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse miles driven %s: %w", tripMileageStr, err)
	}

	return float32(tripMileage64), nil
}

func isUsableTripSample(tripMileage float32, tripDuration time.Duration) bool {
	tripSpeedMph := mathutils.ComputeSpeedMph32(tripMileage, tripDuration)

	if tripSpeedMph < 5.0 || tripSpeedMph > 100.0 {
		return false
	}

	return true
}
