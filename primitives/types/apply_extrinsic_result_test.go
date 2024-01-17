package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeApplyExtrinsicResult(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       ApplyExtrinsicResult
		expectation []byte
	}{
		{
			label:       "Encode ApplyExtrinsicResult(NewDispatchOutcome(None))",
			input:       NewApplyExtrinsicResult(NewDispatchOutcome(nil)),
			expectation: []byte{0x00, 0x00},
		},
		{
			label:       "Encode ApplyExtrinsicResult(NewDispatchOutcome(NewDispatchErrorBadOrigin))",
			input:       NewApplyExtrinsicResult(NewDispatchOutcome(NewDispatchErrorBadOrigin())),
			expectation: []byte{0x00, 0x01, 0x02},
		},
		{
			label:       "Encode ApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransactionCall))",
			input:       NewApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransactionCall())),
			expectation: []byte{0x01, 0x00, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeApplyExtrinsicResult(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation ApplyExtrinsicResult
	}{
		{
			label:       "Decode ApplyExtrinsicResult(NewDispatchOutcome(None))",
			expectation: NewApplyExtrinsicResult(NewDispatchOutcome(nil)),
			input:       []byte{0x00, 0x00},
		},
		{
			label:       "Decode ApplyExtrinsicResult(NewDispatchOutcome(NewDispatchErrorBadOrigin))",
			expectation: NewApplyExtrinsicResult(NewDispatchOutcome(NewDispatchErrorBadOrigin())),
			input:       []byte{0x00, 0x01, 0x02},
		},
		{
			label:       "Decode ApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransactionCall)",
			expectation: NewApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransactionCall())),
			input:       []byte{0x01, 0x00, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeApplyExtrinsicResult(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
