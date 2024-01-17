package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type ApiItem struct {
	Name    sc.FixedSequence[sc.U8] // size 8
	Version sc.U32
}

func (ai ApiItem) Encode(buffer *bytes.Buffer) {
	ai.Name.Encode(buffer)
	ai.Version.Encode(buffer)
}

func (ai ApiItem) Bytes() []byte {
	buffer := &bytes.Buffer{}
	ai.Encode(buffer)

	return buffer.Bytes()
}

func DecodeApiItem(buffer *bytes.Buffer) ApiItem {
	return ApiItem{
		Name:    sc.DecodeFixedSequence[sc.U8](8, buffer),
		Version: sc.DecodeU32(buffer),
	}
}

type RuntimeVersion struct {
	SpecName           sc.Str
	ImplName           sc.Str
	AuthoringVersion   sc.U32
	SpecVersion        sc.U32
	ImplVersion        sc.U32
	Apis               sc.Sequence[ApiItem]
	TransactionVersion sc.U32
	StateVersion       sc.U8
}

func (rv RuntimeVersion) Encode(buffer *bytes.Buffer) {
	rv.SpecName.Encode(buffer)
	rv.ImplName.Encode(buffer)
	rv.AuthoringVersion.Encode(buffer)
	rv.SpecVersion.Encode(buffer)
	rv.ImplVersion.Encode(buffer)
	rv.Apis.Encode(buffer)
	rv.TransactionVersion.Encode(buffer)
	rv.StateVersion.Encode(buffer)
}

func DecodeRuntimeVersion(buffer *bytes.Buffer) RuntimeVersion {
	var rv RuntimeVersion

	rv.SpecName = sc.DecodeStr(buffer)
	rv.ImplName = sc.DecodeStr(buffer)
	rv.AuthoringVersion = sc.DecodeU32(buffer)
	rv.SpecVersion = sc.DecodeU32(buffer)
	rv.ImplVersion = sc.DecodeU32(buffer)

	apisLength := sc.DecodeCompact(buffer).ToBigInt().Int64()
	if apisLength != 0 {
		var apis []ApiItem
		for i := 0; i < int(apisLength); i++ {
			apis = append(apis, DecodeApiItem(buffer))
		}
		rv.Apis = apis
	}

	rv.TransactionVersion = sc.DecodeU32(buffer)
	rv.StateVersion = sc.DecodeU8(buffer)

	return rv
}

func (rv RuntimeVersion) Bytes() []byte {
	return sc.EncodedBytes(rv)
}
