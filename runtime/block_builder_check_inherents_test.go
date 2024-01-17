package main

import (
	"bytes"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_CheckInherents(t *testing.T) {
	expectedCheckInherentsResult := types.NewCheckInherentsResult()

	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))

	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	rt, _ := newTestRuntime(t)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	header := gossamertypes.NewHeader(parentHash, stateRoot, extrinsicsRoot, blockNumber, gossamertypes.NewDigest())

	var exts [][]byte
	err = scale.Unmarshal(inherentExt, &exts)
	assert.NoError(t, err)

	block := gossamertypes.Block{
		Header: *header,
		Body:   gossamertypes.BytesArrayToExtrinsics(exts),
	}

	encodedBlock, err := scale.Marshal(block)
	assert.NoError(t, err)

	inputData := append(encodedBlock, ienc...)
	bytesCheckInherentsResult, err := rt.Exec("BlockBuilder_check_inherents", inputData)
	assert.NoError(t, err)

	buffer := &bytes.Buffer{}
	buffer.Write(bytesCheckInherentsResult)
	checkInherentsResult := types.DecodeCheckInherentsResult(buffer)

	assert.Equal(t, expectedCheckInherentsResult, checkInherentsResult)
}
