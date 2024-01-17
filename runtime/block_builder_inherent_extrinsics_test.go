package main

import (
	"bytes"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	timestamp "github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	"github.com/stretchr/testify/assert"
)

func Test_BlockBuilder_Inherent_Extrinsics(t *testing.T) {
	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))

	assert.NoError(t, err)

	call := timestamp.NewSetCall(sc.NewVaryingData(sc.ToCompact(time)))

	expectedExtrinsic := types.NewUnsignedUncheckedExtrinsic(call)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	rt, _ := newTestRuntime(t)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)

	assert.NotNil(t, inherentExt)

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)

	assert.Equal(t, expectedExtrinsic, extrinsic)
}
