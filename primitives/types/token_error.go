package types

import (
	"bytes"

	"github.com/LimeChain/gosemble/primitives/log"

	sc "github.com/LimeChain/goscale"
)

const (
	TokenErrorNoFunds sc.U8 = iota
	TokenErrorWouldDie
	TokenErrorBelowMinimum
	TokenErrorCannotCreate
	TokenErrorUnknownAsset
	TokenErrorFrozen
	TokenErrorUnsupported
)

type TokenError = sc.VaryingData

func NewTokenErrorNoFounds() TokenError {
	return sc.NewVaryingData(TokenErrorNoFunds)
}

func NewTokenErrorWouldDie() TokenError {
	return sc.NewVaryingData(TokenErrorWouldDie)
}

func NewTokenErrorBelowMinimum() TokenError {
	return sc.NewVaryingData(TokenErrorBelowMinimum)
}

func NewTokenErrorCannotCreate() TokenError {
	return sc.NewVaryingData(TokenErrorCannotCreate)
}

func NewTokenErrorUnknownAsset() TokenError {
	return sc.NewVaryingData(TokenErrorUnknownAsset)
}

func NewTokenErrorFrozen() TokenError {
	return sc.NewVaryingData(TokenErrorFrozen)
}

func NewTokenErrorUnsupported() TokenError {
	return sc.NewVaryingData(TokenErrorUnsupported)
}

func DecodeTokenError(buffer *bytes.Buffer) TokenError {
	b := sc.DecodeU8(buffer)

	switch b {
	case TokenErrorNoFunds:
		return NewTokenErrorNoFounds()
	case TokenErrorWouldDie:
		return NewTokenErrorWouldDie()
	case TokenErrorBelowMinimum:
		return NewTokenErrorBelowMinimum()
	case TokenErrorCannotCreate:
		return NewTokenErrorCannotCreate()
	case TokenErrorUnknownAsset:
		return NewTokenErrorUnknownAsset()
	case TokenErrorFrozen:
		return NewTokenErrorFrozen()
	case TokenErrorUnsupported:
		return NewTokenErrorUnsupported()
	default:
		log.Critical("invalid TokenError type")
	}

	panic("unreachable")
}
