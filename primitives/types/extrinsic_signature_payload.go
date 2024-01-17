package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
)

// SignedPayload A payload that has been signed for an unchecked extrinsics.
//
// Note that the payload that we sign to produce unchecked extrinsic signature
// is going to be different than the `SignaturePayload` - so the thing the extrinsic
// actually contains.
//
// TODO: make it generic
// generic::SignedPayload<RuntimeCall, SignedExtra>;
type SignedPayload struct {
	Call  Call
	Extra SignedExtra
	AdditionalSigned
}

type AdditionalSigned struct {
	SpecVersion sc.U32
	// FormatVersion sc.U32

	// Hh(G): a 32-byte array containing the genesis hash.
	GenesisHash H256 // size 32

	// Hh(B): a 32-byte array containing the hash of the block which starts the mortality period, as described in
	BlockHash H256 // size 32

	TransactionVersion sc.U32
}

func (sp SignedPayload) Encode(buffer *bytes.Buffer) {
	sp.Call.Encode(buffer)
	sp.Extra.Encode(buffer)
	sp.SpecVersion.Encode(buffer)
	sp.TransactionVersion.Encode(buffer)
	// sp.FormatVersion.Encode(buffer)
	sp.GenesisHash.Encode(buffer)
	sp.BlockHash.Encode(buffer)
}

func (sp SignedPayload) Bytes() []byte {
	return sc.EncodedBytes(sp)
}

func (sp SignedPayload) UsingEncoded() sc.Sequence[sc.U8] {
	enc := sp.Bytes()

	if len(enc) > 256 {
		return sc.BytesToSequenceU8(hashing.Blake256(enc))
	} else {
		return sc.BytesToSequenceU8(enc)
	}
}
