package dispatchables

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TransferAllCall struct {
	primitives.Callable
}

func NewTransferAllCall(args sc.VaryingData) TransferAllCall {
	call := TransferAllCall{
		Callable: primitives.Callable{
			ModuleId:   balances.ModuleIndex,
			FunctionId: balances.FunctionTransferAllIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c TransferAllCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeBool(buffer),
	)
	return c
}

func (c TransferAllCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c TransferAllCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c TransferAllCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c TransferAllCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c TransferAllCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ TransferAllCall) IsInherent() bool {
	return false
}

func (_ TransferAllCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 34_878 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(35_121_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ TransferAllCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ TransferAllCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ TransferAllCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ TransferAllCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	err := transferAll(origin, args[0].(types.MultiAddress), bool(args[1].(sc.Bool)))
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

// transferAll transfers the entire transferable balance from `origin` to `dest`.
// By transferable it means that any locked or reserved amounts will not be transferred.
// `keepAlive`: A boolean to determine if the `transfer_all` operation should send all
// the funds the account has, causing the sender account to be killed (false), or
// transfer everything except at least the existential deposit, which will guarantee to
// keep the sender account alive (true).
func transferAll(origin types.RawOrigin, dest types.MultiAddress, keepAlive bool) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	transactor := origin.AsSigned()
	reducibleBalance := reducibleBalance(transactor, keepAlive)

	to, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", dest.Bytes()))
		return types.NewDispatchErrorCannotLookup()
	}

	keep := types.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = types.ExistenceRequirementAllowDeath
	}

	return trans(transactor, to, reducibleBalance, keep)
}
