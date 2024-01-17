package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/crypto"
)

type Ed25519 struct {
	H512 // size 64
}

func NewEd25519(values ...sc.U8) Ed25519 {
	return Ed25519{NewH512(values...)}
}

func (s Ed25519) Verify(msg sc.Sequence[sc.U8], signer Address32) sc.Bool {
	sig := sc.FixedSequenceU8ToBytes(s.H512.FixedSequence)
	message := sc.SequenceU8ToBytes(msg)
	key := sc.FixedSequenceU8ToBytes(signer.FixedSequence)
	return sc.Bool(crypto.ExtCryptoEd25519VerifyVersion1(sig, message, key))
}

func (s Ed25519) Encode(buffer *bytes.Buffer) {
	s.H512.Encode(buffer)
}

func DecodeEd25519(buffer *bytes.Buffer) Ed25519 {
	s := Ed25519{}
	s.H512 = DecodeH512(buffer)
	return s
}

func (s Ed25519) Bytes() []byte {
	return sc.EncodedBytes(s)
}

type Sr25519 struct {
	H512 // size 64
}

func NewSr25519(values ...sc.U8) Sr25519 {
	return Sr25519{NewH512(values...)}
}

func (s Sr25519) Verify(msg sc.Sequence[sc.U8], signer Address32) sc.Bool {
	sig := sc.FixedSequenceU8ToBytes(s.H512.FixedSequence)
	message := sc.SequenceU8ToBytes(msg)
	key := sc.FixedSequenceU8ToBytes(signer.FixedSequence)
	return sc.Bool(crypto.ExtCryptoSr25519VerifyVersion2(sig, message, key))
}

func (s Sr25519) Encode(buffer *bytes.Buffer) {
	s.H512.Encode(buffer)
}

func DecodeSr25519(buffer *bytes.Buffer) Sr25519 {
	s := Sr25519{}
	s.H512 = DecodeH512(buffer)
	return s
}

func (s Sr25519) Bytes() []byte {
	return sc.EncodedBytes(s)
}

type Ecdsa struct {
	sc.FixedSequence[sc.U8] // size 65
}

func NewEcdsa(values ...sc.U8) Ecdsa {
	return Ecdsa{sc.NewFixedSequence(65, values...)}
}

func (s Ecdsa) Verify(msg sc.Sequence[sc.U8], signer Address32) sc.Bool {
	// TODO:
	return true
}

func (s Ecdsa) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func DecodeEcdsa(buffer *bytes.Buffer) Ecdsa {
	s := Ecdsa{}
	s.FixedSequence = sc.DecodeFixedSequence[sc.U8](65, buffer)
	return s
}

func (s Ecdsa) Bytes() []byte {
	return sc.EncodedBytes(s)
}
