package main

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func Test_Metadata_Encoding_Success(t *testing.T) {
	runtime, _ := newTestRuntime(t)
	gossamerMetadata := runtimeMetadata(t, runtime)

	bMetadata, err := runtime.Metadata()
	assert.NoError(t, err)

	buffer := bytes.NewBuffer(bMetadata)

	// Decode Compact Length
	_ = sc.DecodeCompact(buffer)

	// Copy bytes for assertion after re-encode.
	bMetadataCopy := make([]byte, buffer.Len())
	copy(bMetadataCopy, buffer.Bytes())

	metadata, err := types.DecodeMetadata(buffer)
	assert.NoError(t, err)

	// Assert encoding of previously decoded
	assert.Equal(t, bMetadataCopy, metadata.Bytes())

	// Encode gossamer Metadata
	bGossamerMetadata, err := codec.Encode(gossamerMetadata)
	assert.NoError(t, err)

	assert.Equal(t, metadata.Bytes(), bGossamerMetadata)
}
