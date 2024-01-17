package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: add more test cases

func Test_EncodeTransactionValidityError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       TransactionValidityError
		expectation []byte
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       NewTransactionValidityError(NewInvalidTransactionPayment()),
			expectation: []byte{0x00, 0x01},
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       NewTransactionValidityError(NewUnknownTransactionCannotLookup()),
			expectation: []byte{0x01, 0x00},
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

func Test_DecodeTransactionValidityError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation TransactionValidityError
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       []byte{0x00, 0x01},
			expectation: NewTransactionValidityError(NewInvalidTransactionPayment()),
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       []byte{0x01, 0x00},
			expectation: NewTransactionValidityError(NewUnknownTransactionCannotLookup()),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeTransactionValidityError(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
