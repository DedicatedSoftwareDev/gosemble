package types

import (
	"bytes"

	"github.com/LimeChain/gosemble/primitives/log"

	sc "github.com/LimeChain/goscale"
)

const (
	// Too many transactional layers have been spawned.
	TransactionalErrorLimitReached sc.U8 = iota
	// A transactional layer was expected, but does not exist.
	TransactionalErrorNoLayer
)

type TransactionalError = sc.VaryingData

func NewTransactionalErrorLimitReached() TransactionalError {
	return sc.NewVaryingData(TransactionalErrorLimitReached)
}

func NewTransactionalErrorNoLayer() TransactionalError {
	return sc.NewVaryingData(TransactionalErrorNoLayer)
}

func DecodeTransactionalError(buffer *bytes.Buffer) TransactionalError {
	b := sc.DecodeU8(buffer)

	switch b {
	case TransactionalErrorLimitReached:
		return NewTransactionalErrorLimitReached()
	case TransactionalErrorNoLayer:
		return NewTransactionalErrorNoLayer()
	default:
		log.Critical("invalid TransactionalError type")
	}

	panic("unreachable")
}
