package eventstore

import (
	"log"
	"time"
)

type VisitableEntity struct {
	DriverFirstName     string
	TotalDurationDriven time.Duration
	TotalMilesDriven    float64
}

type VisitorInterface interface {
	Visit(*VisitableEntity)
}

func (es *EventStore) Visit(visitor VisitorInterface) {
	for driverFirstName, driverSummary := range es.driverSummaries {
		visitor.Visit(&VisitableEntity{
			DriverFirstName:     driverFirstName,
			TotalDurationDriven: driverSummary.totalDurationDriven,
			TotalMilesDriven:    driverSummary.totalMilesDriven,
		})
	}
}

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p Printer) Visit(visitableEntity *VisitableEntity) {
	log.Printf("DriverFirstName: %s TotalDurationDriven: %v TotalMilesDriven: %v\n",
		visitableEntity.DriverFirstName, visitableEntity.TotalDurationDriven, visitableEntity.TotalMilesDriven)
}

type Recorder struct {
	// Eschew []*VisitableEntity since this is primarily meant for testability, and dealing with pointers
	// makes error reporting unclear when the actual and the expected outputs don't match.
	Entities []VisitableEntity
}

func NewRecorder() *Recorder {
	return &Recorder{
		Entities: make([]VisitableEntity, 0),
	}
}

func (r *Recorder) Visit(visitableEntity *VisitableEntity) {
	r.Entities = append(r.Entities, *visitableEntity)
}
