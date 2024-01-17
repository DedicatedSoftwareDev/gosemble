package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type H256 struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewH256(values ...sc.U8) H256 {
	if len(values) != 32 {
		log.Critical("H256 should be of size 32")
	}
	return H256{sc.NewFixedSequence(32, values...)}
}

func (h H256) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
}

func DecodeH256(buffer *bytes.Buffer) H256 {
	h := H256{}
	h.FixedSequence = sc.DecodeFixedSequence[sc.U8](32, buffer)
	return h
}

func (h H256) Bytes() []byte {
	return sc.EncodedBytes(h)
}

type H512 struct {
	sc.FixedSequence[sc.U8] // size 64
}

func NewH512(values ...sc.U8) H512 {
	if len(values) != 64 {
		log.Critical("H512 should be of size 64")
	}
	return H512{sc.NewFixedSequence(64, values...)}
}

func (h H512) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
}

func DecodeH512(buffer *bytes.Buffer) H512 {
	h := H512{}
	h.FixedSequence = sc.DecodeFixedSequence[sc.U8](64, buffer)
	return h
}

func (h H512) Bytes() []byte {
	return sc.EncodedBytes(h)
}

type Blake2bHash struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewBlake2bHash(values ...sc.U8) Blake2bHash {
	if len(values) != 32 {
		log.Critical("Blake2bHash should be of size 32")
	}
	return Blake2bHash{sc.NewFixedSequence(32, values...)}
}

func (h Blake2bHash) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
}

func DecodeBlake2bHash(buffer *bytes.Buffer) Blake2bHash {
	h := Blake2bHash{}
	h.FixedSequence = sc.DecodeFixedSequence[sc.U8](32, buffer)
	return h
}

func (h Blake2bHash) Bytes() []byte {
	return sc.EncodedBytes(h)
}
