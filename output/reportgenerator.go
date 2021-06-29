package output

import (
	"container/heap"
	"fmt"

	"root.challenge/eventstore"
	"root.challenge/mathutils"
)

type ReportGenerator struct {
	heapElements []*eventstore.VisitableEntity
}

func NewReportGenerator() *ReportGenerator {
	rg := &ReportGenerator{
		heapElements: make([]*eventstore.VisitableEntity, 0),
	}

	heap.Init(rg)

	return rg
}

type GeneratedReport []string

func (rg *ReportGenerator) Generate() GeneratedReport {
	generatedReport := make(GeneratedReport, 0)

	for len(rg.heapElements) > 0 {
		ve := heap.Pop(rg).(*eventstore.VisitableEntity)

		var averageSpeedDisplayStr string
		if ve.TotalDurationDriven > 0 {
			averageSpeedDisplayStr = fmt.Sprintf(" @ %v mph", mathutils.RoundFloat64ToInt64(
				mathutils.ComputeSpeedMph64(ve.TotalMilesDriven, ve.TotalDurationDriven)))
		}

		generatedReport = append(generatedReport, fmt.Sprintf("%s: %v miles%s",
			ve.DriverFirstName, mathutils.RoundFloat64ToInt64(ve.TotalMilesDriven), averageSpeedDisplayStr))
	}

	return generatedReport
}

// Conforms to eventstore.Visitor.
func (rg *ReportGenerator) Visit(visitableEntity *eventstore.VisitableEntity) {
	heap.Push(rg, visitableEntity)
}

// Conforms to heap.Interface.
func (rg ReportGenerator) Len() int {
	return len(rg.heapElements)
}

// Conforms to heap.Interface.
func (rg ReportGenerator) Less(i, j int) bool {
	return rg.heapElements[i].TotalMilesDriven > rg.heapElements[j].TotalMilesDriven
}

// Conforms to heap.Interface.
func (rg ReportGenerator) Swap(i, j int) {
	rg.heapElements[i], rg.heapElements[j] = rg.heapElements[j], rg.heapElements[i]
}

// Conforms to heap.Interface.
func (rg *ReportGenerator) Push(x interface{}) {
	rg.heapElements = append(rg.heapElements, x.(*eventstore.VisitableEntity))
}

// Conforms to heap.Interface.
func (rg *ReportGenerator) Pop() interface{} {
	old := rg.heapElements
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // Avoid memory leak.
	rg.heapElements = old[0 : n-1]
	return x
}
