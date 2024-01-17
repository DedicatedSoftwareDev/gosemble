package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type MetadataModule struct {
	Name      sc.Str
	Storage   sc.Option[MetadataModuleStorage]
	Call      sc.Option[sc.Compact]
	Event     sc.Option[sc.Compact]
	Constants sc.Sequence[MetadataModuleConstant]
	Error     sc.Option[sc.Compact]
	Index     sc.U8
}

func (mm MetadataModule) Encode(buffer *bytes.Buffer) {
	mm.Name.Encode(buffer)
	mm.Storage.Encode(buffer)
	mm.Call.Encode(buffer)
	mm.Event.Encode(buffer)
	mm.Constants.Encode(buffer)
	mm.Error.Encode(buffer)
	mm.Index.Encode(buffer)
}

func DecodeMetadataModule(buffer *bytes.Buffer) MetadataModule {
	return MetadataModule{
		Name:      sc.DecodeStr(buffer),
		Storage:   sc.DecodeOptionWith(buffer, DecodeMetadataModuleStorage),
		Call:      sc.DecodeOption[sc.Compact](buffer),
		Event:     sc.DecodeOption[sc.Compact](buffer),
		Constants: sc.DecodeSequenceWith(buffer, DecodeMetadataModuleConstant),
		Error:     sc.DecodeOption[sc.Compact](buffer),
		Index:     sc.DecodeU8(buffer),
	}
}

func (mm MetadataModule) Bytes() []byte {
	return sc.EncodedBytes(mm)
}

type MetadataModuleStorage struct {
	Prefix sc.Str
	Items  sc.Sequence[MetadataModuleStorageEntry]
}

func (mms MetadataModuleStorage) Encode(buffer *bytes.Buffer) {
	mms.Prefix.Encode(buffer)
	mms.Items.Encode(buffer)
}

func DecodeMetadataModuleStorage(buffer *bytes.Buffer) MetadataModuleStorage {
	return MetadataModuleStorage{
		Prefix: sc.DecodeStr(buffer),
		Items:  sc.DecodeSequenceWith(buffer, DecodeMetadataModuleStorageEntry),
	}
}

func (mms MetadataModuleStorage) Bytes() []byte {
	return sc.EncodedBytes(mms)
}

type MetadataModuleStorageEntry struct {
	Name       sc.Str
	Modifier   MetadataModuleStorageEntryModifier
	Definition MetadataModuleStorageEntryDefinition
	Fallback   sc.Sequence[sc.U8]
	Docs       sc.Sequence[sc.Str]
}

func NewMetadataModuleStorageEntry(name string, modifier MetadataModuleStorageEntryModifier, definition MetadataModuleStorageEntryDefinition, docs string) MetadataModuleStorageEntry {
	return MetadataModuleStorageEntry{
		Name:       sc.Str(name),
		Modifier:   modifier,
		Definition: definition,
		Fallback:   sc.Sequence[sc.U8]{},
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mmse MetadataModuleStorageEntry) Encode(buffer *bytes.Buffer) {
	mmse.Name.Encode(buffer)
	mmse.Modifier.Encode(buffer)
	mmse.Definition.Encode(buffer)
	mmse.Fallback.Encode(buffer)
	mmse.Docs.Encode(buffer)
}

func DecodeMetadataModuleStorageEntry(buffer *bytes.Buffer) MetadataModuleStorageEntry {
	return MetadataModuleStorageEntry{
		Name:       sc.DecodeStr(buffer),
		Modifier:   DecodeMetadataModuleStorageEntryModifier(buffer),
		Definition: DecodeMetadataModuleStorageEntryDefinition(buffer),
		Fallback:   sc.DecodeSequence[sc.U8](buffer),
		Docs:       sc.DecodeSequence[sc.Str](buffer),
	}
}

func (mmse MetadataModuleStorageEntry) Bytes() []byte {
	return sc.EncodedBytes(mmse)
}

const (
	MetadataModuleStorageEntryModifierOptional MetadataModuleStorageEntryModifier = iota
	MetadataModuleStorageEntryModifierDefault                                     = 1
)

type MetadataModuleStorageEntryModifier = sc.U8

func DecodeMetadataModuleStorageEntryModifier(buffer *bytes.Buffer) MetadataModuleStorageEntryModifier {
	b := sc.DecodeU8(buffer)

	switch b {
	case MetadataModuleStorageEntryModifierOptional:
		return MetadataModuleStorageEntryModifierOptional
	case MetadataModuleStorageEntryModifierDefault:
		return MetadataModuleStorageEntryModifierDefault
	default:
		log.Critical("invalid DecodeMetadataModuleStorageEntryModifier type")
	}

	panic("unreachable")
}

const (
	MetadataModuleStorageEntryDefinitionPlain sc.U8 = iota
	MetadataModuleStorageEntryDefinitionMap
)

type MetadataModuleStorageEntryDefinition = sc.VaryingData

func NewMetadataModuleStorageEntryDefinitionPlain(key sc.Compact) MetadataModuleStorageEntryDefinition {
	return sc.NewVaryingData(MetadataModuleStorageEntryDefinitionPlain, key)
}

func NewMetadataModuleStorageEntryDefinitionMap(storageHashFuncs sc.Sequence[MetadataModuleStorageHashFunc], key, value sc.Compact) MetadataModuleStorageEntryDefinition {
	return sc.NewVaryingData(MetadataModuleStorageEntryDefinitionMap, storageHashFuncs, key, value)
}

func DecodeMetadataModuleStorageEntryDefinition(buffer *bytes.Buffer) MetadataModuleStorageEntryDefinition {
	b := sc.DecodeU8(buffer)

	switch b {
	case MetadataModuleStorageEntryDefinitionPlain:
		return NewMetadataModuleStorageEntryDefinitionPlain(sc.DecodeCompact(buffer))
	case MetadataModuleStorageEntryDefinitionMap:
		return NewMetadataModuleStorageEntryDefinitionMap(sc.DecodeSequenceWith(buffer, DecodeMetadataModuleStorageHashFunc), sc.DecodeCompact(buffer), sc.DecodeCompact(buffer))
	default:
		log.Critical("invalid MetadataModuleStorageEntryDefinition type")
	}

	panic("unreachable")
}

type MetadataModuleConstant struct {
	Name  sc.Str
	Type  sc.Compact
	Value sc.Sequence[sc.U8]
	Docs  sc.Sequence[sc.Str]
}

func NewMetadataModuleConstant(name string, id sc.Compact, value sc.Sequence[sc.U8], docs string) MetadataModuleConstant {
	return MetadataModuleConstant{
		Name:  sc.Str(name),
		Type:  id,
		Value: value,
		Docs:  sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mmc MetadataModuleConstant) Encode(buffer *bytes.Buffer) {
	mmc.Name.Encode(buffer)
	mmc.Type.Encode(buffer)
	mmc.Value.Encode(buffer)
	mmc.Docs.Encode(buffer)
}

func DecodeMetadataModuleConstant(buffer *bytes.Buffer) MetadataModuleConstant {
	return MetadataModuleConstant{
		Name:  sc.DecodeStr(buffer),
		Type:  sc.DecodeCompact(buffer),
		Value: sc.DecodeSequence[sc.U8](buffer),
		Docs:  sc.DecodeSequence[sc.Str](buffer),
	}
}

func (mmc MetadataModuleConstant) Bytes() []byte {
	return sc.EncodedBytes(mmc)
}

const (
	MetadataModuleStorageHashFuncBlake128 MetadataModuleStorageHashFunc = iota
	MetadataModuleStorageHashFuncBlake256
	MetadataModuleStorageHashFuncMultiBlake128Concat
	MetadataModuleStorageHashFuncXX128
	MetadataModuleStorageHashFuncXX256
	MetadataModuleStorageHashFuncMultiXX64
	MetadataModuleStorageHashFuncIdentity
)

type MetadataModuleStorageHashFunc = sc.U8

func DecodeMetadataModuleStorageHashFunc(buffer *bytes.Buffer) MetadataModuleStorageHashFunc {
	b := sc.DecodeU8(buffer)

	switch b {
	case MetadataModuleStorageHashFuncBlake128:
		return MetadataModuleStorageHashFuncBlake128
	case MetadataModuleStorageHashFuncBlake256:
		return MetadataModuleStorageHashFuncBlake256
	case MetadataModuleStorageHashFuncMultiBlake128Concat:
		return MetadataModuleStorageHashFuncMultiBlake128Concat
	case MetadataModuleStorageHashFuncXX128:
		return MetadataModuleStorageHashFuncXX128
	case MetadataModuleStorageHashFuncXX256:
		return MetadataModuleStorageHashFuncXX256
	case MetadataModuleStorageHashFuncMultiXX64:
		return MetadataModuleStorageHashFuncMultiXX64
	case MetadataModuleStorageHashFuncIdentity:
		return MetadataModuleStorageHashFuncIdentity

	default:
		log.Critical("invalid MetadataModuleStorageHashFunc type")
	}

	panic("unreachable")
}
