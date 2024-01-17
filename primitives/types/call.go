package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Call interface {
	sc.Encodable

	ModuleIndex() sc.U8
	FunctionIndex() sc.U8
	Args() sc.VaryingData
	Dispatch(origin RuntimeOrigin, args sc.VaryingData) DispatchResultWithPostInfo[PostDispatchInfo]
	IsInherent() bool
	BaseWeight(...any) Weight
	ClassifyDispatch(baseWeight Weight) DispatchClass
	PaysFee(baseWeight Weight) Pays
	WeightInfo(baseWeight Weight) Weight
	DecodeArgs(buffer *bytes.Buffer) Call
}

type Callable struct {
	ModuleId   sc.U8
	FunctionId sc.U8
	Arguments  sc.VaryingData
}

func (c Callable) Encode(buffer *bytes.Buffer) {
	c.ModuleId.Encode(buffer)
	c.FunctionId.Encode(buffer)
	c.Arguments.Encode(buffer)
}

func (c Callable) Bytes() []byte {
	return sc.EncodedBytes(c)
}

func (c Callable) ModuleIndex() sc.U8 {
	return c.ModuleId
}

func (c Callable) FunctionIndex() sc.U8 {
	return c.FunctionId
}

func (c Callable) Args() sc.VaryingData {
	return c.Arguments
}
