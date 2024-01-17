package types

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"
)

type TransactionPriority = sc.U64
type TransactionLongevity = sc.U64

type TransactionTag = sc.Sequence[sc.U8]

// ValidTransaction Contains information concerning a valid transaction.
type ValidTransaction struct {
	// Priority of the transaction.
	//
	// Priority determines the ordering of two transactions that have all
	// their dependencies (required tags) satisfied.
	Priority TransactionPriority

	// Transaction dependencies
	//
	// A non-empty list signifies that some other transactions which provide
	// given tags are required to be included before that one.
	Requires sc.Sequence[TransactionTag]

	// Provided tags
	//
	// A list of tags this transaction provides. Successfully importing the transaction
	// will enable other transactions that depend on (require) those tags to be included as well.
	// Provided and required tags allow Substrate to build a dependency graph of transactions
	// and import them in the right (linear) order.
	Provides sc.Sequence[TransactionTag]

	// Transaction longevity
	//
	// Longevity describes minimum number of blocks the validity is correct.
	// After this period transaction should be removed from the pool or revalidated.
	Longevity TransactionLongevity

	// A flag indicating if the transaction should be propagated to other peers.
	//
	// By setting `false` here the transaction will still be considered for
	// including in blocks that are authored on the current node, but will
	// never be sent to other peers.
	Propagate sc.Bool
}

func (tx ValidTransaction) Encode(buffer *bytes.Buffer) {
	tx.Priority.Encode(buffer)
	tx.Requires.Encode(buffer)
	tx.Provides.Encode(buffer)
	tx.Longevity.Encode(buffer)
	tx.Propagate.Encode(buffer)
}

func DecodeValidTransaction(buffer *bytes.Buffer) ValidTransaction {
	return ValidTransaction{
		Priority:  sc.DecodeU64(buffer),
		Requires:  sc.DecodeSequence[TransactionTag](buffer),
		Provides:  sc.DecodeSequence[TransactionTag](buffer),
		Longevity: sc.DecodeU64(buffer),
		Propagate: sc.DecodeBool(buffer),
	}
}

func (tx ValidTransaction) Bytes() []byte {
	return sc.EncodedBytes(tx)
}

func DefaultValidTransaction() ValidTransaction {
	return ValidTransaction{
		Priority:  TransactionPriority(0),
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}
}

// Combine two instances into one, as a best effort. This will take the superset of each of the
// `provides` and `requires` tags, it will sum the priorities, take the minimum longevity and
// the logic *And* of the propagate flags.
func (vt ValidTransaction) CombineWith(otherVt ValidTransaction) ValidTransaction {
	priority := vt.Priority.SaturatingAdd(otherVt.Priority)
	requires := append(vt.Requires, otherVt.Requires...)
	provides := append(vt.Provides, otherVt.Provides...)
	longevity := vt.Longevity.Min(otherVt.Longevity)
	propagate := vt.Propagate && otherVt.Propagate

	return ValidTransaction{
		Priority:  priority,
		Requires:  requires,
		Provides:  provides,
		Longevity: longevity,
		Propagate: propagate,
	}
}
