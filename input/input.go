package input

import (
	"bufio"
	"fmt"
	"os"
)

func StartReading() <-chan *EventEnvelope {
	inputC := make(chan *EventEnvelope)

	go func() {
		defer close(inputC)

		inputFile, err := openInputFile()
		if err != nil {
			inputC <- newEventEnvelopeForError(
				fmt.Errorf("error opening input file: %w", err))
			return
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			body := Event(scanner.Text())
			inputC <- newEventEnvelopeForBody(&body)
		}

		if err := scanner.Err(); err != nil {
			inputC <- newEventEnvelopeForError(
				fmt.Errorf("error reading input file %s: %w", inputFile.Name(), err))
			return
		}
	}()

	return inputC
}

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
