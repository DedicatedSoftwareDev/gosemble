package system

import (
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonce sc.U32

func (n CheckNonce) AdditionalSigned() (ok sc.Empty, err primitives.TransactionValidityError) {
	ok = sc.Empty{}
	return ok, err
}

func (n CheckNonce) Validate(who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	// TODO: check if we can use just who
	account := system.StorageGetAccount((*who).FixedSequence)

	if sc.U32(n) < account.Nonce {
		err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
		return ok, err
	}

	encoded := (*who).Bytes()
	encoded = append(encoded, sc.ToCompact(sc.U32(n)).Bytes()...)
	provides := sc.Sequence[primitives.TransactionTag]{sc.BytesToSequenceU8(encoded)}

	var requires sc.Sequence[primitives.TransactionTag]
	if account.Nonce < sc.U32(n) {
		encoded := (*who).Bytes()
		encoded = append(encoded, sc.ToCompact(sc.U32(n)-1).Bytes()...)
		requires = sc.Sequence[primitives.TransactionTag]{sc.BytesToSequenceU8(encoded)}
	} else {
		requires = sc.Sequence[primitives.TransactionTag]{}
	}

	ok = primitives.ValidTransaction{
		Priority:  0,
		Requires:  requires,
		Provides:  provides,
		Longevity: primitives.TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}

	return ok, err
}

func (n CheckNonce) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	account := system.StorageGetAccount(who.FixedSequence)

	if sc.U32(n) != account.Nonce {
		if sc.U32(n) < account.Nonce {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
		} else {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionFuture())
		}
		return ok, err
	}

	account.Nonce += 1
	system.StorageSetAccount(who.FixedSequence, account)

	ok = primitives.Pre{}
	return ok, err
}
