package main

import (
	"fmt"
	"log"
	"os"

	"root.challenge/eventprocessor"
	"root.challenge/eventstore"
	"root.challenge/input"
	"root.challenge/output"
)

func main() {
	inputFile, err := openInputFile()
	if err != nil {
		log.Fatalf("Error opening input file: %s", err)
		return
	}

	eventStore := eventstore.New()

	for err := range eventprocessor.New().Process(input.StartReading(inputFile), eventStore) {
		log.Printf("Error processing events: %s", err)
	}

	reportGenerator := output.NewReportGenerator()
	eventStore.Visit(reportGenerator)
	for _, reportEntry := range reportGenerator.Generate() {
		log.Println(reportEntry)
	}
}

func openInputFile() (*os.File, error) {
	// Default to reading from standard input.
	inputFile := os.Stdin

	// Give preference to files explicitly specified as input on the command line.
	if len(os.Args) == 2 {
		inputFileName := os.Args[1]

		f, err := os.Open(inputFileName)
		if err != nil {
			return nil, fmt.Errorf("error opening input file %s: %w", inputFileName, err)
		}

		inputFile = f
	}

	return inputFile, nil
}
