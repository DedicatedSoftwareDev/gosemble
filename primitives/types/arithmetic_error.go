package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	ArithmeticErrorUnderflow sc.U8 = iota
	ArithmeticErrorOverflow
	ArithmeticErrorDivisionByZero
)

type ArithmeticError = sc.VaryingData

func NewArithmeticErrorUnderflow() ArithmeticError {
	return sc.NewVaryingData(ArithmeticErrorUnderflow)
}

func NewArithmeticErrorOverflow() ArithmeticError {
	return sc.NewVaryingData(ArithmeticErrorOverflow)
}

func NewArithmeticErrorDivisionByZero() ArithmeticError {
	return sc.NewVaryingData(ArithmeticErrorDivisionByZero)
}

func DecodeArithmeticError(buffer *bytes.Buffer) ArithmeticError {
	b := sc.DecodeU8(buffer)

	switch b {
	case ArithmeticErrorUnderflow:
		return NewArithmeticErrorUnderflow()
	case ArithmeticErrorOverflow:
		return NewArithmeticErrorOverflow()
	case ArithmeticErrorDivisionByZero:
		return NewArithmeticErrorDivisionByZero()
	default:
		log.Critical("invalid ArithmeticError type")
	}

	panic("unreachable")
}
