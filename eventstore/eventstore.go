package eventstore

import (
	"time"
)

// driverSummary represents the information about a driver that is pertinent to retain in `EventStore`.
//
// ============================================== Maintainer Notes ==============================================
//
// Note that this is a purely internal data structure, and is never exposed through any of the external
// contracts of this package -- this provides the flexibility to change it at will without any impact to
// clients of the package. Think of it as a database's storage format on disk.
type driverSummary struct {
	totalDurationDriven time.Duration
	// totalMilesDriven is a `float64` to avoid overflow (`TripInfo.TripMileage` is a `float32`, and all
	// those `float32`s are aggregated into this field).
	totalMilesDriven float64
}

// EventStore is a repository of relevant information about all the events that have flowed into the system.
//
// The only way for clients to access any of the information stored within is via `VisitorInterface`.
type EventStore struct {
	driverSummaries map[string]*driverSummary
}

// New creates a new `EventStore`.
func New() *EventStore {
	return &EventStore{
		driverSummaries: make(map[string]*driverSummary),
	}
}

// ============================================== Maintainer Notes ==============================================
//
// The *Info structs below are very much a part of the contract that `EventStore` presents to all the other
// components in the system (most notably, the concrete `eventhandler.Interface` implementations), and thus
// need to be maintained in a backwards-compatible manner -- it's not very important to worry about the layout
// of these structs, and it's fine if, with time, they end up containing many unstructured fields; it's the
// responsibility of the public methods on `EventStore` to provide the separation of domains between what clients
// deal with, and what is actually stored internally.

// DriverInfo encapsulates all the information about a driver that can be provided by clients of `EventStore`.
type DriverInfo struct {
	FirstName string
}

// TripInfo encapsulates all the information about a trip that can be provided by clients of `EventStore`.
type TripInfo struct {
	DriverFirstName string
	TripDuration    time.Duration
	TripMileage     float32
}

// RegisterDriver stores information about a new driver in the system.
//
// It's perfectly fine if multiple calls to this method are made with the same `DriverInfo` -- the
// implementation is guaranteed to be idempotent.
func (es *EventStore) RegisterDriver(driverInfo *DriverInfo) error {
	if _, exists := es.driverSummaries[driverInfo.FirstName]; !exists {
		es.registerDriver(driverInfo)
	}

	return nil
}

// registerDriver contains the core of what it takes to register a new driver -- all the public
// methods that invoke it should wholly delegate all the relevant functionality to this method.
func (es *EventStore) registerDriver(driverInfo *DriverInfo) *driverSummary {
	newDriverSummary := &driverSummary{}

	es.driverSummaries[driverInfo.FirstName] = newDriverSummary

	return newDriverSummary
}

// RecordTrip stores information about a new trip in the system.
//
// While it's currently perfectly fine for `TripInfo.DriverFirstName` to reference a driver that
// hasn't previously been registered via a call to `RegisterDriver`, that's not behavior that clients
// should come to depend upon -- this method currently takes care of performing the registration
// lazily because all the information that's needed for that operation is present in `TripInfo`, but
// that will almost certainly change when `DriverInfo` expands to become richer, at which time,
// this lazy-registration behavior of this method will cease to work, and will result in an error
// instead.
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
