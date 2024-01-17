package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type MetadataExtrinsic struct {
	Type             sc.Compact
	Version          sc.U8
	SignedExtensions sc.Sequence[MetadataSignedExtension]
}

func (me MetadataExtrinsic) Encode(buffer *bytes.Buffer) {
	me.Type.Encode(buffer)
	me.Version.Encode(buffer)
	me.SignedExtensions.Encode(buffer)
}

func DecodeMetadataExtrinsic(buffer *bytes.Buffer) MetadataExtrinsic {
	return MetadataExtrinsic{
		Type:             sc.DecodeCompact(buffer),
		Version:          sc.DecodeU8(buffer),
		SignedExtensions: sc.DecodeSequenceWith(buffer, DecodeMetadataSignedExtension),
	}
}

func (me MetadataExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(me)
}

type MetadataSignedExtension struct {
	Identifier       sc.Str
	Type             sc.Compact
	AdditionalSigned sc.Compact
}

func NewMetadataSignedExtension(identifier sc.Str, typeIndex, additionalSigned int) MetadataSignedExtension {
	return MetadataSignedExtension{
		Identifier:       identifier,
		Type:             sc.ToCompact(typeIndex),
		AdditionalSigned: sc.ToCompact(additionalSigned),
	}
}

func (mse MetadataSignedExtension) Encode(buffer *bytes.Buffer) {
	mse.Identifier.Encode(buffer)
	mse.Type.Encode(buffer)
	mse.AdditionalSigned.Encode(buffer)
}

func DecodeMetadataSignedExtension(buffer *bytes.Buffer) MetadataSignedExtension {
	return MetadataSignedExtension{
		Identifier:       sc.DecodeStr(buffer),
		Type:             sc.DecodeCompact(buffer),
		AdditionalSigned: sc.DecodeCompact(buffer),
	}
}

func (mse MetadataSignedExtension) Bytes() []byte {
	return sc.EncodedBytes(mse)
}
