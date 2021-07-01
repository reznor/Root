package output

import (
	"container/heap"
	"fmt"

	"root.challenge/eventstore"
	"root.challenge/mathutils"
)

// ReportGenerator is used to generate a summary report of all the events that entered the system.
//
// ============================================== Maintainer Notes ==============================================
//
// The implementation visits `EventStore` and maintains a max-heap of `eventstore.VisitableEntity`
// objects (based on their `TotalMilesDriven` field) -- while a simple sort would suffice at small
// scale, the technique used here should remain fairly performant even at medium scale, and at large
// scale, the first thing to tweak will likely be capping the size of the max-heap (to control memory
// usage) and leveraging disk space to store the entire working set, working on sub-sections as needed
// (akin to an n-way merge sort); once the resource limits of a single machine are hit, the
// implementation will need to substantially change to run a distributed algorithm.
type ReportGenerator struct {
	heapElements []*eventstore.VisitableEntity
}

// NewReportGenerator creates a new `ReportGenerator`.
//
// This is where the (empty) max-heap is initialized.
func NewReportGenerator() *ReportGenerator {
	rg := &ReportGenerator{
		heapElements: make([]*eventstore.VisitableEntity, 0),
	}

	heap.Init(rg)

	return rg
}

// GeneratedReport defines the output type of `Generate`.
type GeneratedReport []string

// Generate returns a `GeneratedReport` containing the desired output-ready summary information.
//
// This is where we pop from the max-heap, retrieving the elements in descending order of
// `TotalMilesDriven`.
func (rg *ReportGenerator) Generate() GeneratedReport {
	generatedReport := make(GeneratedReport, 0)

	for len(rg.heapElements) > 0 {
		ve := heap.Pop(rg).(*eventstore.VisitableEntity)

		// Only display a speed component in the `GeneratedReport` format if the driver actually drove.
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

// Conforms to `eventstore.VisitorInterface`.
//
// This is where we push to the max-heap, building up the sorted data structure.
func (rg *ReportGenerator) Visit(visitableEntity *eventstore.VisitableEntity) {
	heap.Push(rg, visitableEntity)
}

// Conforms to `heap.Interface`.
func (rg ReportGenerator) Len() int {
	return len(rg.heapElements)
}

// Conforms to `heap.Interface`.
func (rg ReportGenerator) Less(i, j int) bool {
	// Maintain a max-heap (which is why the `Less` comparison uses '>') based on `TotalMilesDriven`.
	return rg.heapElements[i].TotalMilesDriven > rg.heapElements[j].TotalMilesDriven
}

// Conforms to `heap.Interface`.
func (rg ReportGenerator) Swap(i, j int) {
	rg.heapElements[i], rg.heapElements[j] = rg.heapElements[j], rg.heapElements[i]
}

// Conforms to `heap.Interface`.
func (rg *ReportGenerator) Push(x interface{}) {
	rg.heapElements = append(rg.heapElements, x.(*eventstore.VisitableEntity))
}

// Conforms to `heap.Interface`.
func (rg *ReportGenerator) Pop() interface{} {
	old := rg.heapElements
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // Avoid memory leak.
	rg.heapElements = old[0 : n-1]
	return x
}
