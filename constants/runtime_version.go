package constants

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// If the runtime behavior changes, increment spec_version and set impl_version to 0.
// If only runtime implementation changes and behavior does not,
// then leave spec_version as is and increment impl_version.

const SpecName = "node-template"
const ImplName = "node-template"
const AuthoringVersion = 1
const SpecVersion = 100
const ImplVersion = 1
const TransactionVersion = 1
const StateVersion = 1
const StorageVersion = 0

const BlockHashCount = sc.U32(2400)

// RuntimeVersion contains the version identifiers of the Runtime.
var RuntimeVersion = types.RuntimeVersion{
	SpecName:         sc.Str(SpecName),
	ImplName:         sc.Str(ImplName),
	AuthoringVersion: sc.U32(AuthoringVersion),
	SpecVersion:      sc.U32(SpecVersion),
	ImplVersion:      sc.U32(ImplVersion),
	// Api Names are Blake2bHash8("ApiName")
	// Example: common.MustBlake2b8([]byte("Core") -> [223 106 203 104 153 7 96 155]
	Apis: sc.Sequence[types.ApiItem]{
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 223, 106, 203, 104, 153, 7, 96, 155), // Core
			Version: sc.U32(4),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 55, 227, 151, 252, 124, 145, 245, 228), // Metadata
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 64, 254, 58, 212, 1, 248, 149, 154), // BlockBuilder
			Version: sc.U32(6),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 210, 188, 152, 151, 238, 208, 143, 21), // TaggedTransactionQueue
			Version: sc.U32(3),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 247, 139, 39, 139, 229, 63, 69, 76), // OffchainWorkerApi
			Version: sc.U32(2),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 221, 113, 141, 92, 197, 50, 98, 212), // AuraApi
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 171, 60, 5, 114, 41, 31, 235, 139), // SessionKeys
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 237, 153, 197, 172, 178, 94, 237, 245), // GrandpaApi
			Version: sc.U32(3),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 188, 157, 137, 144, 79, 91, 146, 63), // AccountNonceApi
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 55, 200, 187, 19, 80, 169, 162, 168), // TransactionPaymentApi
			Version: sc.U32(3),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 243, 255, 20, 213, 171, 82, 112, 89), // TransactionPaymentCallApi
			Version: sc.U32(3),
		},
	},
	TransactionVersion: sc.U32(TransactionVersion),
	StateVersion:       sc.U8(StateVersion),
}
