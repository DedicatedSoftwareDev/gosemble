package dispatchables

import (
	"bytes"
	"fmt"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ForceFreeCall struct {
	primitives.Callable
}

func NewForceFreeCall(args sc.VaryingData) ForceFreeCall {
	call := ForceFreeCall{
		Callable: primitives.Callable{
			ModuleId:   balances.ModuleIndex,
			FunctionId: balances.FunctionForceFreeIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c ForceFreeCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeU128(buffer),
	)
	return c
}

func (c ForceFreeCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c ForceFreeCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c ForceFreeCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c ForceFreeCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c ForceFreeCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ ForceFreeCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 16_790 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_029_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ ForceFreeCall) IsInherent() bool {
	return false
}

func (_ ForceFreeCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ ForceFreeCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ ForceFreeCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ ForceFreeCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	amount := args[1].(sc.U128)

	err := forceFree(origin, args[0].(types.MultiAddress), amount.ToBigInt())
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

// forceFree frees some balance from a user by force.
// Can only be called by ROOT.
// Consider Substrate fn force_unreserve
func forceFree(origin types.RawOrigin, who types.MultiAddress, amount *big.Int) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	target, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", who.Bytes()))
		return types.NewDispatchErrorCannotLookup()
	}

	force(target, amount)

	return nil
}

// forceFree frees some funds, returning the amount that has not been freed.
func force(who types.Address32, value *big.Int) *big.Int {
	if value.Cmp(constants.Zero) == 0 {
		return big.NewInt(0)
	}

	if totalBalance(who).Cmp(constants.Zero) == 0 {
		return value
	}

	result := system.Mutate(who, func(accountData *types.AccountInfo) sc.Result[sc.Encodable] {
		actual := accountData.Data.Reserved.ToBigInt()
		if value.Cmp(actual) < 0 {
			actual = value
		}

		newReserved := new(big.Int).Sub(accountData.Data.Reserved.ToBigInt(), actual)
		accountData.Data.Reserved = sc.NewU128FromBigInt(newReserved)

		// TODO: defensive_saturating_add
		newFree := new(big.Int).Add(accountData.Data.Free.ToBigInt(), actual)
		accountData.Data.Free = sc.NewU128FromBigInt(newFree)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    sc.NewU128FromBigInt(actual),
		}
	})

	actual := result.Value.(sc.U128)

	if result.HasError {
		return value
	}

	system.DepositEvent(events.NewEventUnreserved(who.FixedSequence, actual))

	return new(big.Int).Sub(value, actual.ToBigInt())
}
