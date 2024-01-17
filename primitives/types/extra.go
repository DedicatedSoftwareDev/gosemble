package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// SignedExtra implements SignedExtension
// Extra data, E, is a tuple containing additional metadata about the extrinsic and the system it is meant to be executed in.
type SignedExtra struct {
	Era Era

	// a compact integer containing the nonce of the sender.
	// The nonce must be incremented by one for each extrinsic created,
	// otherwise the Polkadot network will reject the extrinsic.
	Nonce sc.U32 // encode as Compact

	// a compact integer containing the transactor pay including tip.
	Fee sc.U128 // encode as Compact
}

func (e SignedExtra) Encode(buffer *bytes.Buffer) {
	e.Era.Encode(buffer)
	sc.ToCompact(e.Nonce).Encode(buffer)
	sc.Compact(e.Fee).Encode(buffer)
}

func DecodeExtra(buffer *bytes.Buffer) SignedExtra {
	e := SignedExtra{}
	e.Era = DecodeEra(buffer)
	e.Nonce = sc.U32(sc.U128(sc.DecodeCompact(buffer)).ToBigInt().Uint64())
	e.Fee = sc.U128(sc.DecodeCompact(buffer))
	return e
}

func (e SignedExtra) Bytes() []byte {
	return sc.EncodedBytes(e)
}
