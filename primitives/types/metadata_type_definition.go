package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	MetadataTypeDefinitionComposite sc.U8 = iota
	MetadataTypeDefinitionVariant
	MetadataTypeDefinitionSequence
	MetadataTypeDefinitionFixedSequence
	MetadataTypeDefinitionTuple
	MetadataTypeDefinitionPrimitive
	MetadataTypeDefinitionCompact
	MetadataTypeDefinitionBitSequence
)

type MetadataTypeDefinition = sc.VaryingData

func NewMetadataTypeDefinitionComposite(fields sc.Sequence[MetadataTypeDefinitionField]) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionComposite, fields)
}

func NewMetadataTypeDefinitionVariant(variants sc.Sequence[MetadataDefinitionVariant]) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionVariant, variants)
}

func NewMetadataTypeDefinitionSequence(compact sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionSequence, compact)
}

func NewMetadataTypeDefinitionFixedSequence(length sc.U32, typeId sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionFixedSequence, length, typeId)
}

func NewMetadataTypeDefinitionTuple(compacts sc.Sequence[sc.Compact]) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionTuple, compacts)
}

func NewMetadataTypeDefinitionPrimitive(primitive MetadataDefinitionPrimitive) MetadataTypeDefinition {
	// TODO: type safety
	return sc.NewVaryingData(MetadataTypeDefinitionPrimitive, primitive)
}

func NewMetadataTypeDefinitionCompact(compact sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionCompact, compact)
}

func NewMetadataTypeDefinitionBitSequence(storeOrder, orderType sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionBitSequence, storeOrder, orderType)
}

func DecodeMetadataTypeDefinition(buffer *bytes.Buffer) MetadataTypeDefinition {
	b := sc.DecodeU8(buffer)

	switch b {
	case MetadataTypeDefinitionComposite:
		return NewMetadataTypeDefinitionComposite(sc.DecodeSequenceWith(buffer, DecodeMetadataTypeDefinitionField))
	case MetadataTypeDefinitionVariant:
		return NewMetadataTypeDefinitionVariant(sc.DecodeSequenceWith(buffer, DecodeMetadataTypeDefinitionVariant))
	case MetadataTypeDefinitionSequence:
		return NewMetadataTypeDefinitionSequence(sc.DecodeCompact(buffer))
	case MetadataTypeDefinitionFixedSequence:
		return NewMetadataTypeDefinitionFixedSequence(sc.DecodeU32(buffer), sc.DecodeCompact(buffer))
	case MetadataTypeDefinitionTuple:
		return NewMetadataTypeDefinitionTuple(sc.DecodeSequence[sc.Compact](buffer))
	case MetadataTypeDefinitionPrimitive:
		return NewMetadataTypeDefinitionPrimitive(DecodeMetadataDefinitionPrimitive(buffer))
	case MetadataTypeDefinitionCompact:
		return NewMetadataTypeDefinitionCompact(sc.DecodeCompact(buffer))
	case MetadataTypeDefinitionBitSequence:
		return NewMetadataTypeDefinitionBitSequence(sc.DecodeCompact(buffer), sc.DecodeCompact(buffer))
	default:
		log.Critical("invalid MetadataTypeDefinition type")
	}

	panic("unreachable")
}

type MetadataTypeDefinitionField struct {
	Name     sc.Option[sc.Str]
	Type     sc.Compact
	TypeName sc.Option[sc.Str]
	Docs     sc.Sequence[sc.Str]
}

func NewMetadataTypeDefinitionField(id int) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.NewOption[sc.Str](nil),
		Type:     sc.ToCompact(id),
		TypeName: sc.NewOption[sc.Str](nil),
		Docs:     sc.Sequence[sc.Str]{},
	}
}

func NewMetadataTypeDefinitionFieldWithNames(id int, name sc.Str, idName sc.Str) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.NewOption[sc.Str](name),
		Type:     sc.ToCompact(id),
		TypeName: sc.NewOption[sc.Str](idName),
		Docs:     sc.Sequence[sc.Str]{},
	}
}

func NewMetadataTypeDefinitionFieldWithName(id int, idName sc.Str) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.NewOption[sc.Str](nil),
		Type:     sc.ToCompact(id),
		TypeName: sc.NewOption[sc.Str](idName),
		Docs:     sc.Sequence[sc.Str]{},
	}
}

func (mtdf MetadataTypeDefinitionField) Encode(buffer *bytes.Buffer) {
	mtdf.Name.Encode(buffer)
	mtdf.Type.Encode(buffer)
	mtdf.TypeName.Encode(buffer)
	mtdf.Docs.Encode(buffer)
}

func DecodeMetadataTypeDefinitionField(buffer *bytes.Buffer) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.DecodeOption[sc.Str](buffer),
		Type:     sc.DecodeCompact(buffer),
		TypeName: sc.DecodeOption[sc.Str](buffer),
		Docs:     sc.DecodeSequence[sc.Str](buffer),
	}
}

func (mtdf MetadataTypeDefinitionField) Bytes() []byte {
	return sc.EncodedBytes(mtdf)
}

type MetadataDefinitionVariant struct {
	Name   sc.Str
	Fields sc.Sequence[MetadataTypeDefinitionField]
	Index  sc.U8
	Docs   sc.Sequence[sc.Str]
}

func NewMetadataDefinitionVariant(name string, fields sc.Sequence[MetadataTypeDefinitionField], index sc.U8, docs string) MetadataDefinitionVariant {
	return MetadataDefinitionVariant{
		Name:   sc.Str(name),
		Fields: fields,
		Index:  index,
		Docs:   sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mdv MetadataDefinitionVariant) Encode(buffer *bytes.Buffer) {
	mdv.Name.Encode(buffer)
	mdv.Fields.Encode(buffer)
	mdv.Index.Encode(buffer)
	mdv.Docs.Encode(buffer)
}

func DecodeMetadataTypeDefinitionVariant(buffer *bytes.Buffer) MetadataDefinitionVariant {
	return MetadataDefinitionVariant{
		Name:   sc.DecodeStr(buffer),
		Fields: sc.DecodeSequenceWith(buffer, DecodeMetadataTypeDefinitionField),
		Index:  sc.DecodeU8(buffer),
		Docs:   sc.DecodeSequence[sc.Str](buffer),
	}
}

func (mdv MetadataDefinitionVariant) Bytes() []byte {
	return sc.EncodedBytes(mdv)
}

const (
	MetadataDefinitionPrimitiveBoolean MetadataDefinitionPrimitive = iota
	MetadataDefinitionPrimitiveChar
	MetadataDefinitionPrimitiveString
	MetadataDefinitionPrimitiveU8
	MetadataDefinitionPrimitiveU16
	MetadataDefinitionPrimitiveU32
	MetadataDefinitionPrimitiveU64
	MetadataDefinitionPrimitiveU128
	MetadataDefinitionPrimitiveU256
	MetadataDefinitionPrimitiveI8
	MetadataDefinitionPrimitiveI16
	MetadataDefinitionPrimitiveI32
	MetadataDefinitionPrimitiveI64
	MetadataDefinitionPrimitiveI128
	MetadataDefinitionPrimitiveI256
)

type MetadataDefinitionPrimitive = sc.U8

func DecodeMetadataDefinitionPrimitive(buffer *bytes.Buffer) MetadataDefinitionPrimitive {
	b := sc.DecodeU8(buffer)

	switch b {
	case MetadataDefinitionPrimitiveBoolean:
		return MetadataDefinitionPrimitiveBoolean
	case MetadataDefinitionPrimitiveChar:
		return MetadataDefinitionPrimitiveChar
	case MetadataDefinitionPrimitiveString:
		return MetadataDefinitionPrimitiveString
	case MetadataDefinitionPrimitiveU8:
		return MetadataDefinitionPrimitiveU8
	case MetadataDefinitionPrimitiveU16:
		return MetadataDefinitionPrimitiveU16
	case MetadataDefinitionPrimitiveU32:
		return MetadataDefinitionPrimitiveU32
	case MetadataDefinitionPrimitiveU64:
		return MetadataDefinitionPrimitiveU64
	case MetadataDefinitionPrimitiveU128:
		return MetadataDefinitionPrimitiveU128
	case MetadataDefinitionPrimitiveU256:
		return MetadataDefinitionPrimitiveU256
	case MetadataDefinitionPrimitiveI8:
		return MetadataDefinitionPrimitiveI8
	case MetadataDefinitionPrimitiveI16:
		return MetadataDefinitionPrimitiveI16
	case MetadataDefinitionPrimitiveI32:
		return MetadataDefinitionPrimitiveI32
	case MetadataDefinitionPrimitiveI64:
		return MetadataDefinitionPrimitiveI64
	case MetadataDefinitionPrimitiveI128:
		return MetadataDefinitionPrimitiveI128
	case MetadataDefinitionPrimitiveI256:
		return MetadataDefinitionPrimitiveI256

	default:
		log.Critical("invalid MetadataDefinitionPrimitive type")
	}

	panic("unreachable")
}
