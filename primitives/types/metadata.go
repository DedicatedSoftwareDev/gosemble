package types

import (
	"bytes"
	"errors"
	"fmt"

	sc "github.com/LimeChain/goscale"
)

const (
	MetadataReserved sc.U32 = 0x6174656d // "meta"
	MetadataVersion  sc.U8  = 14
)

type Metadata struct {
	Data RuntimeMetadataV14
}

func NewMetadata(data RuntimeMetadataV14) Metadata {
	return Metadata{Data: data}
}

func (m Metadata) Encode(buffer *bytes.Buffer) {
	MetadataReserved.Encode(buffer)
	MetadataVersion.Encode(buffer)
	m.Data.Encode(buffer)
}

func DecodeMetadata(buffer *bytes.Buffer) (Metadata, error) {
	metaReserved := sc.DecodeU32(buffer)
	if metaReserved != MetadataReserved {
		return Metadata{}, errors.New(fmt.Sprintf("metadata reserved mismatch: expect [%d], actual [%d]", MetadataReserved, metaReserved))
	}

	version := sc.DecodeU8(buffer)
	if version != MetadataVersion {
		return Metadata{}, errors.New(fmt.Sprintf("metadata version mismatch: expect [%d], actual [%d]", MetadataVersion, version))
	}

	return Metadata{
		Data: DecodeRuntimeMetadataV14(buffer),
	}, nil
}

func (m Metadata) Bytes() []byte {
	return sc.EncodedBytes(m)
}

type RuntimeMetadataV14 struct {
	Types     sc.Sequence[MetadataType]
	Modules   sc.Sequence[MetadataModule]
	Extrinsic MetadataExtrinsic
	Type      sc.Compact
}

func (rm RuntimeMetadataV14) Encode(buffer *bytes.Buffer) {
	rm.Types.Encode(buffer)
	rm.Modules.Encode(buffer)
	rm.Extrinsic.Encode(buffer)
	rm.Type.Encode(buffer)
}

func DecodeRuntimeMetadataV14(buffer *bytes.Buffer) RuntimeMetadataV14 {
	return RuntimeMetadataV14{
		Types:     sc.DecodeSequenceWith(buffer, DecodeMetadataType),
		Modules:   sc.DecodeSequenceWith(buffer, DecodeMetadataModule),
		Extrinsic: DecodeMetadataExtrinsic(buffer),
		Type:      sc.DecodeCompact(buffer),
	}
}

func (rm RuntimeMetadataV14) Bytes() []byte {
	return sc.EncodedBytes(rm)
}

type MetadataType struct {
	Id         sc.Compact
	Path       sc.Sequence[sc.Str]
	Params     sc.Sequence[MetadataTypeParameter]
	Definition MetadataTypeDefinition
	Docs       sc.Sequence[sc.Str]
}

func NewMetadataType(id int, docs string, definition MetadataTypeDefinition) MetadataType {
	return MetadataType{
		Id:         sc.ToCompact(id),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func NewMetadataTypeWithPath(id int, docs string, path sc.Sequence[sc.Str], definition MetadataTypeDefinition) MetadataType {
	return MetadataType{
		Id:         sc.ToCompact(id),
		Path:       path,
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func NewMetadataTypeWithParam(id int, docs string, path sc.Sequence[sc.Str], definition MetadataTypeDefinition, param MetadataTypeParameter) MetadataType {
	return MetadataType{
		Id:   sc.ToCompact(id),
		Path: path,
		Params: sc.Sequence[MetadataTypeParameter]{
			param,
		},
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func NewMetadataTypeWithParams(id int, docs string, path sc.Sequence[sc.Str], definition MetadataTypeDefinition, params sc.Sequence[MetadataTypeParameter]) MetadataType {
	return MetadataType{
		Id:         sc.ToCompact(id),
		Path:       path,
		Params:     params,
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mt MetadataType) Encode(buffer *bytes.Buffer) {
	mt.Id.Encode(buffer)
	mt.Path.Encode(buffer)
	mt.Params.Encode(buffer)
	mt.Definition.Encode(buffer)
	mt.Docs.Encode(buffer)
}

func DecodeMetadataType(buffer *bytes.Buffer) MetadataType {
	return MetadataType{
		Id:         sc.DecodeCompact(buffer),
		Path:       sc.DecodeSequence[sc.Str](buffer),
		Params:     sc.DecodeSequenceWith(buffer, DecodeMetadataTypeParameter),
		Definition: DecodeMetadataTypeDefinition(buffer),
		Docs:       sc.DecodeSequence[sc.Str](buffer),
	}
}

func (mt MetadataType) Bytes() []byte {
	return sc.EncodedBytes(mt)
}

type MetadataTypeParameter struct {
	Text sc.Str
	Type sc.Option[sc.Compact]
}

func NewMetadataTypeParameter(id int, text string) MetadataTypeParameter {
	return MetadataTypeParameter{
		Text: sc.Str(text),
		Type: sc.NewOption[sc.Compact](sc.ToCompact(id)),
	}
}

func NewMetadataEmptyTypeParameter(text string) MetadataTypeParameter {
	return MetadataTypeParameter{
		Type: sc.NewOption[sc.Compact](nil),
		Text: sc.Str(text),
	}
}

func (mtp MetadataTypeParameter) Encode(buffer *bytes.Buffer) {
	mtp.Text.Encode(buffer)
	mtp.Type.Encode(buffer)
}

func DecodeMetadataTypeParameter(buffer *bytes.Buffer) MetadataTypeParameter {
	return MetadataTypeParameter{
		Text: sc.DecodeStr(buffer),
		Type: sc.DecodeOption[sc.Compact](buffer),
	}
}

func (mtp MetadataTypeParameter) Bytes() []byte {
	return sc.EncodedBytes(mtp)
}
