package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ForceTransferCall struct {
	primitives.Callable
}

func NewForceTransferCall(args sc.VaryingData) ForceTransferCall {
	call := ForceTransferCall{
		Callable: primitives.Callable{
			ModuleId:   balances.ModuleIndex,
			FunctionId: balances.FunctionForceTransferIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c ForceTransferCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c ForceTransferCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c ForceTransferCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c ForceTransferCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c ForceTransferCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c ForceTransferCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ ForceTransferCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `135`
	//  Estimated: `6196`
	// Minimum execution time: 39_713 nanoseconds.
	r := constants.DbWeight.Reads(2)
	w := constants.DbWeight.Writes(2)
	e := types.WeightFromParts(0, 6196)
	return types.WeightFromParts(40_360_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ ForceTransferCall) IsInherent() bool {
	return false
}

func (_ ForceTransferCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ ForceTransferCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ ForceTransferCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ ForceTransferCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[2].(sc.Compact))

	err := forceTransfer(origin, args[0].(types.MultiAddress), args[1].(types.MultiAddress), value)
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

// forceTransfer transfers liquid free balance from `source` to `dest`.
// Can only be called by ROOT.
func forceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	sourceAddress, err := types.DefaultAccountIdLookup().Lookup(source)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}
	destinationAddress, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return trans(sourceAddress, destinationAddress, value, types.ExistenceRequirementAllowDeath)
}
