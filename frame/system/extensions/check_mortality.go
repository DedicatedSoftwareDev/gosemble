package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func (e CheckMortality) AdditionalSigned() (ok primitives.H256, err primitives.TransactionValidityError) {
	current := sc.U64(system.StorageGetBlockNumber()) // TODO: impl saturated_into::<u64>()
	n := sc.U32(primitives.Era(e).Birth(current))     // TODO: impl saturated_into::<T::BlockNumber>()

	if !system.StorageExistsBlockHash(n) {
		err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionAncientBirthBlock())
		return ok, err
	} else {
		ok = primitives.H256(system.StorageGetBlockHash(n))
	}

	return ok, err
}

// TODO: to be able to provide a custom implementation of the Validate function
type CheckMortality primitives.Era

func (e CheckMortality) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	currentU64 := sc.U64(system.StorageGetBlockNumber()) // TODO: per module implementation

	validTill := primitives.Era(e).Death(currentU64)

	ok = primitives.DefaultValidTransaction()
	ok.Longevity = validTill.SaturatingSub(currentU64)

	return ok, err
}

func (e CheckMortality) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = e.Validate(who, call, info, length)
	return ok, err
}
