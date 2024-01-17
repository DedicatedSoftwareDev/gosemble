package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"

	"github.com/LimeChain/gosemble/constants/testable"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TestCall struct {
	primitives.Callable
}

func NewTestCall(args sc.VaryingData) TestCall {
	call := TestCall{
		Callable: primitives.Callable{
			ModuleId:   testable.ModuleIndex,
			FunctionId: testable.FunctionTestIndex,
		},
	}

	if len(args) != 0 {
		call.Arguments = args
	}

	return call
}

func (c TestCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(sc.DecodeSequence[sc.U8](buffer))
	return c
}

func (c TestCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c TestCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c TestCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c TestCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c TestCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ TestCall) BaseWeight(args ...any) types.Weight {
	return types.WeightFromParts(1_000_000, 0)
}

func (_ TestCall) IsInherent() bool {
	return false
}

func (_ TestCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ TestCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ TestCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (_ TestCall) Dispatch(origin types.RuntimeOrigin, _ sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	storage.Set([]byte("testvalue"), []byte{1})

	support.WithStorageLayer(func() (ok types.PostDispatchInfo, err types.DispatchError) {
		storage.Set([]byte("testvalue"), []byte{2})
		return ok, types.NewDispatchErrorOther("revert")
	})

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{Ok: types.PostDispatchInfo{}}
}
