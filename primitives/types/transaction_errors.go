package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	TransactionValidityErrorInvalidTransaction sc.U8 = iota
	TransactionValidityErrorUnknownTransaction
)

// TransactionValidityError Errors that can occur while checking the validity of a transaction.
type TransactionValidityError sc.VaryingData

func NewTransactionValidityError(value sc.Encodable) TransactionValidityError {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction:
	default:
		log.Critical("invalid TransactionValidityError type")
	}

	return TransactionValidityError(sc.NewVaryingData(value))
}

func (e TransactionValidityError) Encode(buffer *bytes.Buffer) {
	value := e[0]

	switch reflect.TypeOf(value) {
	case reflect.TypeOf(*new(InvalidTransaction)):
		buffer.Write([]byte{0x00})
	case reflect.TypeOf(*new(UnknownTransaction)):
		buffer.Write([]byte{0x01})
	default:
		log.Critical("invalid TransactionValidityError type")
	}

	value.Encode(buffer)
}

func DecodeTransactionValidityError(buffer *bytes.Buffer) TransactionValidityError {
	b := sc.DecodeU8(buffer)

	switch b {
	case TransactionValidityErrorInvalidTransaction:
		value := DecodeInvalidTransaction(buffer)
		return NewTransactionValidityError(value)
	case TransactionValidityErrorUnknownTransaction:
		value := DecodeUnknownTransaction(buffer)
		return NewTransactionValidityError(value)
	default:
		log.Critical("invalid TransactionValidityError type")
	}

	panic("unreachable")
}

func (e TransactionValidityError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

const (
	// The call of the transaction is not expected. Reject
	InvalidTransactionCall sc.U8 = iota

	// General error to do with the inability to pay some fees (e.g. account balance too low). Reject
	InvalidTransactionPayment

	// General error to do with the transaction not yet being valid (e.g. nonce too high). Don't Reject
	InvalidTransactionFuture

	// General error to do with the transaction being outdated (e.g. nonce too low). Reject
	InvalidTransactionStale

	// General error to do with the transaction's proofs (e.g. signature). Reject
	//
	// # Possible causes
	//
	// When using a signed extension that provides additional data for signing, it is required
	// that the signing and the verifying side use the same additional data. Additional
	// data will only be used to generate the signature, but will not be part of the transaction
	// itself. As the verifying side does not know which additional data was used while signing
	// it will only be able to assume a bad signature and cannot express a more meaningful error.
	InvalidTransactionBadProof

	// The transaction birth block is ancient. Reject
	//
	// # Possible causes
	//
	// For `FRAME`-based runtimes this would be caused by `current block number
	// - Era::birth block number > BlockHashCount`. (e.g. in Polkadot `BlockHashCount` = 2400, so
	//   a
	// transaction with birth block number 1337 would be valid up until block number 1337 + 2400,
	// after which point the transaction would be considered to have an ancient birth block.)
	InvalidTransactionAncientBirthBlock

	// The transaction would exhaust the resources of the current block. Don't Reject
	//
	// The transaction might be valid, but there are not enough resources
	// left in the current block.
	InvalidTransactionExhaustsResources

	// Any other custom invalid validity that is not covered by this enum. Reject
	InvalidTransactionCustom // + sc.U8

	// An extrinsic with mandatory dispatch resulted in an error. Reject
	// This is indicative of either a malicious validator or a buggy `provide_inherent`.
	// In any case, it can result in dangerously overweight blocks and therefore if
	// found, invalidates the block.
	InvalidTransactionBadMandatory

	// An extrinsic with a mandatory dispatch tried to be validated.
	// This is invalid; only inherent extrinsics are allowed to have mandatory dispatches.
	InvalidTransactionMandatoryValidation

	// The sending address is disabled or known to be invalid.
	InvalidTransactionBadSigner
)

type InvalidTransaction struct {
	sc.VaryingData
}

func NewInvalidTransactionCall() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionCall)}
}

func NewInvalidTransactionPayment() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionPayment)}
}

func NewInvalidTransactionFuture() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionFuture)}
}

func NewInvalidTransactionStale() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionStale)}
}

func NewInvalidTransactionBadProof() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadProof)}
}

func NewInvalidTransactionAncientBirthBlock() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionAncientBirthBlock)}
}

func NewInvalidTransactionExhaustsResources() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionExhaustsResources)}
}

func NewInvalidTransactionCustom(customError sc.U8) InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionCustom, customError)}
}

func NewInvalidTransactionBadMandatory() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadMandatory)}
}

func NewInvalidTransactionMandatoryValidation() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionMandatoryValidation)}
}

func NewInvalidTransactionBadSigner() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadSigner)}
}

func DecodeInvalidTransaction(buffer *bytes.Buffer) InvalidTransaction {
	b := sc.DecodeU8(buffer)

	switch b {
	case InvalidTransactionCall:
		return NewInvalidTransactionCall()
	case InvalidTransactionPayment:
		return NewInvalidTransactionPayment()
	case InvalidTransactionFuture:
		return NewInvalidTransactionFuture()
	case InvalidTransactionStale:
		return NewInvalidTransactionStale()
	case InvalidTransactionBadProof:
		return NewInvalidTransactionBadProof()
	case InvalidTransactionAncientBirthBlock:
		return NewInvalidTransactionAncientBirthBlock()
	case InvalidTransactionExhaustsResources:
		return NewInvalidTransactionExhaustsResources()
	case InvalidTransactionCustom:
		v := sc.DecodeU8(buffer)
		return NewInvalidTransactionCustom(v)
	case InvalidTransactionBadMandatory:
		return NewInvalidTransactionBadMandatory()
	case InvalidTransactionMandatoryValidation:
		return NewInvalidTransactionMandatoryValidation()
	case InvalidTransactionBadSigner:
		return NewInvalidTransactionBadSigner()
	default:
		log.Critical("invalid InvalidTransaction type")
	}

	panic("unreachable")
}

func (e InvalidTransaction) Bytes() []byte {
	return sc.EncodedBytes(e)
}

const (
	// Could not lookup some information that is required to validate the transaction. Reject
	UnknownTransactionCannotLookup sc.U8 = iota

	// No validator found for the given unsigned transaction. Reject
	UnknownTransactionNoUnsignedValidator

	// Any other custom unknown validity that is not covered by this type. Reject
	UnknownTransactionCustomUnknownTransaction // + sc.U8
)

type UnknownTransaction struct {
	sc.VaryingData
}

func NewUnknownTransactionCannotLookup() UnknownTransaction {
	return UnknownTransaction{sc.NewVaryingData(UnknownTransactionCannotLookup)}
}

func NewUnknownTransactionNoUnsignedValidator() UnknownTransaction {
	return UnknownTransaction{sc.NewVaryingData(UnknownTransactionNoUnsignedValidator)}
}

func NewUnknownTransactionCustomUnknownTransaction(unknown sc.U8) UnknownTransaction {
	return UnknownTransaction{sc.NewVaryingData(UnknownTransactionCustomUnknownTransaction, unknown)}
}

func DecodeUnknownTransaction(buffer *bytes.Buffer) UnknownTransaction {
	b := sc.DecodeU8(buffer)

	switch b {
	case UnknownTransactionCannotLookup:
		return NewUnknownTransactionCannotLookup()
	case UnknownTransactionNoUnsignedValidator:
		return NewUnknownTransactionNoUnsignedValidator()
	case UnknownTransactionCustomUnknownTransaction:
		v := sc.DecodeU8(buffer)
		return NewUnknownTransactionCustomUnknownTransaction(v)
	default:
		log.Critical("invalid UnknownTransaction type")
	}

	panic("unreachable")
}
