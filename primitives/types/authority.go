package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

type Authority struct {
	Id     PublicKey
	Weight sc.U64
}

func (a Authority) Encode(buffer *bytes.Buffer) {
	a.Id.Encode(buffer)
	a.Weight.Encode(buffer)
}

func DecodeAuthority(buffer *bytes.Buffer) Authority {
	return Authority{
		Id:     DecodePublicKey(buffer),
		Weight: sc.DecodeU64(buffer),
	}
}

func (a Authority) Bytes() []byte {
	return sc.EncodedBytes(a)
}
