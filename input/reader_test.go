package input_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"root.challenge/input"
)

func TestStartReading(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedOutput []*input.EventEnvelope
	}{
		"EmptyInput": {
			input:          "",
			expectedOutput: []*input.EventEnvelope{},
		},
		"EmptyLineAsInput": {
			input: "      ",
			expectedOutput: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("      ")),
			},
		},
		"SimpleInput": {
			input: "ABC DEF  GHI  JKL",
			expectedOutput: []*input.EventEnvelope{
				input.NewEventEnvelopeForBody(input.NewEventFromString("ABC DEF  GHI  JKL")),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualOutput := make([]*input.EventEnvelope, 0)
			for eventEnvelope := range input.StartReading(io.NopCloser(strings.NewReader(tc.input))) {
				actualOutput = append(actualOutput, eventEnvelope)
			}

			if !reflect.DeepEqual(actualOutput, tc.expectedOutput) {
				t.Fatalf("expected: %#v, got: %#v", tc.expectedOutput, actualOutput)
			}
		})
	}
}
