package system

import (
	"reflect"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var ZeroAddress = primitives.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

type CheckNonZeroAddress primitives.Address32

func (a CheckNonZeroAddress) AdditionalSigned() (ok sc.Empty, err primitives.TransactionValidityError) {
	ok = sc.Empty{}
	return ok, err
}

func (who CheckNonZeroAddress) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	// TODO:
	// Not sure when this is possible.
	// Checks signed transactions but will fail
	// before this check if the address is all zeros.
	if !reflect.DeepEqual(who, ZeroAddress) {
		ok = primitives.DefaultValidTransaction()
		return ok, err
	}

	err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadSigner())

	return ok, err
}

func (a CheckNonZeroAddress) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (ok primitives.Pre, err primitives.TransactionValidityError) {
	_, err = a.Validate(who, call, info, length)
	return ok, err
}
