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

// Withdraw withdraws `value` free balance from `who`, respecting existence requirements.
// Does not do anything if value is 0.
func Withdraw(who types.Address32, value sc.U128, reasons sc.U8, liveness types.ExistenceRequirement) (types.Balance, types.DispatchError) {
	if value.ToBigInt().Cmp(constants.Zero) == 0 {
		return sc.NewU128FromUint64(uint64(0)), nil
	}

	result := tryMutateAccount(who, func(account *types.AccountData, _ bool) sc.Result[sc.Encodable] {
		newFromAccountFree := new(big.Int).Sub(account.Free.ToBigInt(), value.ToBigInt())

		if newFromAccountFree.Cmp(constants.Zero) < 0 {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: types.NewDispatchErrorModule(types.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorInsufficientBalance),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		existentialDeposit := balances.ExistentialDeposit
		sumNewFreeReserved := new(big.Int).Add(newFromAccountFree, account.Reserved.ToBigInt())
		sumFreeReserved := new(big.Int).Add(account.Free.ToBigInt(), account.Reserved.ToBigInt())

		wouldBeDead := sumNewFreeReserved.Cmp(existentialDeposit) < 0
		wouldKill := wouldBeDead && (sumFreeReserved.Cmp(existentialDeposit) >= 0)

		if !(liveness == types.ExistenceRequirementAllowDeath || !wouldKill) {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: types.NewDispatchErrorModule(types.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorKeepAlive),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		err := ensureCanWithdraw(who, value.ToBigInt(), types.Reasons(reasons), newFromAccountFree)
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}

		account.Free = sc.NewU128FromBigInt(newFromAccountFree)

		system.DepositEvent(events.NewEventWithdraw(who.FixedSequence, value))

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    value,
		}
	})

	if result.HasError {
		return types.Balance{}, result.Value.(types.DispatchError)
	}

	return value, nil
}
