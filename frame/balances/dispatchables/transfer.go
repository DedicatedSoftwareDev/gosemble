package dispatchables

import (
	"bytes"
	"math/big"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TransferCall struct {
	primitives.Callable
}

func NewTransferCall(args sc.VaryingData) TransferCall {
	call := TransferCall{
		Callable: primitives.Callable{
			ModuleId:   balances.ModuleIndex,
			FunctionId: balances.FunctionTransferIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c TransferCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c TransferCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c TransferCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c TransferCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c TransferCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c TransferCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ TransferCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 37_815 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(38_109_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ TransferCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ TransferCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ TransferCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ TransferCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := transfer(origin, args[0].(types.MultiAddress), value)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
		Ok:       types.PostDispatchInfo{},
	}
}

func (_ TransferCall) IsInherent() bool {
	return false
}

// transfer transfers liquid free balance from `source` to `dest`.
// Increases the free balance of `dest` and decreases the free balance of `origin` transactor.
// Must be signed by the transactor.
func transfer(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	to, e := types.DefaultAccountIdLookup().Lookup(dest)
	if e != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	transactor := origin.AsSigned()

	return trans(transactor, to, value, types.ExistenceRequirementAllowDeath)
}

// trans transfers `value` free balance from `from` to `to`.
// Does not do anything if value is 0 or `from` and `to` are the same.
func trans(from types.Address32, to types.Address32, value sc.U128, existenceRequirement types.ExistenceRequirement) types.DispatchError {
	bnInt := value.ToBigInt()
	if bnInt.Cmp(constants.Zero) == 0 || reflect.DeepEqual(from, to) {
		return nil
	}

	result := tryMutateAccountWithDust(to, func(toAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
		return tryMutateAccountWithDust(from, func(fromAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
			newFromAccountFree := new(big.Int).Sub(fromAccount.Free.ToBigInt(), value.ToBigInt())

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
			fromAccount.Free = sc.NewU128FromBigInt(newFromAccountFree)

			newToAccountFree := new(big.Int).Add(toAccount.Free.ToBigInt(), value.ToBigInt())
			toAccount.Free = sc.NewU128FromBigInt(newToAccountFree)

			existentialDeposit := balances.ExistentialDeposit
			if toAccount.Total().Cmp(existentialDeposit) < 0 {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   balances.ModuleIndex,
						Error:   sc.U32(errors.ErrorExistentialDeposit),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}

			err := ensureCanWithdraw(from, value.ToBigInt(), types.ReasonsAll, fromAccount.Free.ToBigInt())
			if err != nil {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value:    err,
				}
			}

			allowDeath := existenceRequirement == types.ExistenceRequirementAllowDeath
			allowDeath = allowDeath && system.CanDecProviders(from)

			if !(allowDeath || fromAccount.Total().Cmp(existentialDeposit) > 0) {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   balances.ModuleIndex,
						Error:   sc.U32(errors.ErrorKeepAlive),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}

			return sc.Result[sc.Encodable]{}
		})
	})

	if result.HasError {
		return result.Value.(types.DispatchError)
	}

	system.DepositEvent(events.NewEventTransfer(from.FixedSequence, to.FixedSequence, value))
	return nil
}

// ensureCanWithdraw checks that an account can withdraw from their balance given any existing withdraw restrictions.
func ensureCanWithdraw(who types.Address32, amount *big.Int, reasons types.Reasons, newBalance *big.Int) types.DispatchError {
	if amount.Cmp(constants.Zero) == 0 {
		return nil
	}

	accountInfo := system.StorageGetAccount(who.FixedSequence)
	minBalance := accountInfo.Frozen(reasons)
	if minBalance.Cmp(newBalance) > 0 {
		return types.NewDispatchErrorModule(types.CustomModuleError{
			Index:   balances.ModuleIndex,
			Error:   sc.U32(errors.ErrorLiquidityRestrictions),
			Message: sc.NewOption[sc.Str](nil),
		})
	}

	return nil
}

// mutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func mutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	return tryMutateAccount(who, f)
}

// tryMutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func tryMutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := tryMutateAccountWithDust(who, f)
	if result.HasError {
		return result
	}

	r := result.Value.(sc.VaryingData)

	// TODO: Convert this to an Option and uncomment it.
	// Check Substrate implementation for reference.
	//dustCleaner := r[1].(DustCleanerValue)
	//dustCleaner.Drop()

	return sc.Result[sc.Encodable]{HasError: false, Value: r[0].(sc.Encodable)}
}

func tryMutateAccountWithDust(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := system.TryMutateExists(who, func(maybeAccount *types.AccountData) sc.Result[sc.Encodable] {
		account := &types.AccountData{}
		isNew := true
		if !reflect.DeepEqual(maybeAccount, types.AccountData{}) {
			account = maybeAccount
			isNew = false
		}

		result := f(account, isNew)
		if result.HasError {
			return result
		}

		maybeEndowed := sc.NewOption[types.Balance](nil)
		if isNew {
			maybeEndowed = sc.NewOption[types.Balance](account.Free)
		}
		maybeAccountWithDust, imbalance := postMutation(*account)
		if !maybeAccountWithDust.HasValue {
			maybeAccount = &types.AccountData{}
		} else {
			maybeAccount.Free = maybeAccountWithDust.Value.Free
			maybeAccount.MiscFrozen = maybeAccountWithDust.Value.MiscFrozen
			maybeAccount.FeeFrozen = maybeAccountWithDust.Value.FeeFrozen
			maybeAccount.Reserved = maybeAccountWithDust.Value.Reserved
		}

		r := sc.NewVaryingData(maybeEndowed, imbalance, result)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    r,
		}
	})
	if result.HasError {
		return result
	}

	resultValue := result.Value.(sc.VaryingData)
	maybeEndowed := resultValue[0].(sc.Option[types.Balance])
	if maybeEndowed.HasValue {
		system.DepositEvent(events.NewEventEndowed(who.FixedSequence, maybeEndowed.Value))
	}
	maybeDust := resultValue[1].(sc.Option[NegativeImbalance])
	dustCleaner := DustCleanerValue{
		AccountId:         who,
		NegativeImbalance: maybeDust.Value,
	}

	r := sc.NewVaryingData(resultValue[2], dustCleaner)

	return sc.Result[sc.Encodable]{HasError: false, Value: r}
}

func postMutation(
	new types.AccountData) (sc.Option[types.AccountData], sc.Option[NegativeImbalance]) {
	total := new.Total()

	if total.Cmp(balances.ExistentialDeposit) < 0 {
		if total.Cmp(constants.Zero) == 0 {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[NegativeImbalance](nil)
		} else {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[NegativeImbalance](NewNegativeImbalance(sc.NewU128FromBigInt(total)))
		}
	}

	return sc.NewOption[types.AccountData](new), sc.NewOption[NegativeImbalance](nil)
}

// totalBalance returns the total storage balance of an account id.
func totalBalance(who types.Address32) *big.Int {
	return system.StorageGetAccount(who.FixedSequence).Data.Total()
}

func reducibleBalance(who types.Address32, keepAlive bool) types.Balance {
	accountData := system.StorageGetAccount(who.FixedSequence).Data

	lockedOrFrozen := accountData.FeeFrozen
	if accountData.FeeFrozen.ToBigInt().Cmp(accountData.MiscFrozen.ToBigInt()) < 0 {
		lockedOrFrozen = accountData.MiscFrozen
	}

	liquid := new(big.Int).Sub(accountData.Free.ToBigInt(), lockedOrFrozen.ToBigInt())
	if liquid.Cmp(accountData.Free.ToBigInt()) > 0 {
		liquid = big.NewInt(0)
	}

	if system.CanDecProviders(who) && !keepAlive {
		return sc.NewU128FromBigInt(liquid)
	}

	existentialDeposit := balances.ExistentialDeposit
	diff := new(big.Int).Sub(accountData.Total(), liquid)

	mustRemainToExist := new(big.Int).Sub(existentialDeposit, diff)

	result := new(big.Int).Sub(liquid, mustRemainToExist)
	if result.Cmp(liquid) > 0 {
		return sc.NewU128FromBigInt(big.NewInt(0))
	}

	return sc.NewU128FromBigInt(result)
}
