package input

import (
	"bufio"
	"fmt"
	"io"
)

// StartReading scans `eventSource` for `Event`s in the background, and streams them out
// (encapsulated in `EventEnvelope`s) over the returned channel.
func StartReading(eventSource io.ReadCloser) <-chan *EventEnvelope {
	eventC := make(chan *EventEnvelope)

	go func() {
		// Close in the reverse order of opening (in keeping with the generally-wise "destructors must
		// run in the reverse order of constructors" principle).
		defer eventSource.Close()
		defer close(eventC)

		scanner := bufio.NewScanner(eventSource)
		for scanner.Scan() {
			body := Event(scanner.Text())
			eventC <- NewEventEnvelopeForBody(&body)
		}

		if err := scanner.Err(); err != nil {
			eventC <- NewEventEnvelopeForError(
				fmt.Errorf("error reading input: %w", err))
			return
		}
	}()

	return eventC
}
