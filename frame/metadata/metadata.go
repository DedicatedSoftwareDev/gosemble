package metadata

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/constants/transaction_payment"
	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

// Metadata returns the metadata of the runtime.
// Returns a pointer-size of the SCALE-encoded metadata of the runtime.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-metadata-metadata)
func Metadata() int64 {
	metadata := buildMetadata()
	bMetadata := sc.BytesToSequenceU8(metadata.Bytes())

	return utils.BytesToOffsetAndSize(bMetadata.Bytes())
}

func buildMetadata() primitives.Metadata {
	metadataTypes := append(primitiveTypes(), basicTypes()...)
	metadataTypes = append(metadataTypes, runtimeTypes()...)

	var modules sc.Sequence[primitives.MetadataModule]

	for _, module := range config.Modules {
		mTypes, mModule := module.Metadata()

		metadataTypes = append(metadataTypes, mTypes...)
		modules = append(modules, mModule)
	}

	extrinsic := primitives.MetadataExtrinsic{
		Type:    sc.ToCompact(metadata.UncheckedExtrinsic),
		Version: types.ExtrinsicFormatVersion,
		SignedExtensions: sc.Sequence[primitives.MetadataSignedExtension]{
			primitives.NewMetadataSignedExtension("CheckNonZeroSender", metadata.CheckNonZeroSender, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("CheckSpecVersion", metadata.CheckSpecVersion, metadata.PrimitiveTypesU32),
			primitives.NewMetadataSignedExtension("CheckTxVersion", metadata.CheckTxVersion, metadata.PrimitiveTypesU32),
			primitives.NewMetadataSignedExtension("CheckGenesis", metadata.CheckGenesis, metadata.TypesH256),
			primitives.NewMetadataSignedExtension("CheckMortality", metadata.CheckMortality, metadata.TypesH256),
			primitives.NewMetadataSignedExtension("CheckNonce", metadata.CheckNonce, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("CheckWeight", metadata.CheckWeight, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("ChargeTransactionPayment", metadata.ChargeTransactionPayment, metadata.TypesEmptyTuple),
		},
	}

	runtimeV14Metadata := primitives.RuntimeMetadataV14{
		Types:     metadataTypes,
		Modules:   modules,
		Extrinsic: extrinsic,
		Type:      sc.ToCompact(metadata.Runtime),
	}

	return primitives.NewMetadata(runtimeV14Metadata)
}

// primitiveTypes returns all primitive types
func primitiveTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.PrimitiveTypesBool, "bool", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveBoolean)),
		primitives.NewMetadataType(metadata.PrimitiveTypesChar, "char", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveChar)),
		primitives.NewMetadataType(metadata.PrimitiveTypesString, "string", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveString)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU8, "U8", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU8)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU16, "U16", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU16)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU32, "U32", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU32)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU64, "U64", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU64)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU128, "U128", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU128)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU256, "U256", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU256)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI8, "I8", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI8)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI16, "I16", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI16)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI32, "I32", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI32)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI64, "I64", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI64)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI128, "I128", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI128)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI256, "I256", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI256)),
	}
}

func basicTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.TypesFixedSequence4U8, "[4]byte", primitives.NewMetadataTypeDefinitionFixedSequence(4, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence20U8, "[20]byte", primitives.NewMetadataTypeDefinitionFixedSequence(20, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence32U8, "[32]byte", primitives.NewMetadataTypeDefinitionFixedSequence(32, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence64U8, "[64]byte", primitives.NewMetadataTypeDefinitionFixedSequence(64, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence65U8, "[65]byte", primitives.NewMetadataTypeDefinitionFixedSequence(65, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesSequenceU8, "[]byte", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesCompactU32, "CompactU32", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.TypesCompactU64, "CompactU64", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))),
		primitives.NewMetadataType(metadata.TypesCompactU128, "CompactU128", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))),

		primitives.NewMetadataType(metadata.TypesFixedU128, "FixedU128", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU128, "u128"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesH256, "primitives H256", sc.Sequence[sc.Str]{"primitive_types", "H256"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),

		primitives.NewMetadataTypeWithPath(metadata.TypesAddress32, "Address32", sc.Sequence[sc.Str]{"sp_core", "crypto", "AccountId32"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence32U8, "[u8; 32]")},
		)),

		primitives.NewMetadataTypeWithPath(metadata.TypesAccountData, "AccountData", sc.Sequence[sc.Str]{"pallet_balances", "AccountData"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "Balance"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "Balance"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "misc_frozen", "Balance"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "fee_frozen", "Balance"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesAccountInfo, "AccountInfo", sc.Sequence[sc.Str]{"frame_system", "AccountInfo"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "nonce", "Index"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "consumers", "RefCount"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "providers", "RefCount"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "sufficients", "RefCount"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAccountData, "data", "AccountData"),
			},
		)),

		primitives.NewMetadataTypeWithPath(metadata.TypesWeight, "Weight", sc.Sequence[sc.Str]{"sp_weights", "weight_v2", "Weight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "ref_time", "u64"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "proof_size", "u64"),
			},
		)),
		primitives.NewMetadataTypeWithParam(metadata.TypesOptionWeight, "Option<Weight>", sc.Sequence[sc.Str]{"Option"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"None",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"Option<Weight>(nil)"),
				primitives.NewMetadataDefinitionVariant(
					"Some",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight),
					},
					1,
					"Option<Weight>(value)"),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight, "T"),
		),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassU32,
			"PerDispatchClass[U32]",
			sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "normal", "T"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "operational", "T"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "mandatory", "T"),
				},
			),
			primitives.NewMetadataTypeParameter(metadata.PrimitiveTypesU32, "T"),
		),

		primitives.NewMetadataTypeWithPath(metadata.TypesSignatureEd25519, "SignatureEd25519", sc.Sequence[sc.Str]{"sp_core", "ed25519", "Signature"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence64U8, "[u8; 64]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesSignatureSr25519, "SignatureSr25519", sc.Sequence[sc.Str]{"sp_core", "sr25519", "Signature"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence64U8, "[u8; 64]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesSignatureEcdsa, "SignatureEcdsa", sc.Sequence[sc.Str]{"sp_core", "ecdsa", "Signature"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence65U8, "[u8; 65]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesMultiSignature, "MultiSignature", sc.Sequence[sc.Str]{"sp_runtime", "MultiSignature"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Ed25519",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSignatureEd25519, "ed25519::Signature"),
						},
						primitives.MultiSignatureEd25519,
						"MultiSignature.Ed25519"),
					primitives.NewMetadataDefinitionVariant(
						"Sr25519",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSignatureSr25519, "sr25519::Signature"),
						},
						primitives.MultiSignatureSr25519,
						"MultiSignature.Sr25519"),
					primitives.NewMetadataDefinitionVariant(
						"Ecdsa",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSignatureEcdsa, "ecdsa::Signature"),
						},
						primitives.MultiSignatureEcdsa,
						"MultiSignature.Ecdsa"),
				})),

		primitives.NewMetadataType(metadata.TypesEmptyTuple, "EmptyTuple", primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{},
		)),

		primitives.NewMetadataType(metadata.TypesTupleU32U32, "(U32, U32)",
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.PrimitiveTypesU32), sc.ToCompact(metadata.PrimitiveTypesU32)})),

		primitives.NewMetadataTypeWithParams(metadata.TypesMultiAddress, "MultiAddress", sc.Sequence[sc.Str]{"sp_runtime", "multiaddress", "MultiAddress"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Id",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesAddress32, "AccountId"),
					},
					primitives.MultiAddressId,
					"MultiAddress.Id"),
				primitives.NewMetadataDefinitionVariant(
					"Index",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU32, "AccountIndex"),
					},
					primitives.MultiAddressIndex,
					"MultiAddress.Index"),
				primitives.NewMetadataDefinitionVariant(
					"Raw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.MultiAddressRaw,
					"MultiAddress.Raw"),
				primitives.NewMetadataDefinitionVariant(
					"Address32",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence32U8, "[u8; 32]"),
					},
					primitives.MultiAddress32,
					"MultiAddress.Address32"),
				primitives.NewMetadataDefinitionVariant(
					"Address20",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence20U8, "[u8; 20]"),
					},
					primitives.MultiAddress20,
					"MultiAddress.Address20"),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesAddress32, "AccountId"),
				primitives.NewMetadataTypeParameter(metadata.TypesEmptyTuple, "AccountIndex"),
			}),

		primitives.NewMetadataTypeWithParam(metadata.TypesRuntimeApis, "ApisVec = sp_std::borrow::Cow<'static, [(ApiId, u32)]>;", sc.Sequence[sc.Str]{"Cow"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesRuntimeVecApis),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesRuntimeVecApis, "T"),
		),

		primitives.NewMetadataType(
			metadata.TypesRuntimeVecApis,
			"[(ApiId, u32)]",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleApiIdU32))),

		primitives.NewMetadataType(
			metadata.TypesTupleApiIdU32,
			"(ApiId, u32)",
			primitives.NewMetadataTypeDefinitionTuple(
				sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesApiId), sc.ToCompact(metadata.PrimitiveTypesU32)})),

		primitives.NewMetadataType(
			metadata.TypesApiId,
			"ApiId",
			primitives.NewMetadataTypeDefinitionFixedSequence(8, sc.ToCompact(metadata.PrimitiveTypesU8))),

		primitives.NewMetadataTypeWithPath(metadata.TypesDispatchInfo, "DispatchInfo", sc.Sequence[sc.Str]{"frame_support", "dispatch", "DispatchInfo"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "weight", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchClass, "class", "DispatchClass"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPays, "pays_fee", "Pays"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesDispatchClass, "DispatchClass", sc.Sequence[sc.Str]{"frame_support", "dispatch", "DispatchClass"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Normal",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchClassNormal,
					"DispatchClass.Normal"),
				primitives.NewMetadataDefinitionVariant(
					"Operational",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchClassOperational,
					"DispatchClass.Operational"),
				primitives.NewMetadataDefinitionVariant(
					"Mandatory",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchClassMandatory,
					"DispatchClass.Mandatory"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesPays, "Pays", sc.Sequence[sc.Str]{"frame_support", "dispatch", "Pays"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Yes",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.PaysYes,
					"Pays.Yes"),
				primitives.NewMetadataDefinitionVariant(
					"No",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.PaysNo,
					"Pays.No"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesDispatchError, "DispatchError", sc.Sequence[sc.Str]{"sp_runtime", "DispatchError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Other",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorOther,
					"DispatchError.Other"),
				primitives.NewMetadataDefinitionVariant(
					"Cannotlookup",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorCannotLookup,
					"DispatchError.Cannotlookup"),
				primitives.NewMetadataDefinitionVariant(
					"BadOrigin",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorBadOrigin,
					"DispatchError.BadOrigin"),
				primitives.NewMetadataDefinitionVariant(
					"Module",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesModuleError, "ModuleError"),
					},
					primitives.DispatchErrorModule,
					"DispatchError.Module"),
				primitives.NewMetadataDefinitionVariant(
					"ConsumerRemaining",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorConsumerRemaining,
					"DispatchError.ConsumerRemaining"),
				primitives.NewMetadataDefinitionVariant(
					"NoProviders",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorNoProviders,
					"DispatchError.NoProviders"),
				primitives.NewMetadataDefinitionVariant(
					"TooManyConsumers",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorTooManyConsumers,
					"DispatchError.TooManyConsumers"),
				primitives.NewMetadataDefinitionVariant(
					"Token",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTokenError, "TokenError"),
					},
					primitives.DispatchErrorToken,
					"DispatchError.Token"),
				primitives.NewMetadataDefinitionVariant(
					"Arithmetic",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesArithmeticError, "ArithmeticError"),
					},
					primitives.DispatchErrorArithmetic,
					"DispatchError.Arithmetic"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionalError",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTransactionalError, "TransactionalError"),
					},
					primitives.DispatchErrorTransactional,
					"DispatchError.TransactionalError"),
				primitives.NewMetadataDefinitionVariant(
					"Exhausted",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorExhausted,
					"DispatchError.Exhausted"),
				primitives.NewMetadataDefinitionVariant(
					"Corruption",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorCorruption,
					"DispatchError.Corruption"),
				primitives.NewMetadataDefinitionVariant(
					"Unavailable",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorUnavailable,
					"DispatchError.Unavailable"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesModuleError, "ModuleError", sc.Sequence[sc.Str]{"sp_runtime", "ModuleError"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU8, "index", "u8"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesFixedSequence4U8, "error", "[u8; MAX_MODULE_ERROR_ENCODED_SIZE]"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesTokenError, "TokenError", sc.Sequence[sc.Str]{"sp_runtime", "TokenError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"NoFunds",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorNoFunds,
					"TokenError.NoFunds"),
				primitives.NewMetadataDefinitionVariant(
					"WouldDie",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorWouldDie,
					"TokenError.WouldDie"),
				primitives.NewMetadataDefinitionVariant(
					"Mandatory",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorBelowMinimum,
					"TokenError.BelowMinimum"),
				primitives.NewMetadataDefinitionVariant(
					"CannotCreate",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorCannotCreate,
					"TokenError.CannotCreate"),
				primitives.NewMetadataDefinitionVariant(
					"UnknownAsset",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorUnknownAsset,
					"TokenError.UnknownAsset"),
				primitives.NewMetadataDefinitionVariant(
					"Frozen",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorFrozen,
					"TokenError.Frozen"),
				primitives.NewMetadataDefinitionVariant(
					"Unsupported",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorUnsupported,
					"TokenError.Unsupported"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesArithmeticError, "ArithmeticError", sc.Sequence[sc.Str]{"sp_arithmetic", "ArithmeticError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Underflow",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.ArithmeticErrorUnderflow,
					"ArithmeticError.Underflow"),
				primitives.NewMetadataDefinitionVariant(
					"Overflow",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.ArithmeticErrorOverflow,
					"ArithmeticError.Overflow"),
				primitives.NewMetadataDefinitionVariant(
					"DivisionByZero",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.ArithmeticErrorDivisionByZero,
					"ArithmeticError.DivisionByZero"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesTransactionalError, "TransactionalError", sc.Sequence[sc.Str]{"sp_runtime", "TransactionalError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"LimitReached",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TransactionalErrorLimitReached,
					"TransactionalError.LimitReached"),
				primitives.NewMetadataDefinitionVariant(
					"NoLayer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TransactionalErrorNoLayer,
					"TransactionalError.NoLayer"),
			})),

		primitives.NewMetadataType(metadata.TypesVecTopics, "Vec<Topics>", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesH256))),

		primitives.NewMetadataTypeWithPath(metadata.TypesDigestItem, "DigestItem", sc.Sequence[sc.Str]{"sp_runtime", "generic", "digest", "DigestItem"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"PreRuntime",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestTypePreRuntime,
					"DigestItem.PreRuntime"),
				primitives.NewMetadataDefinitionVariant(
					"Consensus",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestTypeConsensusMessage,
					"DigestItem.Consensus"),
				primitives.NewMetadataDefinitionVariant(
					"Seal",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestTypeSeal,
					"DigestItem.Seal"),
				primitives.NewMetadataDefinitionVariant(
					"Other",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					0,
					"DigestItem.Other"),
				primitives.NewMetadataDefinitionVariant(
					"RuntimeEnvironmentUpdated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DigestTypeRuntimeEnvironmentUpgraded,
					"DigestItem.RuntimeEnvironmentUpdated"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesDigest, "sp_runtime generic digest Digest", sc.Sequence[sc.Str]{"sp_runtime", "generic", "digest", "Digest"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesSliceDigestItem, "logs", "Vec<DigestItem>"),
				})),
		primitives.NewMetadataType(metadata.TypesSliceDigestItem, "Vec<DigestItem>", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesDigestItem))),
	}
}

func runtimeTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesRuntimeEvent, "node_template_runtime RuntimeEvent", sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeEvent"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"System",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSystemEvent, "frame_system::Event<Runtime>"),
					},
					system.ModuleIndex,
					"Events.System"),
				primitives.NewMetadataDefinitionVariant( // TODO:
					"Grandpa",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					grandpa.ModuleIndex,
					"Events.Grandpa"),
				primitives.NewMetadataDefinitionVariant(
					"Balances",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesBalancesEvent, "pallet_balances::Event<Runtime>"),
					},
					balances.ModuleIndex,
					"Events.Balances"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionPayment",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTransactionPaymentEvent, "pallet_transaction_payment::Event<Runtime>"),
					},
					transaction_payment.ModuleIndex,
					"Events.TransactionPayment"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesRuntimeVersion, "sp_version RuntimeVersion", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString), // spec_name
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString), // impl_name
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // authoring_version
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // spec_version
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // impl_version
				primitives.NewMetadataTypeDefinitionField(metadata.TypesRuntimeApis),     // apis
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // transaction_version
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),     // state_version
			})),

		primitives.NewMetadataType(metadata.SignedExtra, "SignedExtra", primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{
				sc.ToCompact(metadata.CheckNonZeroSender),
				sc.ToCompact(metadata.CheckSpecVersion),
				sc.ToCompact(metadata.CheckTxVersion),
				sc.ToCompact(metadata.CheckGenesis),
				sc.ToCompact(metadata.CheckMortality),
				sc.ToCompact(metadata.CheckNonce),
				sc.ToCompact(metadata.CheckWeight),
				sc.ToCompact(metadata.ChargeTransactionPayment),
			})),

		primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
			sc.Sequence[sc.Str]{"sp_runtime", "generic", "unchecked_extrinsic", "UncheckedExtrinsic"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesMultiAddress, "Address"),
				primitives.NewMetadataTypeParameter(metadata.RuntimeCall, "Call"),
				primitives.NewMetadataTypeParameter(metadata.TypesMultiSignature, "Signature"),
				primitives.NewMetadataTypeParameter(metadata.SignedExtra, "Extra"),
			},
		),

		primitives.NewMetadataTypeWithPath(metadata.RuntimeCall, "RuntimeCall", sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeCall"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"System",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SystemCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<System, Runtime>"),
					},
					system.ModuleIndex,
					"Call.System"),
				primitives.NewMetadataDefinitionVariant(
					"Timestamp",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TimestampCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Timestamp, Runtime>"),
					},
					timestamp.ModuleIndex,
					"Call.Timestamp"),
				primitives.NewMetadataDefinitionVariant(
					"Grandpa",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TimestampCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Grandpa, Runtime>"),
					},
					grandpa.ModuleIndex,
					"Call.Grandpa"),
				primitives.NewMetadataDefinitionVariant(
					"Balances",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.BalancesCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Balances, Runtime>"),
					},
					balances.ModuleIndex,
					"Call.Balances"),
			})),
		primitives.NewMetadataType(metadata.Runtime, "Runtime", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
	}
}
