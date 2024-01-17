package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewMortalEra(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []sc.U64
		expectation Era
	}{
		{
			label: "NewMortalEra(64, 42)",
			input: []sc.U64{64, 42},
			expectation: Era{
				EraPeriod: 64,
				EraPhase:  42,
			},
		},
		{
			label: "NewMortalEra(32768, 20000)",
			input: []sc.U64{32768, 20000},
			expectation: Era{
				EraPeriod: 32768,
				EraPhase:  20000,
			},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			eraResult := NewMortalEra(testExample.input[0], testExample.input[1])
			assert.Equal(t, testExample.expectation, eraResult)
		})
	}
}

func Test_EncodeEra(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Era
		expectation []byte
	}{
		{
			label:       "Encode Era(ImmortalEra)",
			input:       Era{IsImmortal: true},
			expectation: []byte{0x00},
		},
		{
			label: "Encode Era(MortalEra(64, 42))",
			input: Era{
				IsImmortal: false,
				EraPeriod:  64,
				EraPhase:   42,
			},
			expectation: []byte{165, 2},
		},
		{
			label: "Encode Era(MortalEra(32768, 20000))",
			input: Era{
				IsImmortal: false,
				EraPeriod:  32768,
				EraPhase:   20000,
			},
			expectation: []byte{78, 156},
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

func Test_DecodeEra(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation Era
	}{
		{
			label:       "Decode Era(0x00)",
			input:       []byte{0x00},
			expectation: Era{IsImmortal: true},
		},
		{
			label: "Encode Era(165, 2)",
			input: []byte{165, 2},
			expectation: Era{
				IsImmortal: false,
				EraPeriod:  64,
				EraPhase:   42,
			},
		},
		{
			label: "Decode Long Era(78, 156)",
			input: []byte{78, 156},
			expectation: Era{
				IsImmortal: false,
				EraPeriod:  32768,
				EraPhase:   20000,
			},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeEra(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
