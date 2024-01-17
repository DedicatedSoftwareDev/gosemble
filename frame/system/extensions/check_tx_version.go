package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckTxVersion struct{}

func (_ CheckTxVersion) AdditionalSigned() (ok sc.U32, err primitives.TransactionValidityError) {
	return constants.RuntimeVersion.TransactionVersion, err
}

func (_ CheckTxVersion) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	ok = primitives.DefaultValidTransaction()
	return ok, err
}

func (v CheckTxVersion) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = v.Validate(who, call, info, length)
	return ok, err
}
