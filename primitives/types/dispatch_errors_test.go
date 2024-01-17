package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeDispatchError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       DispatchError
		expectation []byte
	}{
		{label: "Encode(DispatchError('unknown error'))", input: NewDispatchErrorOther("unknown error"), expectation: []byte{0x00, 0x34, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72}},
		{label: "Encode(DispatchErrorCannotLookup)", input: NewDispatchErrorCannotLookup(), expectation: []byte{0x01}},
		{label: "Encode(DispatchErrorBadOrigin)", input: NewDispatchErrorBadOrigin(), expectation: []byte{0x02}},
		{label: "Encode(DispatchErrorCustomModule)", input: NewDispatchErrorModule(CustomModuleError{}), expectation: []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeDispatchError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation DispatchError
	}{

		{label: "DecodeDispatchError(0x00, 0x34, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72)", input: []byte{0x00, 0x34, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72}, expectation: NewDispatchErrorOther("unknown error")},
		{label: "DecodeDispatchError(0x01)", input: []byte{0x01}, expectation: NewDispatchErrorCannotLookup()},
		{label: "DecodeDispatchError(0x02)", input: []byte{0x02}, expectation: NewDispatchErrorBadOrigin()},
		{label: "DecodeDispatchError(0x03, 0x00, 0x00, 0x00)", input: []byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00}, expectation: NewDispatchErrorModule(CustomModuleError{})},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeDispatchError(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
