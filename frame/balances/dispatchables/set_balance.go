package dispatchables

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type SetBalanceCall struct {
	primitives.Callable
}

func NewSetBalanceCall(args sc.VaryingData) SetBalanceCall {
	call := SetBalanceCall{
		Callable: primitives.Callable{
			ModuleId:   balances.ModuleIndex,
			FunctionId: balances.FunctionSetBalanceIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c SetBalanceCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c SetBalanceCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c SetBalanceCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c SetBalanceCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c SetBalanceCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c SetBalanceCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ SetBalanceCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 17_474 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_777_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ SetBalanceCall) IsInherent() bool {
	return false
}

func (_ SetBalanceCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ SetBalanceCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ SetBalanceCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ SetBalanceCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	newFree := args[1].(sc.Compact)
	newReserved := args[2].(sc.Compact)

	err := setBalance(origin, args[0].(types.MultiAddress), newFree.ToBigInt(), newReserved.ToBigInt())
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

// setBalance sets the balance of a given account.
// Changes free and reserve balance of `who`,
// including the total issuance.
// Can only be called by ROOT.
func setBalance(origin types.RawOrigin, who types.MultiAddress, newFree *big.Int, newReserved *big.Int) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	address, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	existentialDeposit := balances.ExistentialDeposit
	sum := new(big.Int).Add(newFree, newReserved)

	if sum.Cmp(existentialDeposit) < 0 {
		newFree = big.NewInt(0)
		newReserved = big.NewInt(0)
	}

	result := mutateAccount(address, func(acc *types.AccountData, bool bool) sc.Result[sc.Encodable] {
		oldFree := acc.Free
		oldReserved := acc.Reserved

		acc.Free = sc.NewU128FromBigInt(newFree)
		acc.Reserved = sc.NewU128FromBigInt(newReserved)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    sc.NewVaryingData(oldFree, oldReserved),
		}
	})
	parsedResult := result.Value.(sc.VaryingData)
	oldFree := parsedResult[0].(types.Balance)
	oldReserved := parsedResult[1].(types.Balance)

	if newFree.Cmp(oldFree.ToBigInt()) > 0 {
		diff := new(big.Int).Sub(newFree, oldFree.ToBigInt())

		NewPositiveImbalance(sc.NewU128FromBigInt(diff)).Drop()
	} else if newFree.Cmp(oldFree.ToBigInt()) < 0 {
		diff := new(big.Int).Sub(oldFree.ToBigInt(), newFree)

		NewNegativeImbalance(sc.NewU128FromBigInt(diff)).Drop()
	}

	if newReserved.Cmp(oldReserved.ToBigInt()) > 0 {
		diff := new(big.Int).Sub(newReserved, oldReserved.ToBigInt())

		NewPositiveImbalance(sc.NewU128FromBigInt(diff)).Drop()
	} else if newReserved.Cmp(oldReserved.ToBigInt()) < 0 {
		diff := new(big.Int).Sub(oldReserved.ToBigInt(), newReserved)

		NewNegativeImbalance(sc.NewU128FromBigInt(diff)).Drop()
	}

	system.DepositEvent(
		events.NewEventBalanceSet(
			who.AsAddress32().FixedSequence,
			sc.NewU128FromBigInt(newFree),
			sc.NewU128FromBigInt(newReserved),
		),
	)
	return nil
}
