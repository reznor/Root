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

	log.Printf("\nXXX\n")

	reportGenerator := output.NewReportGenerator()
	eventStore.Visit(reportGenerator)
	reportGenerator.Generate()
}
