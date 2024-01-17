package dispatchables

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

// DepositIntoExisting deposits `value` into the free balance of an existing target account `who`.
// If `value` is 0, it does nothing.
func DepositIntoExisting(who types.Address32, value sc.U128) (types.Balance, types.DispatchError) {
	if value.ToBigInt().Cmp(constants.Zero) == 0 {
		return sc.NewU128FromUint64(uint64(0)), nil
	}

	result := tryMutateAccount(who, func(from *types.AccountData, isNew bool) sc.Result[sc.Encodable] {
		if isNew {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: types.NewDispatchErrorModule(types.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorDeadAccount),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		sum := new(big.Int).Add(from.Free.ToBigInt(), value.ToBigInt())

		from.Free = sc.NewU128FromBigInt(sum)

		system.DepositEvent(events.NewEventDeposit(who.FixedSequence, value))

		return sc.Result[sc.Encodable]{}
	})

	if result.HasError {
		return types.Balance{}, result.Value.(types.DispatchError)
	}

	return value, nil
}
