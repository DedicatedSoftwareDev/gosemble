package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type BlockNumber = sc.U32

type Header struct {
	ParentHash     Blake2bHash
	Number         BlockNumber
	StateRoot      H256
	ExtrinsicsRoot H256
	Digest         Digest
}

func (h Header) Encode(buffer *bytes.Buffer) {
	h.ParentHash.Encode(buffer)
	sc.ToCompact(uint64(h.Number)).Encode(buffer)
	h.StateRoot.Encode(buffer)
	h.ExtrinsicsRoot.Encode(buffer)
	h.Digest.Encode(buffer)
}

func (h Header) Bytes() []byte {
	buffer := &bytes.Buffer{}
	h.Encode(buffer)

	return buffer.Bytes()
}

func DecodeHeader(buffer *bytes.Buffer) Header {
	parentHash := DecodeBlake2bHash(buffer)
	blockNumber := sc.DecodeCompact(buffer)
	stateRoot := DecodeH256(buffer)
	extrinsicRoot := DecodeH256(buffer)
	digest := DecodeDigest(buffer)

	return Header{
		ParentHash:     parentHash,
		Number:         sc.U32(blockNumber.ToBigInt().Int64()),
		StateRoot:      stateRoot,
		ExtrinsicsRoot: extrinsicRoot,
		Digest:         digest,
	}
}
