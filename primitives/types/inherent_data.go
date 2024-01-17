package types

import (
	"bytes"
	"errors"
	"sort"

	sc "github.com/LimeChain/goscale"
)

type InherentData struct {
	Data map[[8]byte]sc.Sequence[sc.U8]
}

func NewInherentData() *InherentData {
	return &InherentData{
		Data: make(map[[8]byte]sc.Sequence[sc.U8]),
	}
}

func (id *InherentData) Encode(buffer *bytes.Buffer) {
	sc.ToCompact(uint64(len(id.Data))).Encode(buffer)

	keys := make([][8]byte, 0)
	for k := range id.Data {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return string(keys[i][:]) < string(keys[j][:]) })

	for _, k := range keys {
		value := id.Data[k]

		buffer.Write(k[:])
		buffer.Write(value.Bytes())
	}
}

func (id *InherentData) Bytes() []byte {
	return sc.EncodedBytes(id)
}

func (id *InherentData) Put(key [8]byte, value sc.Encodable) error {
	if id.Data[key] != nil {
		return NewInherentErrorInherentDataExists(sc.BytesToSequenceU8(key[:]))
	}

	id.Data[key] = sc.BytesToSequenceU8(value.Bytes())

	return nil
}

func (id *InherentData) Clear() {
	id.Data = make(map[[8]byte]sc.Sequence[sc.U8])
}

func DecodeInherentData(buffer *bytes.Buffer) (*InherentData, error) {
	result := NewInherentData()
	length := sc.DecodeCompact(buffer).ToBigInt().Int64()

	for i := 0; i < int(length); i++ {
		key := [8]byte{}
		len, err := buffer.Read(key[:])
		if err != nil {
			return nil, err
		}
		if len != 8 {
			return nil, errors.New("invalid length")
		}
		value := sc.DecodeSequence[sc.U8](buffer)

		result.Data[key] = value
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
