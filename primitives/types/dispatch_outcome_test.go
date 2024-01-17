package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeDispatchOutcome(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       DispatchOutcome
		expectation []byte
	}{
		{label: "Encode DispatchOutcome(None)", input: NewDispatchOutcome(nil), expectation: []byte{0x00}},
		{label: "Encode  DispatchOutcome(DispatchErrorBadOrigin)", input: NewDispatchOutcome(NewDispatchErrorBadOrigin()), expectation: []byte{0x01, 0x02}},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeDispatchOutcome(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation DispatchOutcome
	}{
		{label: "0x00", input: []byte{0x00}, expectation: NewDispatchOutcome(nil)},
		{label: "0x01, 0x02", input: []byte{0x01, 0x02}, expectation: NewDispatchOutcome(NewDispatchErrorBadOrigin())},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeDispatchOutcome(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
