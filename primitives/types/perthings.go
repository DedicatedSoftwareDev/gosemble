package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type Perbill struct {
	Percentage sc.U32
}

func (p Perbill) Encode(buffer *bytes.Buffer) {
	p.Percentage.Encode(buffer)
}

func DecodePerbill(buffer *bytes.Buffer) Perbill {
	p := Perbill{}
	p.Percentage = sc.DecodeU32(buffer)
	return p
}

func (p Perbill) Bytes() []byte {
	return sc.EncodedBytes(p)
}

func (p Perbill) Mul(v sc.Encodable) sc.Encodable {
	switch v := v.(type) {
	case sc.U32:
		return ((v / 100) * p.Percentage)
	case Weight:
		return Weight{
			RefTime:   (v.RefTime / 100) * sc.U64(p.Percentage),
			ProofSize: (v.ProofSize / 100) * sc.U64(p.Percentage),
		}
	default:
		log.Critical("unsupported type")
	}

	panic("unreachable")
}
