package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// ExtrinsicSignature The signature is a varying data type indicating the used signature type,
// followed by the signature created by the extrinsic author (the sender).
type ExtrinsicSignature struct {
	// is the 32-byte address of the sender of the extrinsic
	// as described in https://docs.substrate.io/reference/address-formats/
	Signer    MultiAddress
	Signature MultiSignature
	Extra     SignedExtra
}

func (s ExtrinsicSignature) Encode(buffer *bytes.Buffer) {
	s.Signer.Encode(buffer)
	s.Signature.Encode(buffer)
	s.Extra.Encode(buffer)
}

func DecodeExtrinsicSignature(buffer *bytes.Buffer) ExtrinsicSignature {
	s := ExtrinsicSignature{}
	s.Signer = DecodeMultiAddress(buffer)
	s.Signature = DecodeMultiSignature(buffer)
	s.Extra = DecodeExtra(buffer)
	return s
}

func (s ExtrinsicSignature) Bytes() []byte {
	return sc.EncodedBytes(s)
}
