package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TransferKeepAliveCall struct {
	primitives.Callable
}

func NewTransferKeepAliveCall(args sc.VaryingData) TransferKeepAliveCall {
	call := TransferKeepAliveCall{
		Callable: primitives.Callable{
			ModuleId:   balances.ModuleIndex,
			FunctionId: balances.FunctionTransferKeepAliveIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c TransferKeepAliveCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c TransferKeepAliveCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c TransferKeepAliveCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c TransferKeepAliveCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c TransferKeepAliveCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c TransferKeepAliveCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ TransferKeepAliveCall) IsInherent() bool {
	return false
}

func (_ TransferKeepAliveCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 28_184 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(49_250_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ TransferKeepAliveCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ TransferKeepAliveCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ TransferKeepAliveCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ TransferKeepAliveCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := transferKeepAlive(origin, args[0].(types.MultiAddress), value)
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

// transferKeepAlive is similar to transfer, but includes a check that the origin transactor will not be "killed".
func transferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}
	transactor := origin.AsSigned()

	address, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return trans(transactor, address, value, types.ExistenceRequirementKeepAlive)
}
