package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckGenesis struct{}

func (_ CheckGenesis) AdditionalSigned() (ok primitives.H256, err primitives.TransactionValidityError) {
	ok = primitives.H256(system.StorageGetBlockHash(sc.U32(0)))
	return ok, err
}

func (_ CheckGenesis) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	ok = primitives.DefaultValidTransaction()
	return ok, err
}

func (g CheckGenesis) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = g.Validate(who, call, info, length)
	return ok, err
}
