package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

type SessionKey struct {
	Key    sc.Sequence[sc.U8]
	TypeId sc.FixedSequence[sc.U8]
}

func NewSessionKey(key []byte, typeId [4]byte) SessionKey {
	return SessionKey{
		Key:    sc.BytesToSequenceU8(key),
		TypeId: sc.BytesToFixedSequenceU8(typeId[:]),
	}
}

func (sk SessionKey) Encode(buffer *bytes.Buffer) {
	sk.Key.Encode(buffer)
	sk.TypeId.Encode(buffer)
}

func DecodeSessionKey(buffer *bytes.Buffer) SessionKey {
	return SessionKey{
		Key:    sc.DecodeSequence[sc.U8](buffer),
		TypeId: sc.DecodeFixedSequence[sc.U8](4, buffer),
	}
}

func (sk SessionKey) Bytes() []byte {
	return sc.EncodedBytes(sk)
}
