package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/constants/timestamp"
	ts "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TimestampModule struct {
	functions map[sc.U8]primitives.Call
}

func NewTimestampModule() TimestampModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[ts.FunctionSetIndex] = dispatchables.NewSetCall(nil)

	return TimestampModule{
		functions: functions,
	}
}

func (tm TimestampModule) Functions() map[sc.U8]primitives.Call {
	return tm.functions
}

func (tm TimestampModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (tm TimestampModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (tm TimestampModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return tm.metadataTypes(), primitives.MetadataModule{
		Name: "Timestamp",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "Timestamp",
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"Now",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesU64)),
					"Current time for the current block."),
				primitives.NewMetadataModuleStorageEntry(
					"DidUpdate",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesBool)),
					"Did the timestamp get updated in this block?"),
			},
		}),
		Call:  sc.NewOption[sc.Compact](sc.ToCompact(metadata.TimestampCalls)),
		Event: sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"MinimumPeriod",
				sc.ToCompact(metadata.PrimitiveTypesU64),
				sc.BytesToSequenceU8(sc.U64(ts.MinimumPeriod).Bytes()),
				"The minimum period between blocks. Beware that this is different to the *expected*  period that the block production apparatus provides.",
			),
		},
		Error: sc.NewOption[sc.Compact](nil),
		Index: timestamp.ModuleIndex,
	}
}

func (tm TimestampModule) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParam(metadata.TimestampCalls, "Timestamp calls", sc.Sequence[sc.Str]{"pallet_timestamp", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"set",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "now", "T::Moment"),
					},
					timestamp.FunctionSetIndex,
					"Set the current time."),
			}), primitives.NewMetadataEmptyTypeParameter("T")),
	}
}
