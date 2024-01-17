package session_keys

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/primitives/crypto"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

// GenerateSessionKeys generates a set of session keys with an optional seed.
// The keys should be stored within the keystore exposed by the Host Api.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded optional seed.
// Returns a pointer-size of the SCALE-encoded set of keys.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-sessionkeys_generate_session_keys)
func GenerateSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	seed := sc.DecodeOptionWith(buffer, sc.DecodeSequence[sc.U8])

	auraPubKey := crypto.ExtCryptoSr25519GenerateVersion1(aura.KeyTypeId[:], seed.Bytes())
	grandpaPubKey := crypto.ExtCryptoEd25519GenerateVersion1(grandpa.KeyTypeId[:], seed.Bytes())

	res := sc.BytesToSequenceU8(append(auraPubKey, grandpaPubKey...))

	return utils.BytesToOffsetAndSize(res.Bytes())
}

// DecodeSessionKeys decodes the given session keys.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded keys.
// Returns a pointer-size of the SCALE-encoded set of raw keys and their respective key type.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-sessionkeys_decode_session_keys)
func DecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)
	sequence := sc.DecodeSequenceWith(buffer, sc.DecodeU8)

	buffer = bytes.NewBuffer(sc.SequenceU8ToBytes(sequence))
	sessionKeys := sc.Sequence[types.SessionKey]{
		types.NewSessionKey(sc.FixedSequenceU8ToBytes(types.DecodePublicKey(buffer)), aura.KeyTypeId),
		types.NewSessionKey(sc.FixedSequenceU8ToBytes(types.DecodePublicKey(buffer)), grandpa.KeyTypeId),
	}

	result := sc.NewOption[sc.Sequence[types.SessionKey]](sessionKeys)
	return utils.BytesToOffsetAndSize(result.Bytes())
}
