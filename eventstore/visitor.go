package eventstore

import (
	"log"
	"time"
)

// VisitableEntity provides the read-side contract that `EventStore` presents to the other components in
// the system.
//
// ============================================== Maintainer Notes ==============================================
//
// Note that akin to the `*Info` structs that are a crucial part of the contract of `EventStore`'s write-side
// API, this also provides complete isolation of clients from the internal storage format of `EventStore` -- it
// is the responsibility of `Visit` to perform any translation required to map between this and the internal
// storage format.
type VisitableEntity struct {
	DriverFirstName     string
	TotalDurationDriven time.Duration
	TotalMilesDriven    float64
}

// VisitorInterface specifies the expectations of a client that wishes to make use of `Visit`.
type VisitorInterface interface {
	Visit(*VisitableEntity)
}

// Visit provides a highly-curated and controlled mechanism for clients to get access to the information
// stored within `EventStore`.
//
// See https://en.wikipedia.org/wiki/Visitor_pattern for the benefits of the Visitor design pattern.
func (es *EventStore) Visit(visitor VisitorInterface) {
	for driverFirstName, driverSummary := range es.driverSummaries {
		visitor.Visit(&VisitableEntity{
			DriverFirstName:     driverFirstName,
			TotalDurationDriven: driverSummary.totalDurationDriven,
			TotalMilesDriven:    driverSummary.totalMilesDriven,
		})
	}
}

// Printer is a handy implementation of `VisitorInterface` to help with debugging during development.
type Printer struct{}

// NewPrinter creates a new `Printer`.
func NewPrinter() *Printer {
	return &Printer{}
}

// Conforms to `VisitorInterface`.
func (p Printer) Visit(visitableEntity *VisitableEntity) {
	log.Printf("DriverFirstName: %s TotalDurationDriven: %v TotalMilesDriven: %v\n",
		visitableEntity.DriverFirstName, visitableEntity.TotalDurationDriven, visitableEntity.TotalMilesDriven)
}

// Recorder is a handy implementation of `VisitorInterface` to provide simple programmatic
// introspection of the contents of `EventStore` (primarily to help with writing unit tests).
type Recorder struct {
	// Eschew []*VisitableEntity since this is primarily meant for testability, and dealing with pointers
	// makes error reporting unclear when the actual and the expected outputs don't match.
	Entities []VisitableEntity
}

// NewRecorder creates a new `Recorder`.
func NewRecorder() *Recorder {
	return &Recorder{
		Entities: make([]VisitableEntity, 0),
	}
}

// Conforms to `VisitorInterface`.
func (r *Recorder) Visit(visitableEntity *VisitableEntity) {
	r.Entities = append(r.Entities, *visitableEntity)
}
