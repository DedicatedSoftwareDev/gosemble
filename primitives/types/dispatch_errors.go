package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	DispatchErrorOther sc.U8 = iota
	DispatchErrorCannotLookup
	DispatchErrorBadOrigin
	DispatchErrorModule
	DispatchErrorConsumerRemaining
	DispatchErrorNoProviders
	DispatchErrorTooManyConsumers
	DispatchErrorToken
	DispatchErrorArithmetic
	DispatchErrorTransactional
	DispatchErrorExhausted
	DispatchErrorCorruption
	DispatchErrorUnavailable
)

type DispatchError = sc.VaryingData

func NewDispatchErrorOther(str sc.Str) DispatchError {
	return sc.NewVaryingData(DispatchErrorOther, str)
}

func NewDispatchErrorCannotLookup() DispatchError {
	return sc.NewVaryingData(DispatchErrorCannotLookup)
}

func NewDispatchErrorBadOrigin() DispatchError {
	return sc.NewVaryingData(DispatchErrorBadOrigin)
}

func NewDispatchErrorModule(customModuleError CustomModuleError) DispatchError {
	return sc.NewVaryingData(DispatchErrorModule, customModuleError)
}

func NewDispatchErrorConsumerRemaining() DispatchError {
	return sc.NewVaryingData(DispatchErrorConsumerRemaining)
}

func NewDispatchErrorNoProviders() DispatchError {
	return sc.NewVaryingData(DispatchErrorNoProviders)
}

func NewDispatchErrorTooManyConsumers() DispatchError {
	return sc.NewVaryingData(DispatchErrorTooManyConsumers)
}

func NewDispatchErrorToken(tokenError TokenError) DispatchError {
	// TODO: type safety
	return sc.NewVaryingData(DispatchErrorToken, tokenError)
}

func NewDispatchErrorArithmetic(arithmeticError ArithmeticError) DispatchError {
	// TODO: type safety
	return sc.NewVaryingData(DispatchErrorArithmetic, arithmeticError)
}

func NewDispatchErrorTransactional(transactionalError TransactionalError) DispatchError {
	// TODO: type safety
	return sc.NewVaryingData(DispatchErrorTransactional, transactionalError)
}

func NewDispatchErrorExhausted() DispatchError {
	return sc.NewVaryingData(DispatchErrorExhausted)
}

func NewDispatchErrorCorruption() DispatchError {
	return sc.NewVaryingData(DispatchErrorCorruption)
}

func NewDispatchErrorUnavailable() DispatchError {
	return sc.NewVaryingData(DispatchErrorUnavailable)
}

func DecodeDispatchError(buffer *bytes.Buffer) DispatchError {
	b := sc.DecodeU8(buffer)

	switch b {
	case DispatchErrorOther:
		value := sc.DecodeStr(buffer)
		return NewDispatchErrorOther(value)
	case DispatchErrorCannotLookup:
		return NewDispatchErrorCannotLookup()
	case DispatchErrorBadOrigin:
		return NewDispatchErrorBadOrigin()
	case DispatchErrorModule:
		module := DecodeCustomModuleError(buffer)
		return NewDispatchErrorModule(module)
	case DispatchErrorConsumerRemaining:
		return NewDispatchErrorConsumerRemaining()
	case DispatchErrorNoProviders:
		return NewDispatchErrorNoProviders()
	case DispatchErrorTooManyConsumers:
		return NewDispatchErrorTooManyConsumers()
	case DispatchErrorToken:
		tokenError := DecodeTokenError(buffer)
		return NewDispatchErrorToken(tokenError)
	case DispatchErrorArithmetic:
		arithmeticError := DecodeArithmeticError(buffer)
		return NewDispatchErrorArithmetic(arithmeticError)
	case DispatchErrorTransactional:
		transactionalError := DecodeTransactionalError(buffer)
		return NewDispatchErrorTransactional(transactionalError)
	case DispatchErrorExhausted:
		return NewDispatchErrorExhausted()
	case DispatchErrorCorruption:
		return NewDispatchErrorCorruption()
	case DispatchErrorUnavailable:
		return NewDispatchErrorUnavailable()
	default:
		log.Critical("invalid DispatchError type")
	}

	panic("unreachable")
}

// CustomModuleError A custom error in a module.
type CustomModuleError struct {
	Index   sc.U8             // Module index matching the metadata module index.
	Error   sc.U32            // Module specific error value.
	Message sc.Option[sc.Str] // Varying data type Option (Definition 190). The optional value is a SCALE encoded byte array containing a valid UTF-8 sequence.
}

func (e CustomModuleError) Encode(buffer *bytes.Buffer) {
	e.Index.Encode(buffer)
	e.Error.Encode(buffer)
	//e.Message.Encode(buffer) // Skipped in codec
}

func DecodeCustomModuleError(buffer *bytes.Buffer) CustomModuleError {
	e := CustomModuleError{}
	e.Index = sc.DecodeU8(buffer)
	e.Error = sc.DecodeU32(buffer)
	//e.Message = sc.DecodeOption[sc.Str](buffer) // Skipped in codec
	return e
}

func (e CustomModuleError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

// DispatchErrorWithPostInfo Result of a `Dispatchable` which contains the `DispatchResult` and additional information about
// the `Dispatchable` that is only known post dispatch.
type DispatchErrorWithPostInfo[T sc.Encodable] struct {
	// Additional information about the `Dispatchable` which is only known post dispatch.
	PostInfo T

	// The actual `DispatchResult` indicating whether the dispatch was successful.
	Error DispatchError
}

func (e DispatchErrorWithPostInfo[PostDispatchInfo]) Encode(buffer *bytes.Buffer) {
	e.PostInfo.Encode(buffer)
	e.Error.Encode(buffer)
}

func DecodeErrorWithPostInfo(buffer *bytes.Buffer) DispatchErrorWithPostInfo[PostDispatchInfo] {
	e := DispatchErrorWithPostInfo[PostDispatchInfo]{}
	e.PostInfo = DecodePostDispatchInfo(buffer)
	e.Error = DecodeDispatchError(buffer)
	return e
}

func (e DispatchErrorWithPostInfo[PostDispatchInfo]) Bytes() []byte {
	return sc.EncodedBytes(e)
}
