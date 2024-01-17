package metadata

// Metadata types and their corresponding type id.
const (
	PrimitiveTypesBool = iota
	PrimitiveTypesChar
	PrimitiveTypesString
	PrimitiveTypesU8
	PrimitiveTypesU16
	PrimitiveTypesU32
	PrimitiveTypesU64
	PrimitiveTypesU128
	PrimitiveTypesU256
	PrimitiveTypesI8
	PrimitiveTypesI16
	PrimitiveTypesI32
	PrimitiveTypesI64
	PrimitiveTypesI128
	PrimitiveTypesI256

	TypesFixedSequence4U8
	TypesFixedSequence20U8
	TypesFixedSequence32U8
	TypesFixedSequence64U8
	TypesFixedSequence65U8

	TypesSequenceU8
	TypesFixedU128

	TypesCompactU32
	TypesCompactU64
	TypesCompactU128

	TypesH256

	TypesVecBlockNumEventIndex
	TypesSystemEvent
	TypesSystemEventStorage
	TypesEventRecord

	TypesPhase
	TypesDispatchInfo
	TypesDispatchClass
	TypesPays

	TypesDispatchError
	TypesModuleError
	TypesTokenError
	TypesArithmeticError
	TypesTransactionalError

	TypesBalancesEvent
	TypesBalanceStatus
	TypesVecTopics
	TypesLastRuntimeUpgradeInfo
	TypesSystemErrors
	TypesBlockLength
	TypesPerDispatchClassU32
	TypesDbWeight

	TypesRuntimeApis
	TypesRuntimeVecApis
	TypesApiId

	TypesAuraStorageAuthorities
	TypesSequencePubKeys
	TypesAuthorityId
	TypesSr25519PubKey
	TypesAuraSlot

	TypesBalancesErrors

	TypesTransactionPaymentReleases
	TypesTransactionPaymentEvent

	TypesEmptyTuple
	TypesTupleU32U32
	TypesTupleApiIdU32

	TypesAddress32
	TypesMultiAddress

	TypesAccountData
	TypesAccountInfo

	TypesWeight
	TypesOptionWeight
	TypesPerDispatchClassWeight
	TypesPerDispatchClassWeightsPerClass
	TypesWeightPerClass
	TypesBlockWeights

	TypesDigestItem
	TypesSliceDigestItem
	TypesDigest
	TypesEra

	TypesSignatureEd25519
	TypesSignatureSr25519
	TypesSignatureEcdsa
	TypesMultiSignature

	TypesRuntimeEvent
	TypesRuntimeVersion
	RuntimeCall

	SystemCalls
	TimestampCalls
	GrandpaCalls
	BalancesCalls

	UncheckedExtrinsic
	SignedExtra

	CheckNonZeroSender
	CheckSpecVersion
	CheckTxVersion
	CheckGenesis
	CheckMortality
	CheckNonce
	CheckWeight
	ChargeTransactionPayment

	Runtime
)
