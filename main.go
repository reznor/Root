package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func openInputFile() (*os.File, error) {
	inputFile := os.Stdin
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

type InputEvent string

type InputEventEnvelope struct {
	err error

	// EXPLAIN: Why this is a pointer (for future compatibility, as the type can change and become meatier)
	body *InputEvent
}

func NewInputEventEnvelopeWithError(err error) *InputEventEnvelope {
	return &InputEventEnvelope{
		err: err,
	}
}

func NewInputEventEnvelopeWithBody(inputEvent *InputEvent) *InputEventEnvelope {
	return &InputEventEnvelope{
		body: inputEvent,
	}
}

func startReadingInput() <-chan *InputEventEnvelope {
	inputC := make(chan *InputEventEnvelope)

	go func() {
		defer close(inputC)

		inputFile, err := openInputFile()
		if err != nil {
			inputC <- NewInputEventEnvelopeWithError(
				fmt.Errorf("error opening input file: %w", err))
			return
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			body := InputEvent(scanner.Text())
			inputC <- NewInputEventEnvelopeWithBody(&body)
		}

		if err := scanner.Err(); err != nil {
			inputC <- NewInputEventEnvelopeWithError(
				fmt.Errorf("error reading input file %s: %w", inputFile.Name(), err))
			return
		}

		log.Printf("Finished reading input file %s", inputFile.Name())
	}()

	return inputC
}

func main() {
	inputC := startReadingInput()

	for inputEventEnvelope := range inputC {
		if inputEventEnvelope.err != nil {
			log.Fatalf("Error event --> %s", inputEventEnvelope.err)
			// EXPLAIN: Choose to stop all processing instead of using `continue` to follow the principle of fail-early-fail-hard; better to not mess up any state once the input stream has gotten into an error state, but this constraint can change as the requirements/functionality changes.
			break
		}

		log.Printf("Read event ==> %s", *inputEventEnvelope.body)
	}
}
