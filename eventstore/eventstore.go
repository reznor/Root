package eventstore

import (
	"time"
)

type driverSummary struct {
	totalDurationDriven time.Duration
	// EXPLAIN: This is a float64 because the corresponding field in TripInfo is a float32, and this avoids overflow.
	totalMilesDriven float64
}

type EventStore struct {
	driverSummaries map[string]*driverSummary
}

func New() *EventStore {
	return &EventStore{
		driverSummaries: make(map[string]*driverSummary),
	}
}

// EXPLAIN: How these 2 Info structs are our interface to all the `eventhandlers`, and need to be maintained in a backwards-compatible manner -- it's not very important to worry about the layout of these structs, and it's fine if, with age, they end up having many unstructured fields; it's the responsibility of the methods in this package to provide a separation of domains between what clients deal with, and what is actually stored internally.
type DriverInfo struct {
	FirstName string
}

type TripInfo struct {
	DriverFirstName string
	TripDuration    time.Duration
	TripMileage     float32
}

func (es *EventStore) RegisterDriver(driverInfo *DriverInfo) error {
	if _, exists := es.driverSummaries[driverInfo.FirstName]; !exists {
		es.registerDriver(driverInfo)
	}

	return nil
}

func (es *EventStore) registerDriver(driverInfo *DriverInfo) *driverSummary {
	newDriverSummary := &driverSummary{}

	es.driverSummaries[driverInfo.FirstName] = newDriverSummary

	return newDriverSummary
}

func (es *EventStore) RecordTrip(tripInfo *TripInfo) error {
	// Register the Driver lazily, if needed.
	driverSummary := es.driverSummaries[tripInfo.DriverFirstName]
	if driverSummary == nil {
		driverSummary = es.registerDriver(&DriverInfo{
			FirstName: tripInfo.DriverFirstName,
		})
	}

	driverSummary.totalMilesDriven += float64(tripInfo.TripMileage)
	driverSummary.totalDurationDriven += tripInfo.TripDuration

	return nil
}
