package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RuntimeDispatchInfo struct {
	Weight     Weight
	Class      DispatchClass
	PartialFee Balance
}

func (rdi RuntimeDispatchInfo) Encode(buffer *bytes.Buffer) {
	rdi.Weight.Encode(buffer)
	rdi.Class.Encode(buffer)
	rdi.PartialFee.Encode(buffer)
}

func (rdi RuntimeDispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(rdi)
}

func DecodeRuntimeDispatchInfo(buffer *bytes.Buffer) RuntimeDispatchInfo {
	rdi := RuntimeDispatchInfo{}
	rdi.Weight = DecodeWeight(buffer)
	rdi.Class = DecodeDispatchClass(buffer)
	rdi.PartialFee = sc.DecodeU128(buffer)
	return rdi
}
