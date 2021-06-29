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

type Visitor interface {
	Visit(*VisitableEntity)
}

func (es *EventStore) Visit(visitor Visitor) {
	for driverFirstName, driverSummary := range es.driverSummaries {
		visitor.Visit(&VisitableEntity{
			DriverFirstName:     driverFirstName,
			TotalDurationDriven: driverSummary.totalDurationDriven,
			TotalMilesDriven:    driverSummary.totalMilesDriven,
		})
	}
}

type Printer struct{}

func (p Printer) Visit(visitableEntity *VisitableEntity) {
	log.Printf("DriverFirstName: %s TotalDurationDriven: %v TotalMilesDriven: %v\n",
		visitableEntity.DriverFirstName, visitableEntity.TotalDurationDriven, visitableEntity.TotalMilesDriven)
}
