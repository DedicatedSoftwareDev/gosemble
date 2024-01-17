package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// TransactionValidityResult Information on a transaction's validity and, if valid, on how it relates to other transactions.
type TransactionValidityResult sc.VaryingData

func NewTransactionValidityResult(value sc.Encodable) TransactionValidityResult {
	switch value.(type) {
	case ValidTransaction, TransactionValidityError:
		return TransactionValidityResult(sc.NewVaryingData(value))
	default:
		log.Critical("invalid TransactionValidityResult type")
	}

	panic("unreachable")
}

func (r TransactionValidityResult) Encode(buffer *bytes.Buffer) {
	switch r[0].(type) {
	case ValidTransaction:
		sc.U8(0).Encode(buffer)
	case TransactionValidityError:
		sc.U8(1).Encode(buffer)
	default:
		log.Critical("invalid TransactionValidityResult type")
	}

	r[0].Encode(buffer)
}

func DecodeTransactionValidityResult(buffer *bytes.Buffer) TransactionValidityResult {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return NewTransactionValidityResult(DecodeValidTransaction(buffer))
	case 1:
		return NewTransactionValidityResult(DecodeTransactionValidityError(buffer))
	default:
		log.Critical("invalid TransactionValidityResult type")
	}

	panic("unreachable")
}

func (r TransactionValidityResult) Bytes() []byte {
	return sc.EncodedBytes(r)
}

func (r TransactionValidityResult) IsValidTransaction() sc.Bool {
	switch r[0].(type) {
	case ValidTransaction:
		return true
	default:
		return false
	}
}

func (r TransactionValidityResult) AsValidTransaction() ValidTransaction {
	if r.IsValidTransaction() {
		return r[0].(ValidTransaction)
	} else {
		log.Critical("not a ValidTransaction type")
	}

	panic("unreachable")
}
