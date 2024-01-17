package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	key0 = [8]byte{
		't', 'e', 's', 't', 'i', 'n', 'h', '0',
	}
	key1 = [8]byte{
		't', 'e', 's', 't', 'i', 'n', 'h', '1',
	}

	value0 = sc.Sequence[sc.I32]{1, 2, 3}
	value1 = sc.U32(7)

	expectedEncoded = []byte{8, 116, 101, 115, 116, 105, 110, 104, 48, 52, 12, 1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0, 116, 101, 115, 116, 105, 110, 104, 49, 16, 7, 0, 0, 0}
)

func Test_InherentData_Encode(t *testing.T) {
	inherent := NewInherentData()
	assert.Nil(t, inherent.Put(key0, value0))
	assert.Nil(t, inherent.Put(key1, value1))

	encoded := inherent.Bytes()

	assert.Equal(t, expectedEncoded, encoded)
}

func Test_InherentData_Decode(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.Write(expectedEncoded)

	res, err := DecodeInherentData(buffer)
	assert.Nil(t, err)

	buffer.Reset()
	buffer.Write(sc.SequenceU8ToBytes(res.Data[key0]))

	decodedValue0 := sc.DecodeSequence[sc.I32](buffer)
	assert.Equal(t, value0, decodedValue0)

	buffer.Reset()
	buffer.Write(sc.SequenceU8ToBytes(res.Data[key1]))

	decodedValue1 := sc.DecodeU32(buffer)
	assert.Equal(t, value1, decodedValue1)
}
