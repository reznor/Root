package main

import (
	"log"

	"root.challenge/eventprocessor"
	"root.challenge/eventstore"
	"root.challenge/input"
)

func main() {
	eventStore := eventstore.New()

	for err := range eventprocessor.New().Process(input.StartReading(), eventStore) {
		log.Printf("Error processing events: %s", err)
	}
}
