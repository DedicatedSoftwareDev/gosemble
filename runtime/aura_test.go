package main

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/stretchr/testify/assert"
)

func Test_Aura_Authorities_Empty(t *testing.T) {
	rt, _ := newTestRuntime(t)

	result, err := rt.Exec("AuraApi_authorities", []byte{})
	assert.NoError(t, err)

	assert.Equal(t, []byte{0}, result)
}

func Test_Aura_Authorities(t *testing.T) {
	pubKey1 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee")
	pubKey2 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ef")

	buffer := &bytes.Buffer{}
	buffer.Write(pubKey1)

	bytesPubKey1, err := common.Read32Bytes(buffer)
	assert.NoError(t, err)

	buffer.Write(pubKey2)
	bytesPubKey2, err := common.Read32Bytes(buffer)
	assert.NoError(t, err)

	authorities := [][32]byte{
		bytesPubKey1,
		bytesPubKey2,
	}

	bytesAuthorities, err := scale.Marshal(authorities)
	assert.NoError(t, err)

	rt, storage := newTestRuntime(t)

	err = (*storage).Put(append(keyAuraHash, keyAuthoritiesHash...), bytesAuthorities)
	assert.NoError(t, err)

	result, err := rt.Exec("AuraApi_authorities", []byte{})
	assert.NoError(t, err)

	assert.Equal(t, bytesAuthorities, result)

	var resultAuthorities [][32]byte
	err = scale.Unmarshal(result, &resultAuthorities)
	assert.NoError(t, err)

	assert.Equal(t, authorities, resultAuthorities)
}
