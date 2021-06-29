package main

import (
	"log"

	"root.challenge/eventprocessor"
	"root.challenge/eventstore"
	"root.challenge/input"
	"root.challenge/output"
)

func main() {
	eventStore := eventstore.New()

	for err := range eventprocessor.New().Process(input.StartReading(), eventStore) {
		log.Printf("Error processing events: %s", err)
	}

	eventStore.Visit(eventstore.Printer{})

	reportGenerator := output.NewReportGenerator()
	eventStore.Visit(reportGenerator)
	for _, reportEntry := range reportGenerator.Generate() {
		log.Println(reportEntry)
	}
}
