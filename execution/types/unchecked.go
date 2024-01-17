package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	// ExtrinsicFormatVersion is the current version of the [`UncheckedExtrinsic`] encoded format.
	//
	// This version needs to be bumped if the encoded representation changes.
	// It ensures that if the representation is changed and the format is not known,
	// the decoding fails.
	ExtrinsicFormatVersion = 4
	ExtrinsicBitSigned     = 0b1000_0000
	ExtrinsicUnmaskVersion = 0b0111_1111
)

type UncheckedExtrinsic struct {
	Version sc.U8

	// The signature, address, number of extrinsics have come before from
	// the same signer and an era describing the longevity of this transaction,
	// if this is a signed extrinsic.
	Signature sc.Option[primitives.ExtrinsicSignature]
	Function  primitives.Call
}

func NewUncheckedExtrinsic(function primitives.Call, signedData sc.Option[primitives.ExtrinsicSignature]) UncheckedExtrinsic {
	if signedData.HasValue {
		address, signature, extra := signedData.Value.Signer, signedData.Value.Signature, signedData.Value.Extra
		return NewSignedUncheckedExtrinsic(function, address, signature, extra)
	} else {
		return NewUnsignedUncheckedExtrinsic(function)
	}
}

// NewSignedUncheckedExtrinsic returns a new instance of a signed extrinsic.
func NewSignedUncheckedExtrinsic(function primitives.Call, address primitives.MultiAddress, signature primitives.MultiSignature, extra primitives.SignedExtra) UncheckedExtrinsic {
	return UncheckedExtrinsic{
		Version: sc.U8(ExtrinsicFormatVersion | ExtrinsicBitSigned),
		Signature: sc.NewOption[primitives.ExtrinsicSignature](
			primitives.ExtrinsicSignature{
				Signer:    address,
				Signature: signature,
				Extra:     extra,
			},
		),
		Function: function,
	}
}

// NewUnsignedUncheckedExtrinsic returns a new instance of an unsigned extrinsic.
func NewUnsignedUncheckedExtrinsic(function primitives.Call) UncheckedExtrinsic {
	return UncheckedExtrinsic{
		Version:   sc.U8(ExtrinsicFormatVersion),
		Signature: sc.NewOption[primitives.ExtrinsicSignature](nil),
		Function:  function,
	}
}

func (uxt UncheckedExtrinsic) UnmaskedVersion() sc.U8 {
	return uxt.Version & ExtrinsicUnmaskVersion
}

func (uxt UncheckedExtrinsic) IsSigned() sc.Bool {
	return uxt.Version&ExtrinsicBitSigned == ExtrinsicBitSigned
}

func (uxt UncheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	tempBuffer := &bytes.Buffer{}

	if uxt.Signature.HasValue {
		sc.U8(ExtrinsicFormatVersion | ExtrinsicBitSigned).Encode(tempBuffer)
		uxt.Signature.Value.Encode(tempBuffer)
	} else {
		sc.U8(ExtrinsicFormatVersion & ExtrinsicUnmaskVersion).Encode(tempBuffer)
	}

	uxt.Function.Encode(tempBuffer)

	encodedLen := sc.ToCompact(uint64(tempBuffer.Len()))
	encodedLen.Encode(buffer)
	buffer.Write(tempBuffer.Bytes())
}

func DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic {
	// This is a little more complicated than usual since the binary format must be compatible
	// with SCALE's generic `Vec<u8>` type. Basically this just means accepting that there
	// will be a prefix of vector length.
	expectedLength := int(sc.DecodeCompact(buffer).ToBigInt().Int64())
	beforeLength := buffer.Len()

	version, _ := buffer.ReadByte()
	isSigned := version&ExtrinsicBitSigned != 0

	if version&ExtrinsicUnmaskVersion != ExtrinsicFormatVersion {
		log.Critical("invalid Extrinsic version")
	}

	var extSignature sc.Option[primitives.ExtrinsicSignature]
	if isSigned {
		extSignature = sc.NewOption[primitives.ExtrinsicSignature](primitives.DecodeExtrinsicSignature(buffer))
	}

	// Decodes the dispatch call, including its arguments.
	function := DecodeCall(buffer)

	afterLength := buffer.Len()

	if expectedLength != beforeLength-afterLength {
		log.Critical("invalid length prefix")
	}

	return UncheckedExtrinsic{
		Version:   sc.U8(version),
		Signature: extSignature,
		Function:  function,
	}
}

func (uxt UncheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(uxt)
}
